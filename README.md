# pr

Interactively create a pull request.

## Installation

```bash
go install github.com/mskelton/pr@latest
```

## Prerequisites

To parse and format the Jira ticket from your branch name, add a `JIRA_PREFIXES`
environment variable to your shell with a list of supported Jira prefixes.

```bash
export JIRA_PREFIXES="ABC,SUPPORT"
```

## Usage

```bash
pr
```
