# Gradle for Agents (`gfa`)

[![GoReleaser](https://github.com/silverAndroid/gradle-for-agents/actions/workflows/release.yml/badge.svg)](https://github.com/silverAndroid/gradle-for-agents/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Gradle for Agents** (`gfa`) is a smart, noise-reducing wrapper for the Gradle build tool designed specifically for LLM and AI coding agents. 

When agents run standard Gradle builds (`./gradlew build`), they are flooded with thousands of lines of output—download progress, task success statements, and irrelevant compile warning noise. This wastes context window tokens, increases costs, and makes it incredibly difficult for the agent to find and resolve compilation errors.

`gfa` intercepts Gradle's output, displays a clean progress indicator, and outputs a concise **TOON (Task-Oriented Output & Noise-reduction)** summary when a build fails or when explicitly asked.

---

## Features

- 🤫 **Noise Suppression**: Suppresses verbose download progress, task listings, and non-critical details during execution.
- 🎯 **TOON Output summaries**: Automatically extracts the exact tasks that failed along with their immediate error context (such as compiler errors or stack traces).
- 📁 **Full Output Preservation**: Logs the complete verbose Gradle output to a temporary log file (e.g., `/tmp/gfa-20260711-130000/full_output.log`) so you or the agent can inspect it if needed.
- 🔄 **Auto-Updates**: Automatically checks for updates in the background (when not installed via Homebrew) and prompts you to update.

---

## Installation

### As an Agent Skill (for AI Agents)
You can install this repository as an Agent Skill for compatible AI developer harnesses (like Antigravity, Kilo Code, OpenCode, Cursor, and Claude Code). This teaches agents to autonomously use `gfa` for all Gradle operations in your project.

1. Install the skill in your project workspace:
   ```bash
   npx skills add silverAndroid/gradle-for-agents/skills/gfa
   ```

2. (Optional) Run the setup script to register the `/gfa` custom slash command in Claude Code:
   ```bash
   ./.skills/gfa/scripts/setup.sh
   ```
   *(Note: Depending on your harness, the path may be `./.skills/gfa/scripts/setup.sh` or `./.agents/skills/gfa/scripts/setup.sh`)*

### Via Curl (macOS & Linux)
You can install `gradle-for-agents` and `gfa` globally using the installation script:

```bash
curl -fsSL https://raw.githubusercontent.com/silverAndroid/gradle-for-agents/main/install.sh | bash
```

### Via Homebrew (macOS)
You can install it through the custom Homebrew Tap:

```bash
brew install silverAndroid/tap/gradle-for-agents
```

### From Source
If you have Go installed, you can compile it directly:

```bash
go install github.com/silverAndroid/gradle-for-agents/cmd/gfa@latest
go install github.com/silverAndroid/gradle-for-agents/cmd/gradle-for-agents@latest
```

---

## Usage

Simply substitute `gradle` or `./gradlew` with `gfa` (or `gradle-for-agents`):

```bash
# Run tasks showing progress and summary on completion
gfa assembleDebug

# Run builds and include compiler warnings in the output upon success
gfa build --show-warnings

# Passing through help or version checks
gfa --help
gfa --version
```

### Options

- `--show-warnings`: Output warnings in the summary on successful builds.
- `--version` / `-v`: Prints `gfa` and local gradle version information and exits.
- `--help` / `-h`: Shows help information.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
