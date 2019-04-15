#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

# On OSX, update maxfiles to be more docker friendly
sudo launchctl limit maxfiles 65536 200000

if [ ! -e ~/terminal-scripts ]; 
then
    ln -s $DIR/scripts/ $HOME/terminal-scripts
fi

touch $HOME/.zsh_history
