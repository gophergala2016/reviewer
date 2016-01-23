# reviewer [![Build Status](https://travis-ci.org/gophergala2016/reviewer.svg)](https://travis-ci.org/gophergala2016/reviewer)

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
       myevencoolaprip:
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

TODO!

  [YAML]: http://yaml.org/ "YAML format homepage"
  [TOML]: https://github.com/toml-lang/toml "TOML format definition"
  [JSON]: http://www.json.org/ "JSON format homepage"
  [GitHub API token]: https://github.com/settings/tokens "GitHub profile tokens"
  [GitHub]: https://github.com "GitHub home page"

## Usage

You can use Reviewer basic functionality by simply invoking it directly:

    $ reviewer

TODO!
