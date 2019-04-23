package main

import "sync"

type DockerTerminalConfig struct {
	EnvVars        map[string]string `yaml:"env_vars"`
	Image          string            `yaml:"image"`
	RequiredMounts []string          `yaml:"required_mounts"`
	OptionalMounts []string          `yaml:"optional_mounts"`
	WorkingDir     string            `yaml:"working_dir"`
}

type FanoutControlChan struct {
	sync.Mutex
	error          chan error
	errorListeners []chan error
	quit           chan bool
	quitListeners  []chan bool
}

func NewFanoutControlChan() *FanoutControlChan {
	f := &FanoutControlChan{
		error: make(chan error),
		quit:  make(chan bool),
	}

	f.Start()

	return f
}

func (f FanoutControlChan) Start() {
	go func() {
		select {
		case q := <-f.quit:
			for _, v := range f.quitListeners {
				v <- q
			}

		case e := <-f.error:
			for _, v := range f.errorListeners {
				v <- e
			}
		}

	}()
}

func (f FanoutControlChan) Stop() {
	f.Lock()
	defer f.Unlock()

	// Since we are stopping, nobody else can listen
	for _, v := range f.errorListeners {
		close(v)
	}
	for _, v := range f.quitListeners {
		close(v)
	}
}

func (f *FanoutControlChan) WaitForError() error {
	l := make(chan error)
	f.errorListeners = append(f.errorListeners, l)

	return <-l
}

func (f *FanoutControlChan) WaitForErrorChan() chan error {
	l := make(chan error)
	f.errorListeners = append(f.errorListeners, l)

	return l
}

func (f *FanoutControlChan) WaitForQuit() {
	l := make(chan bool)
	f.quitListeners = append(f.quitListeners, l)

	<-l

	return
}

func (f *FanoutControlChan) WaitForQuitChan() chan bool {
	l := make(chan bool)
	f.quitListeners = append(f.quitListeners, l)

	return l
}

func (f FanoutControlChan) Error(err error) {
	f.error <- err
	f.Stop()
}

func (f FanoutControlChan) Quit() {
	f.quit <- true
	f.Stop()
}
