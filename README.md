[![Github Actions](https://github.com/kick-project/kick/workflows/Go/badge.svg?branch=master)](https://github.com/kick-project/kick/actions) [![Go Report Card](https://goreportcard.com/badge/kick-project/kick)](https://goreportcard.com/report/kick-project/kick)  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/kick-project/kick/blob/master/LICENSE)

# Kick

`kick` is a cli tool to start a project using templates under version control
or stored on local disk.

# Quickstart

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
