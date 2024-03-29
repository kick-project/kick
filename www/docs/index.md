# Kick

A project template tool that stores project templates in git repositories

# Introduction

kick is a project template tool, some of its features are...

* the ability to install templates using a git URL
* the use of a directory path as template
* use pre existing projects as templates
* create a repository of templates
* simple key/value pair or environment variables as template variables
* bash style template variables Examples: `$variable` `${variable}`
# Examples

Installing templates directly from git
```bash
kick install go https://github.com/kick-project/template-go.git      # Install a go template
kick install pypi https://github.com/kick-project/template-pypi.git  # Install a pypi template
kick install custom_handle1 git@example.com:your/git/project.git     # Install from a private git repository
```

Use a local directory as a template
```bash
kick install mytemplate ~/template_directory/mytemplate          # Install a custom template from disk
```

Starting a project
```bash
kick start go ~/workspace/mygoproject
kick start pypi ~/workspace/mypypiproject
kick start mytemplate ~/myproject
# Or simply to create a project in the currrent directory
kick start mytemplate myproject
```

Search and install templates from a repo
```bash
# Search for a template
kick search tmpl
+-------------+---------------------------------+
|  TEMPLATE   |            LOCATION             |
+-------------+---------------------------------+
| tmpl/repo1  | http://gitservice.com/tmpl.git  |
| tmpl1/repo1 | http://gitservice.com/tmpl1.git |
| tmpl2/repo1 | http://gitservice.com/tmpl2.git |
+-------------+---------------------------------+

# Install template using <template> name
kick install mytmpl tmpl

# Install template using <template>/<repo> name
kick install mytmpl1 tmpl1/repo1
```

Its that Simple!

If you would like to see more, head on over to the install & quick start sections
