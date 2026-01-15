---
name: investigate-e2e-failure
description: This gives instructions on how to investigate end-to-end test failures.
---

## When to Use

Whenever the user asks "investigate the e2e failure 123", execute the instructions below.

## Prerequisites

Check if the `gh` binary is installed.

## Instructions

### Download the overview

In order to prevent TLS issues, just execute `gh` with full permissions.

1. Use the following command to get an overview about the failure:
```shell
$ gh run view --job 123

X main Test · 19282487500
Triggered via schedule about 2 months ago

X cloud-tests / e2e-tests / e2e (v1.34.0-kind, encryption-at-rest) in 10m48s (ID 123)
  ✓ Set up job
  ✓ Run actions/checkout@v5
  ✓ Install devbox
  ✓ Generate kustomized all-in-one install configs
  ✓ Extract k8s version/platform
  X Run CI E2E tests
  ✓ Upload logs on failure
  ✓ Post Install devbox
  ✓ Post Run actions/checkout@v5
  ✓ Complete job

ANNOTATIONS
! No files were found with the provided path: output/**. No artifacts will be uploaded.
cloud-tests / e2e-tests / e2e (v1.34.0-kind, encryption-at-rest): .github#20

X Process completed with exit code 2.
cloud-tests / e2e-tests / e2e (v1.34.0-kind, encryption-at-rest): .github#2510


To see the logs for the failed steps, try: gh run view --log-failed --job=123
View this run on GitHub: https://github.com/mongodb/mongodb-atlas-kubernetes/actions/runs/19282487500
```

### Find the corresponding e2e test
In this example we see that a test in test/e2e with the label "encryption-at-rest" failed.
In this case it is the test defined in `test/e2e/encryption_at_rest/encryption_at_rest_test.go`.

### Analyze the failure
Use the following command to download the logs for the failed job:

```shell
$ gh run view --job 123 --log
```

### Investigate
1. Correlate the logs with the test code to identify the root cause of the failure.
2. Look for error messages, stack traces, or any other indicators of what went wrong.
3. Check for any recent changes in the codebase that might have introduced the failure.
4. If necessary, reproduce the failure locally to gain a deeper understanding.
5. Document your findings and suggest potential fixes or improvements to prevent similar failures in the future.
