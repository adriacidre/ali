#!/bin/bash

install:
	go build
	cp -n _tpl/.aliases ~/.aliases
	sed -i 's/{{PATH}}/`pwd`/g' ~/.aliases
	echo "[ -f ~/.aliases ] && source ~/.aliases" >> ~/.zshrc
