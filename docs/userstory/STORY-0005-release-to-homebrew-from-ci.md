# 5. Release to Homebrew from CI

Date: 2026-10-18

## Status

Accepted

## Story

As the detoo maintainer, I want `make release` to publish detoo to Homebrew
through CI, so that people can install it in one step without installing a Go
toolchain, and detoo is more accessible.

## Acceptance criteria

* `make release` tags the release and pushes the tag. The tag push starts a
  CI job that publishes the release to Homebrew.
* The release is published to a Homebrew tap in the `mustardcustarddev`
  organisation.
* Each release has a semantic version, such as `v1.0.0`, which a
  semantic-versioning tool works out from the conventional commits since the
  last release.
* On an Apple silicon Mac, `brew install` installs `adr`, `rfc`, `sty`, `plm`,
  `sop` and `detoo` as prebuilt binaries, with no Go needed.
* `brew install` also installs the tools' shell completion.
* `detoo` reports the release's version, such as `1.0.0`, not a commit hash.
* The README and the install page describe installing with Homebrew.

## Notes

Today the tools are installed by cloning the repository and running
`make deploy`, which needs Go 1.26 or later, `make` and `git`. `make release`
tags the current commit with its short hash and creates a GitHub release. No
release has been made yet.

The aim is to reach homebrew-core later, in another story. For a
self-submission by the repository owner, Homebrew's
[Package Acceptance Policy](https://docs.brew.sh/Package-Acceptance-Policy)
needs at least 90 forks, 90 watchers or 225 stars. 
homebrew-core also needs a stable release with an immutable tag, which
version numbers will give.

Out of scope: Intel Macs.

Open questions:

* Should the release also support Homebrew on Linux?
* Which shells' completion should Homebrew install: bash only, as
  `make deploy` does, or zsh and fish too?
