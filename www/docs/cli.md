# Use & manage templates
## kick start

```bash
generate project scaffolding

Usage:
    kick start <handle> <project>
    kick start (-l|--long)

Options:
    -h --help     print help
    --long        list templates in long format
    <handle>      template handle
    <project>     project path
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
Buid/list/inform on repositories WIP

Usage:
    kick repo build
    kick repo list
    kick repo info <repo>

Options:
    -h --help    print help
    repo         repo subcommand
    build        build repo by downloading the URLS defined in repo.yml and creating the files templates/*.yml
    list         list repositories
    info         repository and/or template information
    <repo>       name of repository
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
