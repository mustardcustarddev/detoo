# Writing a problem playbook

A playbook is read by someone who has just met the problem it plans for. A
system may be down, a customer may be waiting, a deadline may be close, and
the author may not be there to ask. They will do what it says, starting with
the immediate action. So every source it names must exist,
every step must work, and every rule must be the team's. A log name nobody
checked or a menu path that is not there wastes their time when they have
least of it. A target nobody agreed, or a default nobody chose, does worse:
it makes the team's policy something it is not.

So the scenario and the decision come from the user. The rest you research
in the code, the project's documents and SOPs, and the documentation of the
systems involved, and propose with its sources, for the user to confirm. Policies and
targets that no source states, you ask for. See [piosee.md](piosee.md) for
what each section is for.

## Facts first

Before writing any prose, build a fact sheet: every claim the playbook will
make, where it came from, and whether the user has confirmed it. Keep it in
your working notes. Show it to the user with the draft. Do not write it into
the playbook.

| #   | Claim                                                                | Source                                                  | Status    |
| --- | -------------------------------------------------------------------- | ------------------------------------------------------- | --------- |
| 1   | Scenario: a customer reports two charges for one order               | User, in this conversation                              | Given     |
| 2   | Charges are created by `CheckoutService`, retried once on a timeout  | `checkout/service.go`, lines 40–72, read                | Confirmed |
| 3   | Each charge is logged with the order ID and the provider's charge ID | `checkout/service.go`, line 88, read                    | Confirmed |
| 4   | Immediate action: none indicated; nothing worsens while it waits     | User, in this conversation                              | Given     |
| 5   | Counter-indicated: retrying the charge, which charges again          | Provider's idempotency guide, read                      | Confirmed |
| 6   | Find the charges by the order ID with SOP 2                          | `docs/sop/SOP-0002-find-a-customers-charges.md`, read   | Confirmed |
| 7   | Option: refund the duplicate charge, through the dashboard           | Provider's refund guide, linked from `docs/payments.md` | Confirmed |
| 8   | Option: escalate to engineering when the charges cannot be matched   | User, in this conversation                              | Given     |
| 9   | Support may refund up to the order's value without approval          | User, in this conversation                              | Given     |
| 10  | A refund can be made only within the period the provider states     | Provider's refund guide, read                           | Proposed  |
| 11  | Reason code for the refund                                           | **Unknown — ask**                                       |           |
| 12  | Target time to resolve                                               | **Unknown — ask**                                       |           |

The status is **Given** when the user said it, **Proposed** when you found it
and have yet to put it to the user, and **Confirmed** when they agreed to it.
Only given and confirmed entries go into the playbook.

### Acceptable sources

* **The user**, in this conversation: what they stated, or explicitly agreed
  to. A proposal of yours that they did not confirm is not a fact.
* **The repository**: code, configuration, SOPs, runbooks, policies,
  documents and earlier playbooks that you have read. Note the path, and the line where it
  matters.
* **Commands you ran**, and their output, such as the list of a service's
  log sources or the fields of a stored record.
* **A document the user pointed you to**, such as a policy, a support macro
  or an incident ticket.
* **The documentation of a system involved**, such as a payment provider's
  guides and API reference, which you have read. Use the official
  documentation, for the version the project uses. Cite a link only if you
  have read the page, or the user gave it to you.

### Not acceptable

Do not write any of these unless a source above supports it, and the user has
given or confirmed it:

* an identifier: a service, log, path, field, event, table, setting or
  environment;
* an immediate action, or that none is indicated;
* a counter-indicated action, or why it is wrong;
* a reference to an SOP that you have not read, or that does not exist;
* a step through a user interface or a command, and the order of its menus
  and buttons;
* a link;
* a policy, an authority limit, an owner, a reason code or a message to a
  customer;
* an option, when it applies, or its tradeoffs;
* a selection rule, a default, a confidence level or an assumption;
* a deadline, an SLA, a threshold, a target or a review date;
* a claim about cause: that a situation "is caused by" something.

What you remember of a system, or believe is common practice, is not a
source. Read the documentation, or ask.

A cause you inferred from the code is a hypothesis, not a fact. Put it to the
user as one: "Payments are retried on a timeout. Could that charge a customer
twice?" If they agree, the playbook may say what to check to confirm it,
which is an Information step, rather than state the cause.

### Proposals

Research is meant to find the investigation, the options and the steps, so
propose them. Write each proposal as it would read in the playbook, with the
source of each point. It enters the playbook only when the user confirms it.
Steps drawn from a third party's documentation still need the user to
confirm that their team does it that way, such as through the dashboard or
through the API.

A proposal states a policy, such as a limit, an owner, a default or a target,
only when a policy, a contract or the system's documentation states it. Say
which. Otherwise ask.

An idea of your own with no source is not a proposal. Offer one only if the
user asks for ideas, labelled as yours, and write it in only if they adopt
it.

### When a fact is missing

Ask the user. If they do not know, or prefer not to say, do one of these, and
let them choose:

* **Research it further**, if there is somewhere you have not looked, such
  as a page of the documentation or a part of the code.
* **Mark it as an open question** in the playbook, where the reader will meet
  it: "Open question: the reason code for this refund is not yet agreed." The
  team can settle it before the playbook is agreed.
* **Find it out**: say who could be asked, or what could be read.
* **Leave it out**, where the template does not need it.

Never fill the gap with something plausible.

Three gaps cannot be filled this way. If the user cannot say **what the
problem is** or **what has to be decided**, or no one can find **what can
be done about it**, neither the user nor your research, there is no plan to
write yet. Tell them so.

## The sections

The built-in template has a title, a date, a status, the Scenario, and one
section for each step: Problem, Immediate action, Counter-indicated actions,
Information, Options, Selection, Execution and Evaluation. A project's own
template may differ, or follow another method. Fill in each of
its sections by the guidance written in it, with the same standard of
evidence. Follow the house style of the existing playbooks.

### Title

A short phrase naming the problem, not a decision. For example: "Unexpected
charge after account deletion", "Customer cannot cancel subscription",
"Crossing VAT registration threshold". Run `plm list`, and follow the style
of the existing titles. The title becomes the file name, so keep it short.

### Date

The tool writes today's date. Leave it.

### Status

Always `Draft`. The tool writes it; leave it. Whether the playbook is agreed
is for the team to record.

### Scenario

The triggering situation, in one to three sentences, in the present tense,
for any occurrence of it: "A customer contacts support and says their bank
shows two charges for one order." Add the standing conditions that frame it,
when they matter to the plan: "Payments are taken through one provider, in
one currency."

* Describe the situation, not one instance of it: no customer, date or ticket.
* Describe what is observed, not its cause.

### Problem

The exact decision to make, then the risks.

* One clear question, or one statement of the decision: "Whether to refund
  one of the charges, and how, while keeping an audit trail."
* "Primary risks:" as a short list of what goes wrong if the decision is
  wrong. Include each direction, such as refunding when it should not and not
  refunding when it should, and being late or inconsistent, if the user named
  them.
* No solutions. "We need to refund" is an option, not a problem.

### Immediate action

What the person on hand does at once, before gathering information, to make
the situation safe or stop it getting worse.

* One or a few steps, as instructions, with who takes them, by role: "Turn
  off the `checkout_v2` feature flag in `Settings` → `Flags`."
* Refer to an SOP when the action is one: "Put the shop into maintenance
  mode by [SOP 4. Enable maintenance mode](../sop/SOP-0004-enable-maintenance-mode.md)."
* It must be right whatever the investigation finds, safe, and best
  reversible. An action that is right only in some cases is an option.
* Where nothing has to be done at once, write exactly: "No immediate action
  is indicated." Write it only when the user has said so, or confirmed it.
* Do not leave the section out, or empty: the reader must be able to see
  that nothing is needed, rather than wonder whether it was forgotten.

### Counter-indicated actions

Optional. The actions known to make the situation worse, or to waste time.

* A list, each item the action, as the reader might be tempted to put it,
  and why it is wrong: "Restarting the worker. It clears the queue, and with
  it the evidence of which jobs failed."
* Only actions with a source: the documentation, the code, an SOP, an
  earlier incident, or the user. Not general advice, such as "Do not panic".
* Leave the section out, heading and all, when none are known.

### Information

The facts the person handling it must collect before deciding. This is a
list of things to find out, not a list of findings.

* Group them, numbered, as the examples do: what the reporter says, the
  internal evidence, the third party's evidence, the policy that constrains
  the decision.
* Under each, name exactly what to collect, and where: the service and the
  log, the storage path and the fields, the dashboard page, the document. A
  command that shows it, such as `brew list --versions <formula>` for the
  version a reporter has installed, saves the reader looking it up.
  Where an SOP sets out how to collect it, refer to the SOP rather than
  repeating its steps.
* Name who owns each source, or can get at it, if the user said.
* Mark what is assumed rather than checked.
* Name the critical unknowns, and how to settle them quickly.

### Options

The realistic courses of action, numbered, each a distinct action.

* For each: a short name, when it applies, and its tradeoffs, as the sources
  and the user put them. Keep their weight: "usually" stays "usually".
* Include the constraints the documentation sets, such as a time limit, or
  that the action cannot be undone.
* Include a conservative fallback, such as escalating, pausing or asking for
  approval, from a source or from the user.
* Mark the default only if a source states it or the user named one.
* One option is enough when that is all the team has. Its fallback is then
  what to do when it does not apply.

### Selection

How to choose an option when the scenario arises. This is a rule, not a
choice already made.

* The criteria, as confirmed: from the sources, such as the conditions the
  documentation sets for each action, or from the user.
* Which findings from Information lead to which option: "Select **Option 1**
  if either is true: …". Check that the rules leave no case uncovered and do
  not conflict. Raise any gap with the user, rather than closing it yourself.
* The default, if there is one, and when to depart from it.
* The confidence in the rule, and the key assumptions behind it, as the user
  confirmed them.

### Execution

The steps to carry out each option, in order, as instructions.

* Number the steps, with the details under each. Say who takes each step, by
  role.
* Where a step is a routine task with an SOP, make the step a reference to
  it, with anything this playbook needs that the SOP does not say: "Roll
  back the release by [SOP 3. Roll back a release](../sop/SOP-0003-roll-back-a-release.md),
  to the version before the one that raised the alert." Do not copy its
  steps in.
* Be exact where the way matters: the page, the menu path, the setting or
  the command that makes a change only it can make. Write menu paths and
  identifiers in backticks: `Billing` → `Subscriptions`.
* Where any tool will do, as for replying to an issue, labelling it or
  sending a message, say what to do, and leave the tool to the reader: the
  team may use a browser, a command-line tool or an app.
* What to record, and where: who acted, when (in UTC), the identifiers, and a
  reason code if the team uses one.
* Whom to tell, and what, with the wording if a source or the user gave it.
  Long wording
  can go in its own section, such as **Communication Templates**.
* Dependencies, timings and checks, such as "T-1 business day" or "verify the
  resulting state".
* The rollback path, if a step goes wrong.
* What to do so that it does not happen again, such as opening a bug with
  given details, if the user wants that.

### Evaluation

How the team will know, after the playbook has been used, whether it works.

* What to measure, and the targets, from a service level or policy the
  project keeps, or as the user set them.
* How to check the result, referring to an SOP where one covers the check.
* When it is reviewed, and by whom, by role.
* What happens if the targets are missed.

If there is nothing to say under a heading the template does not require,
leave the heading out. Leave out **Counter-indicated actions** when none are known, but never
**Immediate action**. Add extra sections, such as **Communication
Templates** or **Reference Links**, only with the user's agreement and from
evidence. **Reference Links** to the documentation you read are usually
worth proposing.

## Style

* Plain words, in the team's terms. Use the system's own names where the
  reader must find something by them.
* Instructions in the imperative: "Confirm the customer's identity", not "The
  customer's identity was confirmed" or "We might confirm".
* Present tense for the scenario and for what is true of the systems. No past
  tense about an incident.
* Short items in lists, with the detail under them, so the reader can find
  their place quickly.
* No blame, and no individual's name: roles, not people.
* No placeholders and no template guidance left behind. The Communication
  Templates may hold placeholders such as `<amount>` for the reader to fill
  in.

## Examples

Each pair shows what to write instead, if the evidence supports it.

* Instead of the Scenario "On 3 March, the customer on ticket 4471 was charged
  twice", write "A customer reports two charges for one order", and leave the
  incident out.
* Instead of the Information step "Check the logs", find which logs record
  it, propose that, and write what the user confirmed: "`CheckoutService`
  logs: the charge attempts for the order ID, with the provider's charge ID
  for each."
* Instead of adding "Do nothing" as an option because playbooks have one,
  propose the options the documentation and the code support, and ask
  whether the team has others.
* Instead of the selection rule "Refund if it seems like a duplicate", find
  what evidence shows a duplicate, and propose that: "Select Option 1 if both
  charges carry the same order ID."
* Instead of asking "How do refunds work?", read the provider's refund
  documentation, and propose the steps.
* Instead of "Resolve 95% within one business day" because it sounds right,
  ask what the target is. If there is none yet, mark it as an open question.
* Instead of the immediate action "Investigate the issue", which is
  Information, find what stops the harm growing, and propose that, or
  propose "No immediate action is indicated." if nothing does.
* Instead of copying the ten steps of a rollback into Execution, refer to
  the SOP for it, or offer to write one.
* Instead of `gh issue comment 42 --body "..."` or
  `gh issue close 42 --reason "not planned"`, write "Reply on the issue with
  what was found, and close it." The reader handles issues with whatever tool
  they use. Keep a command where it finds a fact, such as
  `brew list --versions <formula>`.
* Instead of a dashboard path from memory, read the provider's documentation,
  or ask the user, and write the path they confirm.

## SOP references

A playbook refers to an SOP by a Markdown link to its file, relative to the
playbook, with its number and title as the link text:
`[SOP 3. Roll back a release](../sop/SOP-0003-roll-back-a-release.md)`.
Run `sop config` for the directory, and `sop list` for the file names.

* Refer only to an SOP you have read, and whose steps fit this use.
* Where the playbook needs a routine task that no SOP covers, ask the user
  whether to draft one with the `create-sop` skill, to write the steps into
  the playbook, or to mark it as an open question.
* Where an SOP is out of date or does not quite fit, say so, and offer to
  update it, rather than working round it in the playbook.

## Before you finish

* Every sentence traces to a line of the fact sheet.
* The playbook is a plan for every occurrence of the problem. It records no
  incident, and no decision made at the time.
* The scenario and the decision come from the user.
* There is an Immediate action section: an action given or confirmed by the
  user, safe whatever the evidence, or "No immediate action is indicated."
* Every counter-indicated action has a source, or the section is left out.
* Every SOP linked to exists at the path linked, and you have read it.
* Every option, selection rule and step was given by the user, or proposed
  from a source and confirmed by them.
* Every policy and target was given by the user, or stated by a source the
  user confirmed.
* Every identifier, step and link was read from a source, or given by the
  user.
* No step prescribes a tool for an action any tool will do, such as replying
  to an issue. Commands appear where they find a fact, check a result, or
  are the only way to act.
* The number in the title is the number in the file name.
* There is a conservative fallback, named by the user.
* The selection rules cover every case, or the gaps are open questions.
* No entry is still Proposed.
* Every "Unknown — ask" line was asked about. Its answer is used, or it is an
  open question, or it was left out, by the user's choice.
* Every section of the template in use is filled in, and no guidance text is
  left.
* The status is `Draft`.
* Links between playbooks were made by the tool, not by hand.
* The user has seen the draft with its sources, and approved it.
