#!/bin/bash

DOCKER_HOME_DIR="/home/gwagner"

CMD="/usr/local/bin/docker run -it --rm" 
CMD+=" -w $DOCKER_HOME_DIR"
CMD+=" -v /var/run/docker.sock:/var/run/docker.sock"
CMD+=" --mount type=bind,source=\"$HOME/\",target=\"$DOCKER_HOME_DIR/share\""
CMD+=" --mount type=bind,source=\"$HOME/terminal-scripts/.gitconfig\",target=\"$DOCKER_HOME_DIR/.gitconfig\""
CMD+=" --mount type=bind,source=\"$HOME/terminal-scripts/.nanorc\",target=\"$DOCKER_HOME_DIR/.nanorc\""
CMD+=" --mount type=bind,source=\"$HOME/terminal-scripts/.zlogout\",target=\"$DOCKER_HOME_DIR/.zlogout\""
CMD+=" --mount type=bind,source=\"$HOME/terminal-scripts/.zshrc\",target=\"$DOCKER_HOME_DIR/.zshrc\""
CMD+=" --mount type=bind,source=\"$HOME/.zsh_history\",target=\"$DOCKER_HOME_DIR/.zsh_history\""

if [ -d ~/.ssh ]; then
    CMD+=" --mount type=bind,source=\"$HOME/.ssh\",target=\"$DOCKER_HOME_DIR/.ssh\""
fi

CMD+=" terminal:latest"

eval $CMD