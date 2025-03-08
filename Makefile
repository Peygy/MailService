BACKDIR = ./back
FRONTDIR = ./front

#mkwin
all: main

.PHONY: main
main:
	$(MAKE) -j 2 back front pass=$(pass)

.PHONY: back
back:
	$(MAKE) -C $(BACKDIR) build pass=$(pass)

.PHONY: front
front:
	$(MAKE) -C $(FRONTDIR)