
# GitHub User Activity CLI

## Installation

To install the CLI, run the following command:

```sh
go install github.com/EssaAlshammri/github-activity
```

## Usage
The CLI provides two output formats: summary and all. By default, it uses the summary format.

Command Syntax
```sh
github-activity [--format=<format>] <username>
```

#### Options
--format: Output format, either summary or all (default: summary).


## Examples
To get a summary of recent activity for a user:
```sh
github-activity EssaAlshammri
```