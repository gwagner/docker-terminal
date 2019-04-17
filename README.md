# docker-terminal

Use docker as your normal every day terminal.  Take your specific config from machine to machine without needing 
to do any of the symbolic linking magic most people do to ensure settings are move.  Startup a terminal in a container, when done exit.  

Forget to exit?  Your container will be cleaned up in 6 hours.  Need your terminal back, no worries, you can always `docker attach`
back to the running container.

# Is this needed?

Probably not, just a fun experiment!!

# How do i use this

1. `git clone git@github.com:gwagner/docker-terminal.git`
1. `cd [PATH TO CHECKOUT]`
1. `./build.sh`
1. `./install.sh`
1. `[OPEN NEW TERMINAL]`

## Optional

* Import iTerm2 profile for fancy colors and graphs

# Future Plans

1. Without needing to change your default shell, get back to your host shell through this PTY
1. Set a config for what container to turn on.  Right now only the current Dockerfile.ubuntu is supported, but i should show Centos and Gentoo some love
1. Container flavors.  The idea would be building pre-dockerfiled containers which fit a specific purpose to reduce size and increase efficiency
    * PHP Developer Container
    * Golang Developer Container
    * Devops Container
    * Other Kinds?
1. Test this on a few linux flavors.  The ultimate idea of this is to provide a familiar shell to users without any need to keep multiple bastion hosts up to date.
1. Work with remote docker on remote boxes?
1. Maybe test on Windows
    * I am not a windows guy, so this is probably pretty low on the list