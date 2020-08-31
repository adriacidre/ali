#!/bin/bash

install:
	go install
	touch ~/.aliases
	echo "[ -f ~/.aliases ] && source ~/.aliases" >> ~/.zshrc
