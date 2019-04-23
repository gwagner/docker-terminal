package main

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var CONFIG DockerTerminalConfig
var logger logrus.FieldLogger

func main() {
	if err := run(); err != nil {
		if logger != nil {
			logger.WithError(err).Fatal("Fatal")
		} else {
			fmt.Println(err)
		}
		os.Exit(2)
		return
	}
}

func run() error {
	if err := setupLogging(); err != nil {
		return errors.Wrapf(err, "Error setting up logging")
	}

	if err := parseConfig(); err != nil {
		return errors.Wrapf(err, "Error parsing configs")
	}

	// Get a new container struct
	ctr, err := NewContainer(logger, nil, CONFIG)
	if err != nil {
		return errors.Wrapf(err, "Error creating new container")
	}

	if err := ctr.Create(); err != nil {
		return errors.Wrapf(err, "Error creating new container")
	}

	ctr.Start()
	if err := ctr.Wait(); err != nil {
		return errors.Wrapf(err, "Error running container %s", ctr.Container.ID)
	}

	return nil
}

func parseConfig() error {
	defaultPath := true
	path := fmt.Sprintf("%s/terminal-scripts/dt-config.yaml", os.Getenv("HOME"))
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

func setupLogging() error {
	f, err := os.OpenFile(fmt.Sprintf("%s/dt.log", os.Getenv("HOME")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	logger = logrus.New()
	logger.(*logrus.Logger).SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	logger.(*logrus.Logger).SetOutput(f)
	logger.(*logrus.Logger).SetLevel(logrus.TraceLevel)

	return nil
}
