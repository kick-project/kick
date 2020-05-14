VERSION=$(shell versionbump list .prjmod.yml)

version:
	@echo $(VERSION)

bumpmajor:
	git fetch --tags
	versionbump --checktags major .prjmod.yml

bumpminor:
	git fetch --tags
	versionbump --checktags minor .prjmod.yml

bumppatch:
	git fetch --tags
	versionbump --checktags patch .prjmod.yml

release:
	git fetch --tags
	git tag v$(VERSION)