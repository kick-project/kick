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
  `${HOME}/.env`. Dotenv files provide a way to create "environment" variables without polluting the Users environment variables as they are only loaded when a program sourcing the `.env` file runs. Note that `.env` variables will _not_ override existing environment variables.
  Using the sample variables below, a template file is easily populated.
  - `${USER}`: From the users environment variables
  - `${PROJECT_NAME}`: The directory name parsed from `<path>` when the command `kick start <handle> <path>` is called. 

```text
cat > Makefile <<EOF
# kick:render <- Tell "kick start ..." to render this file as a template. Stripped from the generated file.
DOTENV:=dotenv .env

AUTHOR="${USER}"

_build:
  go build cmd/${PROJECT_NAME}

# make wrapper - Execute any target prefixed with a underscore if the target is
# not explicitly defined in the Makefile. EG 'make build' will result in the
# execution of 'make _build'.
%:
	@egrep -q '^_$@:' Makefile && $(DOTENV) $(MAKE) _$
EOF
```

# Why?

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
* Custom test tools
  - Python test tools: `pytest`
* Editor config
  - [Editorconfig](https://editorconfig.org/). Project level indentation: `.editorconfig`
* Task automation
  - Make: `Makefile`
  - Python [Invoke](http://www.pyinvoke.org/): `tasks.py`
  - Ruby [Rake](https://github.com/ruby/rake): `Rakefile`

This can all be added to a template which can be called from the command line.

# Quick start

Add a set of variables that will be used in your project. Variables are
stored in `~/.env`. All the variables defined here are passed to templates
along with all environment variables as `${VARIABLE}`

`~/.env`
```dotenv
AUTHOR=First Last <first.last@somedomain.com>
```

Create a project that will be used as an example to generate go projects.
Templates are rendered using a go library that emulates the function of GNUs
envsubst command.
```bash
mkdir -p ~/kicks/kickgo
```

`~/kicks/kickgo/AUTHORS`
```yaml
# kick:render <--- This modeline tells.kick to render file as a template. Line is stripped out from output file.
${AUTHOR}
```

`~/kicks/kickgo/README.md`
```markdown
# kick:render
# ${PROJECT_NAME}
```

`~/kicks/kickgo/.gitignore`
```.gitignore
# Vim
.*.swp

# Mac
.DS_Store

/vendor/
```

Add a binary
```bash
mkdir -p ~/kicks/kickgo/cmd/\${PROJECT_NAME}
touch ~/kicks/kickgo/cmd/\${PROJECT_NAME}/main.go 
```

`~/kicks/kickgo/cmd/\${PROJECT_NAME}/main.go`
```go
// kick:render
package main

import "fmt"

func main() {
    fmt.Println("Project ${PROJECT_NAME}")
}
```

Add the following file `~/.kick/templates.yml`
```yaml
- name: goproject
  url: ~/kicks/kickgo
```

Create the project using go
```bash
kick start goproject ~/mynewproject
```

# Git templates

```bash
cd ~/kicks/kickgo
git init
git add .
git commit -m "first commit"
git push --set-upstream git@github.com/owner/kickgo.git master
```

Modify `~/.kick/templates.yml`
```yaml
- name: goproject
  url: http://github.com/owner/kickgo.git
```

Start a new project with the recently checked in boilerplate
```bash
kick start goproject ~/myproject
```
