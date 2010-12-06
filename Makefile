include $(GOROOT)/src/Make.inc

GOFMT=gofmt -s -spaces=true -tabindent=false -tabwidth=4

TARG=simpleconfig
GOFILES=\
	simpleconfig.go\

include $(GOROOT)/src/Make.pkg

format:
	$(GOFMT) -w simpleconfig.go
	$(GOFMT) -w simpleconfig_test.go

