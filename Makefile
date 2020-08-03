#!/bin/bash

install:
	@go install
	if [[ ! -e ~/.aliases ]]; then
		touch ~/.aliases
	fi
	echo "[ -f ~/.aliases ] && source ~/.aliases"