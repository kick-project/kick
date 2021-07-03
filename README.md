[![Github Actions](https://github.com/kick-project/kick/workflows/Go/badge.svg?branch=master)](https://github.com/kick-project/kick/actions) [![Go Report Card](https://goreportcard.com/badge/kick-project/kick)](https://goreportcard.com/report/kick-project/kick)  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/kick-project/kick/blob/master/LICENSE)

# Kick

For docs and getting started see https://kick-project.github.io/kick/

## About

Kick is a cli tool to start a project using templates under version control
or from a path on local disk.

Its features include

* A way to add templates using a git remote location, URL or local path.

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

Kick saves time and supercharges the creation of a project by using cookie
cutter "starter files" that will work with the CI, add task automation using
make, Rake, Invoke, add support to create packages such as gems, pypi, rpms,
debs and many more files and tools that are needed in the development
process.

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
