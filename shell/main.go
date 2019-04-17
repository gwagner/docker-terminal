package main

import (
	"fmt"
	"github.com/blinkist/go-dockerpty"
	"github.com/fsouza/go-dockerclient"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var HOST_HOME = os.Getenv("HOME")
var CONFIG DockerTerminalConfig

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(2)
		return
	}
}

func run() error {
	endpoint := "unix:///var/run/docker.sock"
	cli, _ := docker.NewClient(endpoint)

	// Find the real path for this binary
	path, _ := os.Executable()
	path, _ = filepath.EvalSymlinks(path)
	path = filepath.Dir(path)
	hostProjectDir, _ := filepath.Abs(fmt.Sprintf("%s/../", path))
	containerProjectDir := strings.Replace(hostProjectDir, HOST_HOME, "/home/docker/share", -1)

	if err := parseConfig(); err != nil {
		return err
	}

	binds, err := getBinds()
	if err != nil {
		return err
	}

	// Create container
	ctr, err := cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Env: []string{
				fmt.Sprintf("HOST_HOME=%s", HOST_HOME),
				fmt.Sprintf("HOST_PROJECT_DIR=%s", hostProjectDir),
				fmt.Sprintf("CONTAINER_PROJECT_DIR=%s", containerProjectDir),
				"IS_CONTAINER=true",
			},
			Image:      "terminal:latest",
			OpenStdin:  true,
			StdinOnce:  true,
			Tty:        true,
			WorkingDir: "/home/docker",
		},
		HostConfig: &docker.HostConfig{
			AutoRemove:      true,
			Binds:           binds,
			PublishAllPorts: true,
			VolumeDriver:    "bind",
		},
	})

	if err != nil {
		return err
	}

	// Cleanup when done
	defer func() {
		cli.RemoveContainer(docker.RemoveContainerOptions{
			ID:    ctr.ID,
			Force: true,
		})
	}()

	// Fire up the console
	if err = dockerpty.Start(cli, ctr, &docker.HostConfig{}); err != nil {
		return err
	}

	return nil
}

func parseConfig() error {
	defaultPath := true
	path := fmt.Sprintf("%s/terminal-scripts/dt-config.yaml", HOST_HOME)
	if p, ok := os.LookupEnv("CONFIG"); ok {
		path = p
		defaultPath = false
	}

	_, err := os.Stat(path)
	if e, ok := err.(*os.PathError); ok {
		if defaultPath {
			return nil
		}

		return e
	}

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(dat, &CONFIG); err != nil {
		return err
	}

	return nil
}

func getBinds() ([]string, error) {
	binds := []string{
		"/var/run/docker.sock:/var/run/docker.sock",
		fmt.Sprintf("%s:/home/docker/share/", HOST_HOME),
	}

	mounts, err := CONFIG.GetMounts()
	if err != nil {
		return []string{}, err
	}
	binds = append(binds, mounts...)

	return binds, nil
}
