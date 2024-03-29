# Build

Build, test and develop kick

## Requirements

* Go 1.16 or later.
* Sqlite3 3.6.19 or later. Required for FOREIGN KEY support.
* GCC compiler.
* [Vagrant](https://www.vagrantup.com/) for test cross device moving of temporary files.

## Development dependencies

Install development dependencies.

```bash
make depsdev
```

## Dependencies

Install program dependencies.

```bash
make deps
```

## Build & Test

```bash
# Build
make build
```

```bash
# Test
make test
```

```bash
# View test results in a browser
make report
```

## Vagrant

To use Vagrant to test cross device moving of temporary files.

```bash
# Bring vagrant up and ssh
make ssh

# Install development dependencies in vagrant
make depsdev

# Install dependencies
make deps

# Test
make test
```

## Help

To get help, run the following after installing development dependencies.
Requires `$PATH` contains `$GOPATH/bin`

```bash
make help
```

```bash
Usage: make <target>

### HELP
  help - Print Help
### DEVELOPMENT
  build   - Build binary
  install - Install to $(GOPATH)/bin
  clean   - Reset project to original state
  test    - Test
  package - Create an RPM, Deb, Homebrew package
  release - Trigger a release by creating a tag and pushing to the upstream repository
  deps    - Install build dependencies
  depsdev - Install development dependencies
  report  - Open reports generated by "make test" in a browser
### VERSION INCREMENT
  bumpmajor - Increment VERSION file ${major}.0.0 - major bump
  bumpminor - Increment VERSION file 0.${minor}.0 - minor bump
  bumppatch - Increment VERSION file 0.0.${patch} - patch bump
### DOCUMENTATION
  builddocs  - Build documents
  deploydocs - Deploy documentaiton to github
  docserver  - Start document server
### VAGRANT
  up      - Start vagrant
  ssh     - Run vagrant ssh and cd to shared directory
  halt    - Run vagrant halt
  destroy - Run vagrant destroy
```