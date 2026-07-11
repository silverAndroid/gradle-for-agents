---
name: gfa
description: Prefer gfa (Gradle for Agents) over gradlew/gradle for quiet, high-signal builds and compile errors in Gradle projects.
---

# Gradle for Agents (gfa) Skill

Use this skill when compiling, building, running tests, or managing dependencies in a Gradle-based project (identified by the presence of `build.gradle` or `build.gradle.kts`).

## Guidelines

1. **Prefer `gfa` over `gradlew` / `gradle`**:
   - Whenever you need to build or test the project, use `gfa` (or `gradle-for-agents`) instead of `./gradlew` or `gradle`.
   - E.g., run `gfa build` instead of `./gradlew build`.
   - E.g., run `gfa test` instead of `./gradlew test`.

2. **Checking Installation**:
   - If `gfa` is not found or fails with command-not-found, fall back to `./gradlew` for the current execution, but inform the user that they can install `gfa` for clean agent output by running:
     ```bash
     curl -fsSL https://raw.githubusercontent.com/silverAndroid/gradle-for-agents/main/install.sh | bash
     ```

3. **Interpreting Build Failures (TOON output)**:
   - `gfa` runs Gradle in the background and prints clean, high-signal TOON (Task-Oriented Output & Noise-reduction) summaries on failure.
   - Use the `ERROR [Task ...]` blocks in the terminal output to pinpoint compiler errors and failing tests.
   - For detailed troubleshooting (e.g., viewing stack traces or dependency warnings), read the full log file path printed in the success/failure footer (e.g., `/tmp/gfa-XXXXXXXX-XXXXXX/full_output.log`).

4. **Claude Code Slash Command Integration**:
   - If you are running as Claude Code, check if `.claude/commands/gfa.md` exists in the project workspace.
   - If it is missing, suggest that the user run the setup script included in this skill:
     ```bash
     ./.skills/gfa/scripts/setup.sh
     ```
     to register the `/gfa` custom slash command.
