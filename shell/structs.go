package main

import (
	"github.com/pkg/errors"
	"os"
	"strings"
)

type DockerTerminalConfig struct {
	RequiredMounts []string `yaml:"required_mounts"`
	OptionalMounts []string `yaml:"optional_mounts"`
}

func (d *DockerTerminalConfig) GetMounts() ([]string, error) {
	required, err := d.GetRequiredMounts()
	if err != nil {
		return []string{}, err
	}

	optional, err := d.GetOptionalMounts()
	if err != nil {
		return []string{}, err
	}

	return append(required, optional...), nil
}

func (d *DockerTerminalConfig) GetRequiredMounts() ([]string, error) {
	var mounts []string
	for _, v := range d.RequiredMounts {
		v = strings.Replace(v, "$HOST_HOME", HOST_HOME, -1)
		parts := strings.Split(v, ":")

		_, err := os.Stat(parts[0])
		if err != nil {
			return []string{}, errors.Wrapf(err, "Invalid path specified: %s", v)
		}

		mounts = append(mounts, v)
	}

	return mounts, nil
}

func (d *DockerTerminalConfig) GetOptionalMounts() ([]string, error) {
	var mounts []string
	for _, v := range d.OptionalMounts {
		v = strings.Replace(v, "$HOST_HOME", HOST_HOME, -1)
		parts := strings.Split(v, ":")

		_, err := os.Stat(parts[0])
		if err != nil {
			continue
		}

		mounts = append(mounts, v)
	}

	return mounts, nil
}
