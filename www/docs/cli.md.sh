#!/usr/bin/env bash

cat <<EOF
# Use & manage templates
## kick start

\`\`\`bash
$(./kick start -h)
\`\`\`

## kick list

\`\`\`bash
$(./kick list -h)
\`\`\`

## kick install

\`\`\`bash
$(./kick install -h)
\`\`\`

## kick remove

\`\`\`bash
$(./kick remove -h)
\`\`\`

## kick search

\`\`\`bash
$(./kick search -h)
\`\`\`

# Management commands

## kick setup

\`\`\`bash
$(./kick setup -h)
\`\`\`

# Repository management

## kick repo

\`\`\`bash
$(./kick repo -h)
\`\`\`

## kick init

\`\`\`bash
$(./kick init -h)
\`\`\`
EOF