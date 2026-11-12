# Changelog CLI specification

## Commands

- Invoking `clbot` with no arguments is supported and prints the CLI's own `Version` and usage help listing the supported commands; this version is distinct from the latest released version.
- `clbot version` prints the calculated next release version by applying the included changes to the latest release (for example, a patch change after `v2.18.1` prints `2.18.2`).
- `clbot released` and `clbot latest` are aliases that print only the latest released version, without its leading `v`.
- `clbot list` lists relevant changes from the latest released version to now, omitting changes labeled `[skip-cl]`. Each entry starts with an `ID TYPE TITLE` row, followed by its description when present; entries are separated by a blank line. Unset-type changes are included in the output, but the command then reports `found N unset changelog entry/entries` and exits with an error.
- `clbot validate {input}` validates the specified input; `clbot verify {input}` is an alias. The special input `-` reads all input from stdin: the first line is the change title and the remainder is its description. Surrounding whitespace is trimmed, and exactly one recognized marker must occur across the title and description combined.

The exported `internal/clbot.Run` API has the signature `Run(ctx context.Context, inputs RuntimeInputs) error`. `RuntimeInputs` contains `Args []string`, `Stdin io.Reader`, `Stdout io.Writer`, `Stderr io.Writer`, and `DataSource DataSource`. Callers pass the reader that supplies the `-` input (the CLI passes `os.Stdin`); a nil reader or read failure is reported as a validation error. For example, `printf '%s\n%s' "$PR_TITLE" "$PR_BODY" | clbot validate -` validates a pull request title and body without interpreting either as a shell command.

## Data and change classification

The datasource contract provides the latest release and a collection of changes, each with an ID, type, title, and description. Git is the default provider; providers must be replaceable through an injectable datasource interface. For Git, a change ID is the full commit SHA. A future file datasource will use the filename as the change ID.

The change type is one of `skip-cl`, `patch`, `feature`, `breaking`, or `unset`. It is derived from exactly one recognized marker across the title and description; `[bugfix]` maps to type `patch`. With no recognized marker or multiple recognized markers, the type is `unset`. `clbot list` includes unset changes in its output (but omits `skip-cl` changes) and returns a count-based error after output; validation continues to reject missing or duplicate markers.

Recognize these labels on changes:

- `[skip-cl]`: omit the change from changelog output and version-bump calculation.
- `[bugfix]` or `[patch]`: patch change.
- `[feature]`: minor change.
- `[breaking]`: major change.

For version calculation, if any included change is labeled `[breaking]`, bump major; otherwise, if any is labeled `[feature]`, bump minor; otherwise, bump patch, including when there are no `[feature]` or `[breaking]` changes. Apply the selected bump to the latest `X.Y.Z` version using semver resets: major increments `X` and resets `Y` and `Z` to zero; minor increments `Y` and resets `Z` to zero; patch increments `Z`.

## Architecture and implementation constraints

Tests must use `github.com/stretchr/testify/assert` and `github.com/stretchr/testify/require`.

`cmd/clbot` contains only a thin `main` that delegates to the exported `Run` in `internal/clbot`; all CLI command and argument logic, as well as changelog implementation, live in `internal/clbot`. Use the standard-library `flag` package, not Cobra.

Phase one builds an injectable, empty shell: define the seams and command structure without implementing runtime behavior. Write tests first (TDD) that specify full intended functionality. Unit tests use a fake Git provider; integration tests and their fixtures live under `internal/clbot/testdata`.
