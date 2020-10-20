VERSION=$(shell versionbump list .kickmod.yml)

version:
	@echo $(VERSION)

bumpmajor:
	git fetch --tags
	versionbump --checktags major .kickmod.yml

bumpminor:
	git fetch --tags
	versionbump --checktags minor .kickmod.yml

bumppatch:
	git fetch --tags
	versionbump --checktags patch .kickmod.yml

release:
	git fetch --tags
	git tag v$(VERSION)