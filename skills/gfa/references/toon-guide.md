# TOON Output & Noise-reduction Guide

This guide describes how to understand and handle output when building or testing Gradle projects using `gfa` (Gradle for Agents).

## Why TOON?

Standard Gradle output contains thousands of lines of progress bars, dependency downloads, and compiler warnings. This results in:
1. Massive token consumption (bloating your context window).
2. Failure to locate compilation errors because they are buried in logs.

`gfa` intercepts the Gradle execution in the background, only printing a single status line (`.....`) until the build completes.

## Success Outputs

On a successful build, `gfa` outputs:
```
INFO: Logging full output to: /tmp/gfa-20260711-130000/full_output.log
......
INFO: Build finished, generating summary...
SUCCESS: Build completed successfully. Full logs at: /tmp/gfa-20260711-130000/full_output.log
```

## Failure Outputs (TOON Summary)

On failure, `gfa` prints a concise, filtered block containing only compile/test failures:
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

FAILURE: Build failed with exit code 1. Full logs at: /tmp/gfa-20260711-130000/full_output.log
```

### Action Plan on Failure:
1. Parse the path and line number from the `ERROR [Task ...]` block (e.g. `/Users/user/workspace/app/src/main/java/MainActivity.java:12`).
2. Open that file and fix the code symbol or syntax error.
3. Re-run `gfa build`.
4. If the error is not clear from the TOON block, open the absolute path shown in `Full logs at: <path>` to read the verbose stack trace. Do NOT re-run raw `./gradlew`.
