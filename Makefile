
include $(HOME)/template/mk/go.v1.mk

SCURRENT = $(shell pwd)/
DCURRENT = /mnt/hgfs/work/go/proj/wingui

.PHONY: put

put:
	$(RSYNC_CMD) $(SCURRENT) $(DCURRENT)

