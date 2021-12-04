# Installation

# MacOS

## homebrew
```bash
wget https://github.com/kick-project/kick/releases/download/v1.1.0/kick.rb
brew install kick.rb
```

# Linux

## homebrew
```bash
wget https://github.com/kick-project/kick/releases/download/v1.1.0/kick.rb
brew install kick.rb
```

## dpkg
```bash
wget https://github.com/kick-project/kick/releases/download/v1.1.0/kick_1.1.0_amd64.deb
dpkg -i ./kick_1.0.0_amd64.deb
```

## rpm
```bash
rpm -ivh https://github.com/kick-project/kick/releases/download/v1.1.0/kick-1.1.0.x86_64.rpm
```

# Go CLI

## go install

Requires go 1.16 or later
```bash
go install github.com/kick-project/kick/cmd/kick@v1.1.0
```

## go get

For version 1.15.x
```bash
go get -u github.com/kick-project/kick/cmd/kick
```
