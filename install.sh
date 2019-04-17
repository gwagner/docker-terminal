#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

# Determine which home to use, inside or outside the container
SHARED_HOME=$HOME
if [ ! -z $IS_CONTAINER ]; then
    SHARED_HOME="/home/docker/share"
    DIR=$CONTAINER_PROJECT_DIR
fi

if [ ! -f $DIR/bin/dt ]; then
    echo
    echo "Please run build.sh first to ensure dt is built"
    exit 2
fi

# Only run outside of a container
if [ -z $IS_CONTAINER ]; then
    # On OSX, update maxfiles to be more docker friendly
    echo
    echo "# Needs root access to set maxfiles (ulimit) on OSX to a reasonable value for docker"
    sudo launchctl limit maxfiles 65536 200000
fi

if [ ! -h $SHARED_HOME/terminal-scripts ];
then
    echo
    echo "# Symlinking scripts folder to ~/terminal-scripts"
    ln -s $HOST_PROJECT_DIR/scripts/ $SHARED_HOME/terminal-scripts
fi

echo
echo "# Ensuring ~/.zsh_history exists"
touch $SHARED_HOME/.zsh_history

# Only run outside of a container
if [ -z $IS_CONTAINER ] && [ ! -e /usr/local/bin/dt ];
then
    echo
    echo "# Ensure dt is coppied to /usr/local/bin"
    ln -s $HOST_PROJECT_DIR/bin/dt /usr/local/bin/dt
fi

echo "# Change the default shell for active user to /usr/local/bin/dt"
echo
chsh -s /usr/local/bin/dt $USER