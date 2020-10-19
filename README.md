[![Github Actions](https://github.com/kick-project/kick/workflows/Go/badge.svg?branch=master)](https://github.com/kick-project/kick/actions) [![Go Report Card](https://goreportcard.com/badge/kick-project/kick)](https://goreportcard.com/report/kick-project/kick)  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/kick-project/kick/blob/master/LICENSE)

# Project Start -.kick

Project Start .kick` is a cli tool to start a project using template 
boilerplates.

# Quickstart

Add a set of variables that will be used in your project.
Variables are stored in `~/.env`. All the variables defined here
are passed to templates along with all environment variables as
`.Env.$VARIABLE`

`~/.env`
```dotenv
AUTHOR=First Last <first.last@somedomain.com>
```

Create a project that will be used as an example to generate go projects.
Templates are rendered using Go's text/template
```bash
mkdir -p ~/prjs/prjgo
```

`~/prjs/prjgo/AUTHORS`
```yaml
# prj:render <--- This modeline tells.kick to render file as a template. Line is stripped out from output file.
{{.Env.AUTHOR}}
```

`~/prjs/prjgo/README.md`
```markdown
# prj:render
# {{.Project.NAME}}
```

`~/prjs/prjgo/.gitignore`
```.gitignore
# Vim
.*.swp

# Mac
.DS_Store

/vendor/
```

Add a binary
```bash
mkdir -p ~/prjs/prjgo/cmd/\{\{.Project.NAME\}\}
touch ~/prjs/prjgo/cmd/\{\{.Project.NAME\}\}/main.go 
```

`~/prjs/prjgo/cmd/\{\{.Project.NAME\}\}/main.go`
```go
// prj:render
package main

import "fmt"

func main() {
    fmt.Println("Project {{.Project.NAME}}")
}
```

Add the following file `~/.kick.yml`
```yaml
templates:
    - name: goproject
      url: ~/prjs/prjgo
```

Create the project using go
```bash
prjstart start goproject ~/mynewproject
```

# Git templates

```bash
cd ~/prjs/prjgo
git init
git add .
git commit -m "first commit"
git push --set-upstream git@github.com/owner/prjgo.git master
```

Modify `~/.kick.yml`
```yaml
templates:
    - name: goproject
      url: http://github.com/owner/prjgo.git
```

Start a new project with the recently checked in boilerplate
```bash
prjstart start goproject ~/myproject
```

# Variables
Variables are either environment variables, variables defined in `~/.env`
or project variables.

To list available project variables run
```bash
prjstart list --vars
```

