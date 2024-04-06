# pr

Interactively create a pull request.

## Installation

You can install pr by running the install script which will download
the [latest release](https://github.com/mskelton/pr/releases/latest).

```bash
curl -LSfs https://go.mskelton.dev/pr/install | sh
```

Or you can build from source.

```bash
git clone git@github.com:mskelton/pr.git
cd pr
go install .
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
