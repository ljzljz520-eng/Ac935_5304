# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	hospitaldesk/audit	[no test files]
?   	hospitaldesk/auth	[no test files]
?   	hospitaldesk/cmd/clinicdesk	[no test files]
?   	hospitaldesk/model	[no test files]
?   	hospitaldesk/policy	[no test files]
?   	hospitaldesk/schedule	[no test files]
?   	hospitaldesk/service	[no test files]
?   	hospitaldesk/storage	[no test files]
?   	hospitaldesk/training	[no test files]
--- FAIL: TestEmployeeCannotDownloadUnpublishedPolicy (0.00s)
    workflow_test.go:144: employee downloaded unpublished policy
FAIL
FAIL	hospitaldesk	0.023s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/clinicdesk): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/clinicdesk): exit `0`
- Frontend build (web): exit `0`
