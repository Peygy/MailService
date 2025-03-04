BACKDIR = ./back
FRONTDIR = ./front

all: main

.PHONY: main
main:
	$(MAKE) -j 2 back front

.PHONY: back
back:
	$(MAKE) -C $(BACKDIR) build

.PHONY: front
front:
	$(MAKE) -C $(FRONTDIR)