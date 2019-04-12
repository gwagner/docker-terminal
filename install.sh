#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

if [ ! -e ~/terminal-scripts ]; 
then
    ln -s $DIR/scripts/ $HOME/terminal-scripts
fi

touch $HOME/.zsh_history