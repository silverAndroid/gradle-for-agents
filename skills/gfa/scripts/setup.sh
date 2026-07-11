#!/usr/bin/env bash

set -euo pipefail

# 1. Determine directories relative to script location
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SKILL_DIR="$(dirname "$SCRIPT_DIR")"
TEMPLATE_PATH="$SKILL_DIR/templates/gfa.md"

# 2. Locate project root (search upwards for .git, build.gradle, or use current dir)
PROJECT_ROOT=""
CURRENT_DIR="$(pwd)"
while [ "$CURRENT_DIR" != "/" ]; do
    if [ -f "$CURRENT_DIR/build.gradle" ] || [ -f "$CURRENT_DIR/build.gradle.kts" ] || [ -d "$CURRENT_DIR/.git" ]; then
        PROJECT_ROOT="$CURRENT_DIR"
        break
    fi
    CURRENT_DIR="$(dirname "$CURRENT_DIR")"
done

if [ -z "$PROJECT_ROOT" ]; then
    PROJECT_ROOT="$(pwd)"
fi

echo "Project root detected: $PROJECT_ROOT"

# 3. Check if gfa is installed in PATH
if ! command -v gfa >/dev/null 2>&1; then
    echo "WARNING: 'gfa' binary was not found in your PATH."
    echo "To install 'gfa' and 'gradle-for-agents' globally, please run:"
    echo "  curl -fsSL https://raw.githubusercontent.com/silverAndroid/gradle-for-agents/main/install.sh | bash"
    echo ""
else
    echo "Check: 'gfa' binary is installed in your PATH."
fi

# 4. Copy Claude Code command template
CLAUDE_COMMANDS_DIR="$PROJECT_ROOT/.claude/commands"
TARGET_FILE="$CLAUDE_COMMANDS_DIR/gfa.md"

if [ -f "$TEMPLATE_PATH" ]; then
    echo "Creating directory: $CLAUDE_COMMANDS_DIR"
    mkdir -p "$CLAUDE_COMMANDS_DIR"
    
    echo "Installing Claude Code command template to $TARGET_FILE..."
    cp "$TEMPLATE_PATH" "$TARGET_FILE"
    echo "SUCCESS: Custom slash command /gfa is now registered for Claude Code in this project."
else
    echo "ERROR: Template file not found at $TEMPLATE_PATH" >&2
    exit 1
fi
