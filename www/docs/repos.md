# Repositories

Kick template repositories are git projects with references to one ore more
templates. Kick sub commands such as `kick search` and `kick install` will search and
install from a kick repository.

# Making a Repository

## Create a template

Download the example template

```bash
mkdir -p ~/templates
cd ~/templates
wget http://github.com/kick-project/kick/releases/latest/website-template.tar.gz -O - | tar -zxvf -
```

To initialize a template for kick repository use, run `kick init template <name>`.
This command creates a `repo.yml`.
```bash
cd website-template
kick init template website
# <STDOUT>
# generated repo.yml
```

Initialize git, commit and push changes to a upstream repository
```bash
cd website-template
git init
git add .
git commit -m 'Initial commit'
git remote add origin git@github.com/example/website-template.git
git push --set-upstream origin master
```

## Create Repository

```bash
mkdir ~/repos
cd ~/repos

mkdir myrepo
cd myrepo
kick init repo myrepo
# <STDOUT>
# generated repo.yml
```

Modify repo.yml to include the new repository
```yaml
# repo.yml
name: myrepo
description: Repository myrepo
templates:
- git@github.com/example/website-template.git
```

## Build repository

Build repository by running the `kick repo build` subcommand. This will clone the repositories defined under the templates section in the yaml file and
copy the contents of `.kick.yml` into a subdirectory `templates/`. Note that this step can also be performed manually by creating the files under `templates/`.

```bash
kick repo build
```