BUILDPATH = $(CURDIR)
GO = $(shell with go)
GOINSTALL = $(GO) install
GOCLEAN = $(GO) clean

EXENAME = [ -f ~/.aliases ] && source ~/.aliases
			ali() {
				$GOPATH/github.com/adriacidre/ali/ali $@
				source ~/.aliases
				if [ [ "$1" == "rm" ] ] then
					unalias $2
				fi
			}

export GOPATH = $(CURDIR)

myname:
	@echo "I am a makefile"

install:
	$(GOINSTALL) $(EXENAME)

#això es el que havia fet abans pero tampoc hem funciona

#install:\
	[ -f ~/.aliases ] && source ~/.aliases\
	ali() {\
		$GOPATH/github.com/adriacidre/ali/ali $@\
		source ~/.aliases\
		@if [ [ "$1" == "rm" ] ]; then\
			unalias $2;\
		fi\
	}\

#havia fet rel mateix amb build i supu