#!/bin/sh

# detach
# ctrl-b + d

# scroll up /down
# ctrl-b + [ , q exit

option="${1}"
case ${option} in
    new) echo new
      tmux new -d -s wrappingToken '/vault/wrappingToken'
      ;;
    kill) echo kill
      tmux kill-session -t wrappingToken
      ;;
    ls) echo ls
      tmux ls
      ;;
    a) echo attach
      tmux a -t wrappingToken
      ;;
    *)  echo 'Unknown!'
      exit 0
    ;;
esac