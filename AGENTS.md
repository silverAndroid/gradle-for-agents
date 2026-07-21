# AI Agents Integration Guide

If you are an AI coding agent (like Antigravity or any LLM-based agent) working on this codebase, or any project that uses Gradle, **please read this document**.

## Why Use `gfa`?

Standard Gradle build output is extremely verbose and noisy. A typical compilation run prints:
- Detailed download progression for dependencies
- Status messages for every executed task
- Dozens of compiler warnings and warnings about deprecated Gradle features

This massive volume of output:
1. **Bloats your context window**, consuming valuable input tokens and raising inference costs.
2. **Obscures compile errors** under thousands of lines of log noise, making it harder to pinpoint why a build failed.

`gfa` is a Gradle wrapper designed specifically for agents. It runs Gradle in the background, prints quiet progress updates (`.....`), and only outputs clean, high-signal information:

- **On success**: Prints a single line confirming build completion and the location of the verbose log.
- **On failure**: Outputs a targeted **TOON (Task-Oriented Output & Noise-reduction)** summary showing exactly which task failed along with its immediate error messages or compiler traces.

---

## Invocation Guideline

Whenever you need to verify if code compiles, run unit tests, or build the project, use `gfa` (or `gradle-for-agents`) instead of `gradle` or `./gradlew`:

```bash
# Good
gfa build

# Good
gfa test

# Avoid (noisy, will pollute your context)
./gradlew build
```

---

## Reading Failures (TOON Summaries)

When a build fails, `gfa` outputs a summary like this:

```
INFO: Logging full output to: /tmp/gfa-20260711-130000/full_output.log
......
INFO: Build finished, generating summary...

--- TOON OUTPUT ---
ERROR [Task :app:compileDebugJavaWithJavac]:
/Users/user/workspace/app/src/main/java/MainActivity.java:12: error: cannot find symbol
    System.out.println(nonExistentVar);
                       ^
  symbol:   variable nonExistentVar
  location: class MainActivity
...

FAILURE: Build failed with exit code 1. Full logs at: /tmp/gfa-20260711-130000/full_output.log
```

Use the `ERROR [Task ...]` output block directly to fix the syntax or compilation issues in your workspace.

---

## Accessing Full Logs

If the TOON output summary is insufficient (e.g. you need to examine the full stack trace, Gradle configuration error, or dependency resolution issue), do not re-run the build with raw `./gradlew`.

Instead, read the full log file path printed in the success/failure footer:
1. Copy the path to the log file (e.g. `/tmp/gfa-XXXXXXXX-XXXXXX/full_output.log`).
2. Use your file reading tool to inspect its content directly.

---

## Agent Sandbox Compatibility

If you are running in a restricted agent environment where you cannot write to `/tmp` (e.g. strict workspace-only access), you can override where `gfa` creates its log directories by setting the `GFA_LOG_DIR` environment variable.

```bash
# E.g. placing logs inside a workspace folder instead of /tmp
GFA_LOG_DIR=./.gfa-logs gfa build
```

---

## Code Modification Rules

**Rule: Test Coverage**
Whenever you make changes to source code (adding new features, fixing bugs, or modifying logic), you MUST ensure the changes have ample test coverage.
- Check existing test files (`*_test.go`, etc.) and update them if they apply.
- If no applicable tests exist, create new test cases to verify the new or modified logic.
- Run the tests to ensure they pass and accurately cover the changes before completing the task.
