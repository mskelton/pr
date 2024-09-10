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

To parse and format a Linear or Jira ticket from your branch name, add a
`PR_TICKET_PREFIXES` environment variable to your shell with a list of supported
prefixes.

```bash
export PR_TICKET_PREFIXES="ABC,SUPPORT"
```

You can also specify prefixes to strip out of your PR names. This is useful if
you include your username as a prefix in your branches.

```bash
export PR_STRIP_PREFIXES='username/'
```

## Usage

```bash
pr
```
