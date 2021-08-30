# auto-change-log

`auto-change-log` automates the process of updating a changelog based
on git history.

## Installation

`go install github.com/tzneal/auto-change-log@latest`

## Features

- Can append to a changelog by parsing the date of the latest release to only add new entries
- Configurable classifiers used to identify which commmits are new features, changes, bug fixes, etc.
- Can generate a changelog for arbitrary points in time specified by tags or dates

## Usage

`auto-change-log generate`

## Sample Config File
repositorypath - path the git repository, if empty it defaults to the current working directory.

classifyrules - rules used to identify the type of change

issueextractor - rules used to extract issue numbers for appending to the change summary, optionally linkin g it in markdown if a linkurl is provided

```yaml
filename: CHANGELOG.md
defaultentrytype: Changed
repositorypath: ""
ignoremerge: true
classifyrules:
- regexp: "(?i)\\b(removed|removing)\\b"
  entrytype: Removed
  priority: 50
- regexp: "(?i)\\b(deprecated|deprecating)"
  entrytype: Deprecated
  priority: 50
- regexp: "(?i)\\b(add|added|adding)\\b"
  entrytype: Added
  priority: 50
- regexp: "(?i)\\b(fix|fixed|fixes)\\b"
  entrytype: Fixed
  priority: 50
- regexp: "(?i)\\bbug\\b"
  entrytype: Fixed
  priority: 50
- regexp: "(?i)\\bsecurity"
  entrytype: Security
  priority: 50
issueextractor:
- regexp: "(ISSUE-\\d+)"
  linkurl: http://my.bug.tracker/$1/view
```
## Help

```
$ auto-change-log generate --help
writes the change log

Usage:
auto-change-log generate [flags]

Flags:
-h, --help           help for generate
-n, --name string    the release name, this must match the newest release if appending (default "Unreleased")
--overwrite      overwrite the current changelog entries. If not set, append.
-s, --since string   only look at commits since this time, specified as YYYY-MM-DD or a tag name
-u, --until string   only look at commits up to this time, specified as YYYY-MM-DD or a tag name

```