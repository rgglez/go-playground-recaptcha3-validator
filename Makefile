.PHONY: current next patch minor major release push

# Current tag (v0.0.0 if none yet)
CURRENT := $(shell git describe --tags --abbrev=0 2>/dev/null || echo v0.0.0)
PART ?= patch

# Next tag: bump PART (patch|minor|major) of CURRENT
NEXT = $(shell echo $(CURRENT) | awk -F. -v p=$(PART) '{ \
	major=substr($$1,2)+0; minor=$$2+0; patch=$$3+0; \
	if (p=="major") { major++; minor=0; patch=0 } \
	else if (p=="minor") { minor++; patch=0 } \
	else { patch++ } \
	printf "v%d.%d.%d", major, minor, patch }')

current:
	@echo $(CURRENT)

next:
	@echo $(NEXT)

patch minor major:
	@$(MAKE) --no-print-directory release PART=$@

# Tag and push to origin
release:
	@test -z "$$(git status --porcelain)" || { echo "working tree dirty"; exit 1; }
	git tag -a $(NEXT) -m "Release $(NEXT)"
	git push origin $(NEXT)

# Push an existing tag: make push TAG=v1.2.3
push:
	@test -n "$(TAG)" || { echo "usage: make push TAG=vX.Y.Z"; exit 1; }
	git push origin $(TAG)
