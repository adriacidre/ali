#!/bin/bash

install:
	@echo "Building binaries"
	go build
	@echo "Replacing paths on your config files"
	@CURRENT_PATH=`pwd`
	@sed -i'' -e "s|{{PATH}}|$CURRENT_PATH|g" _tpl/.aliases
	@echo "Setting up config files"
	@cp -n _tpl/.aliases ~/.aliases || true
	@echo "Setting up ali on your zshrc"
	@echo "[ -f ~/.aliases ] && source ~/.aliases" >> ~/.zshrc
	@echo "\nAli is now ready to be used! type 'ali list' to list your commands"
