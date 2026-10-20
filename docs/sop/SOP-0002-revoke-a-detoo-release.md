# 2. Revoke a detoo release

Date: 2026-10-18

## Status

Accepted

## Steps

1. Make sure you are the detoo maintainer, who releases detoo, and that you
   have `git`, push access to `main` on `mustardcustarddev/detoo` and to
   `mustardcustarddev/homebrew-tap`, and the GitHub CLI, `gh`, signed in. This procedure withdraws a release that has
   been published, as [SOP 1](SOP-0001-release-detoo.md) publishes it. It
   does not decide whether to: decide that first.
2. Write down the version to revoke, `<version>`, such as `v0.2.0`, and the
   same version without its `v`, `<x.y.z>`.
3. Stop the release from being published again while you revoke it. List the
   Release workflow's runs with
   `gh run list --workflow release.yml --repo mustardcustarddev/detoo`, and
   cancel any run for `<version>` that has not finished, with
   `gh run cancel <run-id> --repo mustardcustarddev/detoo`. Run the list
   again, and check that no run for `<version>` is queued or in progress.
4. Clone the tap with `gh repo clone mustardcustarddev/homebrew-tap`, and go
   into it with `cd homebrew-tap`. Check which release the formula installs
   with `grep url Formula/detoo.rb`.
   * If its URL does not contain `/releases/download/<version>/`, the tap
     already installs another release: go on to step 7.
5. List the formula's history with
   `git log --format='%h %s' -- Formula/detoo.rb`. Each release has a commit
   `detoo <a.b.c>`. Find the latest commit for a version other than `<x.y.z>`
   whose release has not been revoked, which
   `gh release view v<a.b.c> --repo mustardcustarddev/detoo` shows.
   * If there is one, restore the formula to it with
     `git restore --source <commit> Formula/detoo.rb`, and check that
     `grep url Formula/detoo.rb` now shows `/releases/download/v<a.b.c>/`.
   * If there is none, `<version>` is the only release: remove the formula
     with `git rm Formula/detoo.rb`.
6. Commit the change with `git commit --all --message "Revoke detoo <x.y.z>"`,
   and push it with `git push`. Check what the tap now serves with
   `gh api --header 'Accept: application/vnd.github.raw' repos/mustardcustarddev/homebrew-tap/contents/Formula/detoo.rb | grep url`,
   which shows `/releases/download/v<a.b.c>/`, or, if you removed the
   formula, reports `Not Found`.
7. Delete the GitHub release, but keep its tag, with
   `gh release delete <version> --repo mustardcustarddev/detoo --yes`.
   Keeping the tag keeps the version used, so that the next release has a
   later one, which Homebrew upgrades to: Homebrew does not upgrade to a
   version that is already installed. Do not use `--cleanup-tag`, push the
   tag again, or run its Release workflow again, which would publish the
   release and its formula again.
8. Check that `gh release view <version> --repo mustardcustarddev/detoo`
   reports `release not found`, that
   `git ls-remote --tags origin <version>` still lists the tag, and that
   `gh release list --repo mustardcustarddev/detoo` marks `v<a.b.c>` as
   `Latest`, if there is one.
9. In the detoo repository, on `main`, write a notice for the documentation
   site in `docs/site/src/content/docs/releases/detoo-<x>-<y>-<z>-revoked.md`,
   such as `detoo-0-2-0-revoked.md`, with no dots in the name. Give it this
   front matter, where `<order>` is one less than the lowest `order` of the
   notices already in that directory, or `0` for the first, so that the
   newest is listed first:

   ```markdown
   ---
   title: detoo <x.y.z> is revoked
   description: detoo <x.y.z> is withdrawn. What was wrong with it, who it affects, and how to go back to an earlier release.
   section: releases
   order: <order>
   ---
   ```

   Under it, say what is wrong with the release, linking to the issue that
   reported it, who is affected, and what to do: run `brew update`, then
   `brew uninstall detoo`, then `brew install detoo`, which installs
   `<a.b.c>`, or, if you removed the formula, that detoo cannot be installed
   until the fix is released.
10. Add a warning to the top of `README.md`, under the badges, that links to
    the notice:

    ```markdown
    > [!WARNING]
    > detoo <x.y.z> is revoked. See
    > [detoo <x.y.z> is revoked](https://mustardcustarddev.github.io/detoo/docs/releases/detoo-<x>-<y>-<z>-revoked/)
    > for what to do if you installed it.
    ```

    Open question: when the warning is taken out of the README is not yet
    agreed.
11. Commit the notice and the warning with
    `git commit --message "docs: revoke detoo <x.y.z>"`, and push `main` with
    `git push`. Find the Site workflow's run with
    `gh run list --workflow site.yml --limit 1`, and wait for it to succeed
    with `gh run watch <run-id> --exit-status`. Check that the README on
    `https://github.com/mustardcustarddev/detoo` shows the warning, and that
    its link opens the notice, listed under Releases in the sidebar.
12. On an Apple silicon Mac with the revoked release installed, update the tap
    with `brew update`, then remove the release with `brew uninstall detoo`.
    * If the tap still has a formula, install the earlier release with
      `brew install detoo`, then run `detoo doctor` in a project, and check
      that its header shows `<a.b.c>` beside the detoo badge.
    * If you removed the formula, check that `brew info detoo` reports
      `No available formula`.
13. Release the fix, when it is ready, with
    [SOP 1](SOP-0001-release-detoo.md). Its version is later than
    `<version>`, because the tag is kept.
