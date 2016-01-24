# reviewer [![Build Status](https://travis-ci.org/gophergala2016/reviewer.svg)](https://travis-ci.org/gophergala2016/reviewer) [![ReportCard][ReportCard-Image]][ReportCard-Url]

Reviewer is a tool for running easy code reviews in GitHub.

## Installation

    go install github.com/gophergala2016/reviewer

## Configuration

For using Reviewer, you will need to get a [GitHub API token].
Reviewer uses a configuration file, by default located at `$HOME/.reviewer.yaml`.
If you want to set this file to other location, name, or accepted format, you can use the option `--config` with an argument specifying the filename.
Accepted formats for the configuration file are:

  - [YAML]
  - [TOML]
  - [JSON]

An example YAML file would have this format:

    authorization:
       token: MYNICEANDSHINYGITHUBAPITOKEN
    repositories:
       mycoolapp:
           username: cooldeveloper
           status: true
           required: 3
       myevencoolapi:
           username: cooldeveloper
           status: true
           required: 2

where:

  - `authorization` contains:
      - `token`: corresponds to user's [GitHub API token]. This key can be also given throught REVIEWER_TOKEN environment variable.
  - `repositories` consists on a set of subsets defined by the repository name in [GitHub], and containing a set of keys with different meanings:
      - `username`: Would correspond to the username holding the repository to be checked.
      - `status`: Defining whether the repository is, or is not, enabled for checking.
      - `required`: Corresponds to the number of approvals required to go on with the merge, in case nothing else blocks it.

You can get Reviewer's configuration by invoking the command configure:

      $ reviewer configure
      Using config file: /home/user/.reviewer.yaml
      - cooldeveloper / mycoolapp ENABLED +1:3
      - cooldeveloper / myevencoolapi ENABLED +1:2

  [YAML]: http://yaml.org/ "YAML format homepage"
  [TOML]: https://github.com/toml-lang/toml "TOML format definition"
  [JSON]: http://www.json.org/ "JSON format homepage"
  [GitHub API token]: https://github.com/settings/tokens "GitHub profile tokens"
  [GitHub]: https://github.com "GitHub home page"

## Usage

You can use Reviewer basic functionality by simply invoking it directly:

      $ reviewer
      Using config file: /home/user/.reviewer.yaml
      + cooldeveloper/mycoolapp
      + cooldeveloper/mycoolapi
        - 47 NOP   (Changes CI badge location on README.md) score 1 of 3 required

Where it first reports the configuration file used, if any.
Then, it prints a list of the repos, with a list of the PRs pending to merge.
For each PR, it shows:
  - Pull request identifier
  - Operation done, i.d. `NOP` if it doesn't satisfies requirements to be done, `MERGED`, if it was merged.
  - Pull Request title, between brackets.
  - Score from the approvals in the pull request comments.
  - How many approvals were required.

[ReportCard-Url]: http://goreportcard.com/report/gophergala2016/reviewer
[ReportCard-Image]: http://goreportcard.com/badge/gophergala2016/reviewer
