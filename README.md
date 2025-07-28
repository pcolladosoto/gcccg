# The Go Conventional Commit Changelog Generator
Typing up some release notes can be quite a daunting task. Now, if you adhere to the
[Conventional Commit Spec](https://www.conventionalcommits.org/en/v1.0.0/) generating
an automatic changelog becomes that much easier.

Changelogs need to be automatically generated some times. A paradigmatic case is the
use of CI/CD pipelines in whatever environment one is working on. However, when
browsing available changelog generators nothing really hit the spot. So... why not
write one?

This task is actually bearable thanks to the great [leodido/go-conventionalcommit][]
parser (which is based on [Ragel][https://www.colm.net/open-source/ragel/]) and the
Go-native Git implementation provided by [go-git/go-git][].

## Installation
Installing this tool is as simple as running

    $ go install github.com/pcolladosoto/gcccg@latest

## Usage
The command ships with (hopefully) sensible defaults. It'll inspect a Git repository on the current
directory and generate the changelog between the latest and second-to-last tags. One can however
manually specify what tags to generate the changelog in between.

The output can be printed to `stdout` (so that it can usually be redirected to a file) and/or stored
on a file.

One can also influence the generated changelog with several options including the addition of
GitHub-specific flavouring, the inclusion of contributor emails and so on.

We'll add a manpage at some point documenting available options, but for now the following examples
should be informative enough: `gcccg` is not really that complex!

### Automatic tag detection
Detecting tags automatically relies on the alphabetical order of tags. That is, the latest tag is
assumed to be the latest one in alphabetical order (as given by `sort.Strings`). This lines up with the
[Semantic Versioning](https://semver.org) specification which is what we believe people usually stick
to. Beware of possible pitfalls if adding tags following other formats! If doing so, we recommend
leveraging the `--to-tag` and `--from-tag` options as seen in the examples below.

### Examples
The following show usual invocations of `gcccg`:

    # Generate the changelog between the latest and second-to-last tags on a repository
    # stored locally at /path/to/repo and print it to stdout
    $ gcccg --repo /path/to/repo --stdout

    # Generate the changelog between tags v2.0.0 and v1.0.0 for a repository stored
    # locally at /path/to/repo and save it to changelog.md on the current directory
    $ gcccg --repo /path/to/repo --from-tag v2.0.0 --to-tag v1.0.0 --out changelog.md

    # Generate the changelog between tags v2.0.0 and v1.0.0 for a repository stored
    # locally at /path/to/repo, including contributor emails and print it to stdout
    $ gcccg --repo /path/to/repo --from-tag v2.0.0 --to-tag v1.0.0 --email --stdout

<!-- REFs -->
[leodido/go-conventionalcommit]: https://github.com/leodido/go-conventionalcommits
[go-git/go-git]: https://github.com/go-git/go-git
