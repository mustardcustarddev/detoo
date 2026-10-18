# 1. Release detoo

Date: 2026-10-18

## Status

Accepted

## Steps

1. Make sure you are the detoo maintainer, who releases detoo, and that you
   have Go 1.26 or later, `make`, `git`, push access to `main` on
   `mustardcustarddev/detoo`, and the GitHub CLI, `gh`, signed in.
2. Make sure the repository's `HOMEBREW_TAP_DEPLOY_KEY` secret is set, which
   the Release workflow uses to push to `mustardcustarddev/homebrew-tap`. If it
   is not, set it up once, as the Release section of
   [CONTRIBUTING.md](../../CONTRIBUTING.md#release) describes.
3. Switch to `main` with `git switch main`, and check that
   `git status --porcelain` prints nothing. `make release` refuses to run
   otherwise.
4. Preview the version that will be released:

   ```bash
   go run github.com/caarlos0/svu/v3@v3.4.1 next
   ```

   svu works it out from the conventional commits since the last release: a
   `fix:` makes a patch release, a `feat:` a minor one, and a breaking change,
   such as `feat!:`, a major one.
5. Release with `make release`, or with `make release VERSION=<version>`, such
   as `v1.0.0`, to set the version yourself. It builds, vets and tests the
   tools, tags `main` with an annotated tag `detoo <version>`, and pushes
   `main` and the tag. Expect `Releasing <version>`.
   * If it reports `release: <version> is already tagged, and no fix or
     feature since needs a release`, stop: there is nothing to release.
   * If it reports any other `release:` error, the build fails, or
     `git push` fails, stop, and escalate. If the push failed, the tag was
     already created locally.
6. Find the Release workflow's run for the tag with
   `gh run list --workflow release.yml --limit 1`, then watch it with
   `gh run watch <run-id> --exit-status`, and wait for it to succeed.
   * If it fails, stop, and escalate. A failed run can be run again: it
     replaces the release's files, and the formula with them.
7. Check the GitHub release with
   `gh release view <version> --repo mustardcustarddev/detoo`. It lists
   `detoo_<x.y.z>_darwin_arm64.tar.gz` and `checksums.txt`, where `<x.y.z>` is
   the version without its `v`.
8. Check that `Formula/detoo.rb` in `mustardcustarddev/homebrew-tap` has a
   commit `detoo <x.y.z>`.
9. On an Apple silicon Mac, install the release with `brew upgrade detoo`, or,
   if detoo is not installed yet, with `brew tap mustardcustarddev/tap` and
   then `brew install detoo`.
10. Run `detoo doctor` in a project, and check that its header shows
    `<x.y.z>` beside the detoo badge, not a commit hash.
