# Copyright 2009 Artem Titoulenko and the Go Authors. All rights reserved.
# blah blah blah
#

#this makefile is a bit diff since it copies stones into the $GOBIN directory
#so it can freely be used in the command prompt

all:%.make install

include $(GOROOT)/src/Make.$(GOARCH)

TARG=stone
GOFILES=\
	main.go\

include $(GOROOT)/src/Make.pkg

%.make:
	$(GOBIN)/*g -o$(TARG).8 $(GOFILES)
	$(GOBIN)/*l -o$(TARG) $(TARG).8

clean:
	rm -f *.$O $(TARG) *.8 8.out a.out
	rm -r -f _obj

install: $(TARG)
	cp $(TARG) $(GOBIN)/$(TARG)
	
nuke: clean
	rm $(GOBIN)/$(TARG)