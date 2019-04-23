package main

import (
	"fmt"
	"github.com/blinkist/go-dockerpty"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Container struct {
	/* Services to facilitate the running container */
	Client    *docker.Client
	Container *docker.Container
	Logger    logrus.FieldLogger

	/** Channels to receive state on */
	Control *FanoutControlChan

	/** Configurations specific to the container being created */
	Binds      []string
	EnvVars    []string
	Image      string
	WorkingDir string
}

func NewContainer(l logrus.FieldLogger, c *docker.Client, dc DockerTerminalConfig) (*Container, error) {
	// Default to the local connection if a nil docker client is passed in
	if c == nil {
		c, _ = docker.NewClient("unix:///var/run/docker.sock")
	}

	ret := &Container{
		Binds: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
			fmt.Sprintf("%s:/home/docker/share/", os.Getenv("HOME")),
		},
		Client: c,
		EnvVars: []string{
			"IS_CONTAINER=true",
		},
		Image:      "terminal:latest",
		Logger:     l,
		WorkingDir: "/home/docker/",

		/** Setup channels */
		Control: NewFanoutControlChan(),
	}

	if err := ret.populateProjectPaths(); err != nil {
		return nil, err
	}

	if err := ret.AddRequiredMount(CONFIG.RequiredMounts...); err != nil {
		return nil, errors.Wrapf(err, "Error adding required mount")
	}

	if err := ret.AddOptionalMount(CONFIG.OptionalMounts...); err != nil {
		return nil, errors.Wrapf(err, "Error adding optional mount")
	}

	if len(dc.EnvVars) > 0 {
		for k, v := range dc.EnvVars {
			if !strings.Contains(v, "=") {
				ret.EnvVars = append(ret.EnvVars, fmt.Sprintf("%s=%s", k, v))
			} else {
				ret.EnvVars = append(ret.EnvVars, v)
			}
		}
	}

	if dc.Image != "" {
		ret.Image = dc.Image
	}

	if dc.WorkingDir != "" {
		ret.WorkingDir = dc.WorkingDir
	}

	_, err := c.InspectImage(ret.Image)
	if err != nil {
		return nil, errors.Wrapf(err, "Image %s does not exist", ret.Image)
	}

	return ret, nil
}

func (c *Container) AddRequiredMount(binds ...string) error {
	for _, v := range binds {
		v = strings.Replace(v, "$HOST_HOME", os.Getenv("HOME"), -1)
		parts := strings.Split(v, ":")

		_, err := os.Stat(parts[0])
		if err != nil {
			return errors.Wrapf(err, "Invalid path specified: %s", v)
		}

		c.Logger.Debugf("Adding required mount %s", v)
		c.Binds = append(c.Binds, v)
	}

	return nil
}

func (c *Container) AddOptionalMount(binds ...string) error {
	for _, v := range binds {
		v = strings.Replace(v, "$HOST_HOME", os.Getenv("HOME"), -1)
		parts := strings.Split(v, ":")

		_, err := os.Stat(parts[0])
		if err != nil {
			c.Logger.Warnf("Invalid optional mount %s", v)
			continue
		}

		c.Logger.Debugf("Adding optional mount %s", v)
		c.Binds = append(c.Binds, v)
	}

	return nil
}

func (c *Container) Create() error {

	// Create container
	ctr, err := c.Client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Env:          c.EnvVars,
			Image:        c.Image,
			OpenStdin:    true,
			StdinOnce:    true,
			Tty:          true,
			WorkingDir:   c.WorkingDir,
		},
		HostConfig: &docker.HostConfig{
			AutoRemove:      true,
			Binds:           c.Binds,
			PublishAllPorts: true,
			VolumeDriver:    "bind",
		},
	})

	// if something bombs during creation, return it
	if err != nil {
		return err
	}

	c.Container = ctr
	logger.Debugf("Created Container: %s", c.Container.ID)

	return nil
}

func (c *Container) Start() {

	// Puke several log entries out in debug mode so we know how the container is starting up
	c.Logger.Debugf("Starting Console for Container: %s", c.Container.ID)
	c.Logger.Debugf("Num Env Vars: %d", len(c.EnvVars))
	for k, v := range c.EnvVars {
		c.Logger.Debugf("EnvVar %d: %s", k, v)
	}
	c.Logger.Debugf("Image: %s", c.Image)
	c.Logger.Debugf("Num Mounts: %d", len(c.Binds))
	for k, v := range c.Binds {
		c.Logger.Debugf("Bind %d: %s", k, v)
	}
	c.Logger.Debugf("Working Dir: %s", c.WorkingDir)

	go c.start()
	go c.watchSignals()
	go c.watchRunning()
}

func (c *Container) start() {
	// Fire up the console
	if err := dockerpty.Start(c.Client, c.Container, &docker.HostConfig{}); err != nil {
		c.Control.Error(err)
	}
}

func (c *Container) watchSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		c.Logger.Debugf("Got Signal: %s", sig.String())
		c.Control.Quit()
	}()
}

func (c *Container) watchRunning() {
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			container, err := c.Client.InspectContainer(c.Container.ID)
			if err != nil {
				c.Control.Error(err)
				return
			}

			if container == nil {
				c.Control.Error(fmt.Errorf("error getting state for container %s", c.Container.ID))
				return
			}

			if container.State.Running == false {
				c.Control.Error(fmt.Errorf("container %s is not running", c.Container.ID))
				return
			}
		case <-c.Control.WaitForQuitChan():
			c.Logger.Debugf("Caught quit in watchRunning")
			return
		}
	}
}

func (c *Container) Wait() error {
	var returnErr error
	select {
	case err := <-c.Control.WaitForErrorChan():
		c.Stop()
		returnErr = err
	case <-c.Control.WaitForQuitChan():
		c.Stop()
	}

	c.Control.Stop()
	return returnErr
}

func (c *Container) Stop() {
	logger.Debugf("Stopping container: %s", c.Container.ID)

	if err := c.Client.StopContainer(c.Container.ID, 10); err != nil &&

		// Only show errors not pertaining to containers not already in a stopped state
		!strings.Contains(err.Error(), "Container not running") {

		logger.WithError(err).Errorf("Could not stop container: %s", c.Container.ID)
		return
	}

	logger.Debugf("Container %s stopped", c.Container.ID)
}

func (c *Container) populateProjectPaths() error {
	// Find the real path for this binary
	path, err := os.Executable()
	if err != nil {
		return err
	}
	c.Logger.Debugf("Relative Executable Path: %s", path)

	// If this is a symlink, find the real path
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}
	c.Logger.Debugf("Absolute Executable Path: %s", path)

	// Get the directory for the fully resolved path
	path = filepath.Dir(path)
	c.Logger.Debugf("Executable Directory: %s", path)

	// Get the absolute path of the parent directory
	hostProjectDir, err := filepath.Abs(fmt.Sprintf("%s/../", path))
	if err != nil {
		return err
	}

	c.EnvVars = append(c.EnvVars, fmt.Sprintf("HOST_PROJECT_DIR=%s", hostProjectDir))

	// Generate a dir for this project shared into the container
	containerProjectDir := strings.Replace(hostProjectDir, os.Getenv("HOME"), "/home/docker/share", -1)
	c.EnvVars = append(c.EnvVars, fmt.Sprintf("CONTAINER_PROJECT_DIR=%s", containerProjectDir))

	// Log some information
	c.Logger.Debugf("Host Project Dir: %s", hostProjectDir)
	c.Logger.Debugf("Container Project Dir: %s", containerProjectDir)

	return nil
}
