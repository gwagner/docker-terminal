#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ ! -f bin/dt ]; then
    echo "Please run build.sh first to ensure dt is built"
    exit 2
fi

# On OSX, update maxfiles to be more docker friendly
echo "# Needs root access to set maxfiles (ulimit) on OSX to a reasonable value for docker"
echo
sudo launchctl limit maxfiles 65536 200000

echo "# Symlinking scripts folder to ~/terminal-scripts"
echo
if [ ! -e ~/terminal-scripts ]; 
then
    ln -s $DIR/scripts/ $HOME/terminal-scripts
fi

echo "# Ensuring ~/.zsh_history exists"
echo
touch $HOME/.zsh_history

echo "# Ensure dt is coppied to /usr/local/bin"
echo
if [ ! -e /usr/local/bin/dt ];
then
    ln -s $DIR/bin/dt /usr/local/bin/dt
fi

echo "# Change the default shell for active user to /usr/local/bin/dt"
echo
chsh -s /usr/local/bin/dt $USER