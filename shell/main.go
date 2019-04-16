package main

import (
	"fmt"
	"github.com/blinkist/go-dockerpty"
	docker "github.com/fsouza/go-dockerclient"
	"os"
)

var OS_HOME = os.Getenv("HOME")

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

	binds := []string{
		"/var/run/docker.sock:/var/run/docker.sock",
		fmt.Sprintf("%s:/home/docker/share/", OS_HOME),
		fmt.Sprintf("%s/terminal-scripts/.gitconfig:/home/docker/.gitconfig", OS_HOME),
		fmt.Sprintf("%s/terminal-scripts/.motd:/home/docker/.motd", OS_HOME),
		fmt.Sprintf("%s/terminal-scripts/.nanorc:/home/docker/.nanorc", OS_HOME),
		fmt.Sprintf("%s/terminal-scripts/.zlogout:/home/docker/.zlogout", OS_HOME),
		fmt.Sprintf("%s/terminal-scripts/.zshrc:/home/docker/.zshrc", OS_HOME),
		fmt.Sprintf("%s/.zsh_history:/home/docker/.zsh_history", OS_HOME),
	}

	if _, err := os.Stat(fmt.Sprintf("%s/.ssh", OS_HOME)); !os.IsNotExist(err) {
		binds = append(binds, fmt.Sprintf("%s/.ssh:/home/docker/.ssh", OS_HOME))
	}

	if _, err := os.Stat(fmt.Sprintf("%s/go", OS_HOME)); !os.IsNotExist(err) {
		binds = append(binds, fmt.Sprintf("%s/go:/home/docker/go", OS_HOME))
	}

	// Create container
	ctr, err := cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Image:        "terminal:latest",
			OpenStdin:    true,
			StdinOnce:    true,
			Tty:          true,
			WorkingDir:   "/home/docker",
		},
		HostConfig: &docker.HostConfig{
			VolumeDriver:    "bind",
			PublishAllPorts: true,
			AutoRemove:      true,
			Binds:           binds,
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

/* func run() error {
	ctx = context.Background()
	cli, err = client.NewEnvClient()
	if err != nil {
		return err
	}

	if err = container_create(ctx, cli); err != nil {
		return err
	}

	if err = container_attach(ctx, cli); err != nil {
		return err
	}
	defer containerConn.Close()

	sigc := container_forward_all_signals(ctx, cli)
	defer signal.StopCatch(sigc)

	if err = container_io(); err != nil {
		return err
	}
	defer terminal.Restore(fd, terminalState)

	ok, err := cli.ContainerWait(ctx, containerId, container.WaitConditionRemoved)
	select{
	case e := <- err:
		return e
	case <- ok:
	}

	containerConn.Close()

	return nil
}*/
