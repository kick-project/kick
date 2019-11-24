[![Go Report Card](https://goreportcard.com/badge/crosseyed/prjstart)](https://goreportcard.com/report/crosseyed/prjstart) [![Build Status](https://travis-ci.org/crosseyed/prjstart.svg?branch=master)](https://travis-ci.org/crosseyed/prjstart) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/crosseyed/prjstart/blob/master/LICENSE)

# Project Start - prjstart

Project Start `prjstart` is a cli tool to start a project using template 
boilerplates.

# Quickstart

Create Hello World project. Templates are rendered using Go's text/template
```bash
mkdir -p ~/prjs/helloworld
cat > ~/prjs/helloworld/project.yaml <<EOF
# prj:render <--- Tell modeline to render file as a template. Line is stripped out from output file.
project: {{.Project.NAME}}
home: {{.Env.HOME}}
EOF

cat > ~/prjs/helloworld/README.md <<EOF
# Hello World 
This file will not be rendered as a template
EOF
```

Simple directory render
```bash
mkdir -p ~/prjs/directoryrender/\{\{.Env.USER\}\}
touch ~/prjs/directoryrender/\{\{.Env.USER\}\}/emptyfile
```

Add the following file `~/.prjstart.yml`
```yaml
templates:
    - name: helloworld
      url: ~/prjs/helloworld
      desc: Hello World
    - name: directoryrender
      url: ~/prjs/directoryrender
      desc: Example Directory render
```
