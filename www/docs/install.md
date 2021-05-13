# Installation

# MacOS

## homebrew
```bash
brew install kick
```

# Linux

## homebrew
```bash
brew install kick
```

## dpkg
```bash
wget https://github.com/kick-project/kick/releases/download/v1.0.0/kick_1.0.0_amd64.deb
dpkg -i ./kick_1.0.0_amd64.deb
```

## rpm
```bash
rpm -ivh https://github.com/kick-project/kick/releases/download/v1.0.0/kick-1.0.0.x86_64.rpm
```

# Go CLI

## go install

Requires go 1.16 or later
```bash
go install github.com/kick-project/kick/cmd/kick@latest
```

## go get

For version 1.15.x 
```bash
go get -u github.com/kick-project/kick/cmd/kick
```
