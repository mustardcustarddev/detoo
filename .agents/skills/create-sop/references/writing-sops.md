# Writing a standard operating procedure

An SOP is followed step by step by someone who may never have done the task,
perhaps in the middle of an incident, because a playbook told them to. They
will type the commands it gives, in the order it gives them. So every
command, flag, path, page and environment name must be one that exists, and
every step must be in the right place.

That detail is in the repository and in the documentation of the tools
involved far more often than it is in anyone's memory. So research it, write
it with its sources, and ask the user only for what you cannot find.

## Where to look

Look for how the project does the task, in roughly this order, and stop
looking for a step once you have a source for it that runs, rather than one
that describes:

1. **What runs.** Build files and task runners: `Makefile`, `justfile`,
   `Taskfile.yml`, `package.json` scripts, `scripts/` and `bin/`. CI and CD
   workflows: `.github/workflows/`, `.gitlab-ci.yml`, `Jenkinsfile` and the
   like. Deployment configuration: Dockerfiles, Compose files, Kubernetes
   manifests, Helm charts, Terraform, Ansible, Procfiles and platform files
   such as `fly.toml` or `app.yaml`.
2. **What is written down.** `README.md`, `CONTRIBUTING.md`, `docs/`,
   runbooks, existing SOPs and playbooks, ADRs that constrain the task, and
   comments beside the code that does it.
3. **What the tools say.** The official documentation of each tool and
   service involved, matched to the version the project uses, as pinned in a
   lock file, a manifest or configuration. A tool's own `--help`, where it
   is installed.
4. **What git says.** `git log` for the files involved, to see how recently
   they changed, and whether a README describes a procedure the
   configuration has since left behind.

Run read-only commands only. Never run a step of the procedure to see what
it does.

## Facts first

Before writing the steps, build a fact sheet in your working notes: every
step and claim, and where it came from. Show the sources with the draft. Do
not write the fact sheet into the SOP.

| #   | Step or claim                                                    | Source                                             | Status   |
| --- | ---------------------------------------------------------------- | -------------------------------------------------- | -------- |
| 1   | Releases are tagged `vX.Y.Z` and deployed by the `deploy` job    | `.github/workflows/release.yml`, lines 1–40, read  | Found    |
| 2   | `make deploy VERSION=…` deploys a given version                  | `Makefile`, lines 52–60, read                      | Found    |
| 3   | The deployed version is shown by `kubectl rollout history`       | Kubernetes documentation, for the pinned version   | Found    |
| 4   | README says to deploy with `./deploy.sh`, which no longer exists | `README.md`, line 88; `git log -- deploy.sh`       | Conflict |
| 5   | Who may deploy to production                                     | **Unknown**                                        |          |

The status is **Found** when a source you read gives it, **Conflict** when
two sources disagree, **Given** when the user said it, and **Unknown** when
nothing settles it. Found and Given entries go into the SOP as steps. A
Conflict is put to the user, or, if they leave it, becomes an open question.
An Unknown becomes an open question.

### Acceptable sources

* **The repository**: code, scripts, configuration, documents, runbooks,
  SOPs and playbooks you have read. Note the path, and the lines where it
  matters.
* **The official documentation** of a tool or service involved, for the
  version the project uses, which you have read. Cite a link only if you have
  read the page.
* **The output of a read-only command you ran**, such as `--help`.
* **The user**, in this conversation, or a document they pointed you to.

### Not acceptable

Do not write any of these unless a source above gives it:

* a command, a flag, an argument or its order;
* a path, a file, a service, a cluster, a namespace, a bucket, an
  environment or a branch;
* a step through a user interface, and the order of its menus and buttons;
* a link;
* who may do a step, or who must approve it;
* a check that a step worked, or what its output should be;
* a time a step takes, or must be done within.

A step you remember from another project, or think is how such tools are
usually used, is not a source.

## When a fact is missing

Do not stop to interview the user for each gap. Write the step as an open
question, where the reader will meet it, and carry on:

```markdown
3. Open question: which cluster serves production. The repository names
   `staging` in `deploy/staging.yaml`, but has no production configuration.
```

Ask the user, in the one question before creating the SOP, only about the
gaps that matter most: a conflict between sources, or a fact without which
the procedure cannot be followed at all. At most three. Leave the rest as
open questions for the team.

If research finds nothing at all for the task, there is nothing to draft.
Say so, with where you looked, and ask the user how to proceed.

## The title

The task, as an instruction, short enough to be a file name: "Roll back a
release", "Rotate the webhook signing secret", "Restore the database from a
backup". Follow the style of the existing SOPs.

## The steps

An ordered list, under the template's heading for them.

* **One action a step**, as an instruction in the imperative: "Find the
  previous release tag", not "The previous tag should be found".
* **Exact.** Put the command in a code block, or in backticks, as the
  reader will type it, with placeholders in angle brackets for what changes
  each time: `make deploy VERSION=<tag>`. Write menu paths and identifiers
  in backticks: `Settings` → `Webhooks`.
* **Preconditions first.** If the task needs access, a tool or an approval,
  the first steps say so: "Make sure you can run `kubectl` against the
  production cluster."
* **A check after each step that can fail**, when a source gives one: how to
  tell it worked, such as the output to expect, and what to do if not.
* **No decisions.** "If the check fails, stop, and escalate to the on-call
  engineer" is a check. "Decide whether to roll back or fix forward" is a
  decision, and belongs in a playbook.
* **Sub-steps** for the detail of a step, nested under it, when they keep the
  main list readable.
* **Links** to the documentation the reader needs at a step, only to pages
  you have read.
* **No history.** The steps say how to do the task now, not how it used to
  be done, or why it changed.

Do not add sections the template does not have unless the user asks, or the
house style of the existing SOPs uses them.

## Examples

* Instead of "Deploy the previous version", write the command the
  repository uses: "Deploy the previous release with
  `make deploy VERSION=<tag>`, where `<tag>` is the tag found in step 2."
* Instead of "Check that it worked", write the check a source gives: "Run
  `kubectl rollout status deployment/web`, and wait for
  `successfully rolled out`."
* Instead of `./deploy.sh`, which the README names but the repository no
  longer has, use what the CI workflow runs, and tell the user the README is
  out of date.
* Instead of asking "How do you roll back?", read the Makefile and the CI
  workflow, and propose the steps.
* Instead of "Production cluster: `prod-eu-1`" because it seems likely,
  leave an open question, and ask if it matters most.

## Before you finish

* Every step traces to a line of the fact sheet.
* Every command, path, page and name was read from a source, or given by the
  user.
* No step was run to test it.
* There is no decision in the SOP; a check says only to stop and escalate.
* Every Unknown and Conflict was asked about or is an open question.
* No guidance text from the template is left.
* The status is `Draft`.
* Links between SOPs were made by the tool, not by hand.
* The user has seen the draft with its sources, and chose to create it.
