# Contributing to detoo

The [README](README.md#for-developers) says how to build the tools from
source.

## Release

```bash
make release                # tag the next version, and push it
make release VERSION=v1.0.0 # or choose the version yourself
```

`make release` builds and tests the tools, and works out the next version with
[svu][] from the [conventional commits][] since the last release: a `fix:`
makes a patch release, a `feat:` a minor one, and a breaking change, such as
`feat!:`, a major one. It then tags `main` with the version, and pushes the
tag. The tag starts the [Release workflow](.github/workflows/release.yml),
which builds the tools for Apple silicon Macs, publishes them as a GitHub
release, tests the Homebrew formula, and publishes it to
[mustardcustarddev/homebrew-tap][tap].

The workflow pushes to the tap with a deploy key, which is set up once:

```bash
ssh-keygen -t ed25519 -N "" -C "detoo release" -f homebrew-tap
gh repo deploy-key add homebrew-tap.pub --repo mustardcustarddev/homebrew-tap --title "detoo release" --allow-write
gh secret set HOMEBREW_TAP_DEPLOY_KEY --repo mustardcustarddev/detoo < homebrew-tap
rm homebrew-tap homebrew-tap.pub
```

GitHub removes a deploy key that `gh` added if `gh`'s authorisation is
revoked. To keep the key regardless, add `homebrew-tap.pub` in the tap's
Settings › Deploy keys instead, with write access.

The formula is [packaging/homebrew/detoo.rb](packaging/homebrew/detoo.rb), which
the workflow fills in with each release's version and checksum.

[svu]: https://github.com/caarlos0/svu
[conventional commits]: https://www.conventionalcommits.org/
[tap]: https://github.com/mustardcustarddev/homebrew-tap
