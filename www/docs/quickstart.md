# Quick Start

# Project templates

Project templates are folders that contain other directories and files.

Given the following directory tree, lets create a `Project Template` which contains the
same files and folders.

```
mytemplate
├── cmd
│   └── ${PROJECT_NAME}
│       └── main.go
├── docs
├── .gitignore
├── README.md
└── scripts
    └── ${PROJECT_NAME}.sh
```

Create a folder to hold templates.
```bash
mkdir ~/templates
cd ~/templates
```

Create the Project Template `mytemplate`.
```bash
mkdir mytemplate
mkdir -p 'mytemplate/cmd/${PROJECT_NAME}/'
touch 'mytemplate/cmd/${PROJECT_NAME}/main.go'

# Create git ignore
cat <<'EOF' > mytemplate/.gitignore
.DS_Store # MacOS files
.*.swp    # Vim buffer file
EOF

touch mytemplate/README.md
mkdir -p mytemplate/docs
mkdir -p mytemplate/scripts
touch 'mytemplate/scripts/${PROJECT_NAME}.sh'
```

Verify that the structure matches the layout above.
```bash
tree -an mytemplate
mytemplate
├── cmd
│   └── ${PROJECT_NAME}
│       └── main.go
├── docs
├── .gitignore
├── README.md
└── scripts
    └── ${PROJECT_NAME}.sh
```

Examine the `.gitignore` file in your project.
```bash
cat mytemplate/.gitignore
```

Install the template using `kick` by specifying the path to the `mytemplate`
directory. In the example below we will use the handle of `myhandle`.

```bash
kick install myhandle ~/templates/mytemplate
```

To see the installed template run the `kick start -l` or `kick start --long`.
```bash
kick start --long
+----------+----------+-------------+------------------------------------+
|  HANDLE  | TEMPLATE | DESCRIPTION |              LOCATION              |
+----------+----------+-------------+------------------------------------+
| myhandle | -        | -           | /home/vagrant/templates/mytemplate |
+----------+----------+-------------+------------------------------------+
```

Create a projects directory to create new projects.
```bash
mkdir ~/projects
```

Lets now use `myhandle` to create directories
```bash
kick start myhandle ~/projects/myproject
```

Upon examining the folder tree, one will see that any file or folder using the
variable `${PROJECT_NAME}` has been replaced with `myproject`. Any file or
folder that contains a variable in its name will be interpolated by the
variables value.
```bash
cd ~/projects
tree -an myproject
myproject
├── cmd
│   └── myproject
│       └── main.go
├── docs
├── .gitignore
├── README.md
└── scripts
    └── myproject.sh

4 directories, 4 files
```

## Variables

Variables are either...

1. predefined variables
1. environment variables
1. variables stored in `~/.env`

The order of precedence is as above.

The variables stored in `~/.env` are key value pairs and take the form
`key=value`.

Create some environment variables which will be used later. Using your favorite
editor add the following variables to `~/.env`
```text
author="JOHN SMITH"
email=john.smith@email.com
```

NOTE: The `~/.env` file isn't just for the `kick` utility but can be used by
any other program in the same way that `~/.profile` is not for just one tool.

### Variable Functions

For a list of supported variable functions see [Supported Variable
Functions](reference.md) in the reference section.

## Template files

Templates files are any text file which contains a modeline. A modeline is a
piece of text that informs kick that the file is a template. A modeline takes up
the form `kick:render` and should be placed within the first 5 lines of a text
file. Modeline lines are stripped from the file and can be placed inside any
comment type.

Example modelines
```text
kick:render
```

```bash
# kick:render
```

```c
// kick:render
```

```html
<!--- kick:render -->
```

The next exercise is to add 4 template variables to the `README.md` file in our
template.

* `${PROJECT_NAME}` - Predefined variable
* `${USER}`         - Environment variable built into shells (E.G. Bash) which
                      contains the current username.
* `${author}`       - Variable from "~/.env"
* `${email}`        - Variable from "~/.env"

Change the contents of the `~/templates/mytemplate/README.md` by cutting and
pasting the whole text below...

```bash
cat <<'EOF' > ~/templates/mytemplate/README.md
<!--- kick:render -->
# ${PROJECT_NAME}

AUTHOR: ${author}
EMAIL: ${email}
USER: ${USER}

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.
EOF
```

Create a new project using the same handle `myhandle` ...

```bash
kick start myhandle ~/projects/project2
```

Inspect the contents of the `README.md` file.

```
cat ~/projects/project2/README.md
```

```markdown
# project2

AUTHOR: JOHN SMITH
EMAIL: john.smith@email.com
USER: vagrant

Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.
```

## Git Project Templates

Kick can use remote git repositories as stores for project templates.  Using our
freshly minted project template we can upload this to a git repository and use
`kick install` to install the directory.

Create a git repository and upload it to a remote repository.
```
cd ~/templates/mytemplate
git init
git add .
git commit -m 'Initial commit'
git remote add origin git@github.com/username/mytemplate.git
git push -u origin master
```

Remove the old handle
```bash
kick remove mytemplate
```

Add a remote location
```bash
kick install mytemplate git@github.com/username/mytemplate.git
```

Start a new project
```bash
kick start mytemplate ~/project/project3
```

If you would like to find out how to create a repository of project templates,
head to the repositories section.
