---
description: Run gradle-for-agents (gfa) build or test commands quietly and parse TOON summaries.
---

Run the `gfa` tool on the codebase with the arguments provided:
!gfa $ARGUMENTS

If the build fails, examine the TOON output summary to locate and fix any compilation or test errors.
If `gfa` is not installed or not found, let the user know they can install it by running:
`curl -fsSL https://raw.githubusercontent.com/silverAndroid/gradle-for-agents/main/install.sh | bash`
