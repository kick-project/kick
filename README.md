[![Github Actions](https://github.com/kick-project/kick/workflows/Go/badge.svg?branch=master)](https://github.com/kick-project/kick/actions) [![Go Report Card](https://goreportcard.com/badge/kick-project/kick)](https://goreportcard.com/report/kick-project/kick)  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/kick-project/kick/blob/master/LICENSE)

# Kick - A work in progress

Kick is a cli tool to start a project using templates under version control
or from a path on local disk.

Its features include

* A way to add templates using a URL or path

```bash
kick install gem git@github.com:kick-project/template-gem.git    # Install a gem template
kick install go git@github.com:kick-project/template-go.git      # Install a go template
kick install pypi git@github.com:kick-project/template-pypi.git  # Install a pypi template
kick install custom_handle1 git@example.com:your/git/project.git # Install a custom template for git
kick install custom_handle2 /path/to/template_directory          # Install a custom template from disk
```

* A cli to start projects 

```bash
kick start gem ~/workspace/mygemproject
kick start go ~/workspace/mygoproject
kick start pypi ~/workspace/mypypiproject
```

* Simple templates using environment variables or dotenv variables from
  `${HOME}/.env`.
  
  Dotenv files provide a way to create "environment" variables without polluting the users environment variables as they are only loaded when a program sourcing the `.env` file runs. Note that `.env` variables will _not_ override existing environment variables.
  Using the sample variables below, a template file is easily populated.
  - `${PROJECT_NAME}`: The directory name parsed from `<path>` when the command `kick start <handle> <path>` is called. 

```text
cat > Makefile <<EOF
# kick:render <- Tell "kick start ..." to render this file as a template. Stripped from the generated file.

_build:
  go build ./cmd/${PROJECT_NAME}

EOF
```

* Template directory paths

```bash
mkdir -p pypi/src/\${PROJECT_NAME}/
touch pypi/src/\${PROJECT_NAME}/\${PROJECT_NAME}.py
```

## Benefits of using Kick

Kick can supercharge the creation of a project to include "starter files"
that will work with the CI of choice or add those additional supporting files
to speed up development.

Starting a project from scratch can be time consuming, there are a few CLI
tools that help (E.G. `go mod init ...` Go project, `rails new` Ruby on
Rails, `pip-init` Pypi project) but they tend to be specific to a type of
project and may not include all the bells and whistles. Consider the
following additions that a project may include...

* CI Integration
  - Github Workflows: `.github/workflows/...`
  - Gitlab CI: `.gitlab-ci.yml` 
* Ignore files
  - Git ignore: `.gitignore`
  - Docker ignore: `.dockerignore`
* Testing tools & libraries
* Linters
* Editor config
  - [Editorconfig](https://editorconfig.org/), project level indentation: `.editorconfig`
* Task automation
  - Make: `Makefile`
  - Python [Invoke](http://www.pyinvoke.org/): `tasks.py`
  - Ruby [Rake](https://github.com/ruby/rake): `Rakefile`

This can all be added to a template which can be called from the command line.

## Getting started

### Install Software

*MacOS*
```bash
brew install kick 
```

*RPMs RHEL/CentOS*
```bash
yum install -y http://github.com/kick-project/kick/archives/kick-1.0.0.rpm
```

*Debs Debian/Ubuntu*
```bash
curl -sLO https://github.com/kick-project/kick/archives/kick-1.0.0.deb
sudo dpkg -i kick-1.0.0.deb
```

### Create variables

Kick uses environment variables or variables defined in `~/.env` to populate
templates. All environment variables and variables defined in `~/.env` are
passed to the templates.

Using an editor add variables to `~/.env`

```dotenv
# ~/.env
AUTHOR=First Last
EMAIL=email@address.com
```

### Create template 

Create a project that will be used as an example to generate go projects.
Templates are rendered using a go library that emulates the function of GNUs
envsubst command.

```bash
mkdir template-go
```

Add an AUTHORS file with ${AUTHOR} and ${EMAIL}

`template-go/AUTHORS`
```yaml
# kick:render <--- This modeline tells.kick to render file as a template. Line is stripped out from output file.
${AUTHOR} ${EMAIL}
```

Add a README.md

`template-go/README.md`
```markdown
# kick:render
# ${PROJECT_NAME}
```

Add the main function

```bash
mkdir -p template-go/cmd/\${PROJECT_NAME}
touch template-go/cmd/\${PROJECT_NAME}/main.go 
```

`template-go/cmd/\${PROJECT_NAME}/main.go`
```go
// kick:render
package main

import "fmt"

func main() {
    fmt.Println("Project ${PROJECT_NAME}")
}
```

### Install the template

```bash
kick install go path/to/template-go
```

### Use the template

```bash
kick start go path/to/project
```

### Upload template to a git repository

```bash
cd template-go
git init
git add .
git commit -m "first commit"
git push --set-upstream git@github.com:owner/template-go.git master
```

### Install the remote template

Remove the installed `go` handle

```bash
kick remove go
```

Install the template using a valid git URL

```bash
kick install go git@github.com:owner/template-go.git
```

### Use the remote template

```bash
kick start go path/to/project
```