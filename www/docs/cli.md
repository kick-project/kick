# Use & manage templates
## kick start

```bash
generate project scaffolding

Usage:
    kick start <handle> <project>

Options:
    -h --help     print help
    <handle>      template handle
    <project>     project path
```

## kick list

```bash
list handles/templates

Usage:
    kick list [-l]

Options:
    -h --help     print help
    -l            print long output
```

## kick install

```bash
Install template

Usage:
    kick install <handle> <location>

Options:
    -h --help        print help
    <handle>         name to use when creating new projects
    <location>       template name, URL or location of template
```

## kick remove

```bash
Remove an installed template

Usage:
    kick remove <handle>

Options:
    -h --help        print help
    <handle>         handle to remove
```

## kick search

```bash
search for templates using a keyword

Usage:
    kick search [-l] [<term>]

Options:
    -h --help  print help
    -l         long output
    <term>     search term
```

# Management commands

## kick setup

```bash
initialize configuration

Usage:
    kick setup

Options:
    -h --help     print help
```

# Repository management

## kick repo

```bash
Build a repo from repo.yml

Usage:
    kick repo build

Options:
    -h --help    print help
    repo         repo subcommand
    build        build repo by downloading the URLS defined in repo.yml and creating the files templates/*.yml
```

## kick init

```bash
Create a repo or template

Usage:
    kick init repo <name> [<path>]
    kick init template <name> [<path>]

Options:
    -h --help    print help
    repo         create repository       
    template     create a template
    <name>       template or repo name
    <path>       directory path. if not set creates files in working directory
```
