# Writing a PIOSEE playbook

A playbook is read by someone who has just met the situation it plans for. A
customer may be waiting, a deadline may be close, and the author may not be
there to ask. They will do what it says. So every source it names must exist,
every step must work, and every rule must be the team's. A log name nobody
checked or a menu path that is not there wastes their time when they have
least of it. A target nobody agreed, or a default nobody chose, does worse:
it makes the team's policy something it is not.

So the plan is the user's: the scenario, the decision, the options, the
selection rules, the policies and the targets. You find the concrete facts
you can, and ask for the rest. See [piosee.md](piosee.md) for what each
section is for.

## Facts first

Before writing any prose, build a fact sheet: every claim the playbook will
make, and where it came from. Keep it in your working notes. Show it to the
user with the draft. Do not write it into the playbook.

| #   | Claim                                                                  | Source                                       |
| --- | ---------------------------------------------------------------------- | -------------------------------------------- |
| 1   | Charges are created by `CheckoutService`, retried once on a timeout    | `checkout/service.go`, lines 40–72, read      |
| 2   | Each charge is logged with the order ID and the provider's charge ID   | `checkout/service.go`, line 88, read          |
| 3   | Provider's refund guide: the link in `docs/payments.md`                | `docs/payments.md`, read; page read           |
| 4   | Scenario: a customer reports two charges for one order                 | User, in this conversation                    |
| 5   | Options: refund the duplicate; explain two orders; escalate            | User, in this conversation                    |
| 6   | Default when the evidence is unclear: escalate to engineering          | User, in this conversation                    |
| 7   | Support may refund up to the order's value without approval            | User, in this conversation                    |
| 8   | Reason code for the refund                                             | **Unknown — ask**                             |
| 9   | Target time to resolve                                                 | **Unknown — ask**                             |

### Acceptable sources

* **The user**, in this conversation: what they stated, or explicitly agreed
  to. A suggestion of yours that they did not take up is not a fact.
* **The repository**: code, configuration, runbooks, policies, documents and
  earlier playbooks that you have read. Note the path, and the line where it
  matters.
* **Commands you ran**, and their output, such as the list of a service's
  log sources or the fields of a stored record.
* **A document the user pointed you to**, such as a policy, a support macro,
  an incident ticket, or a third party's documentation. Cite a link only if
  you have read the page, or the user gave it to you.

### Not acceptable

Do not write any of these unless a source above supports it:

* an identifier: a service, log, path, field, event, table, setting or
  environment;
* a step through a user interface or a command, and the order of its menus
  and buttons;
* a link;
* a policy, an authority limit, an owner, a reason code or a message to a
  customer;
* an option, when it applies, or its tradeoffs;
* a selection rule, a default, a confidence level or an assumption;
* a deadline, an SLA, a threshold, a target or a review date;
* a claim about cause: that a situation "is caused by" something.

A cause you inferred from the code is a hypothesis, not a fact. Put it to the
user as one: "Payments are retried on a timeout. Could that charge a customer
twice?" If they agree, the playbook may say what to check to confirm it,
which is an Information step, rather than state the cause.

### Your own suggestions

The user may ask for your help: "what options are there?", "what should we
measure?", "how do we do a refund in the dashboard?". That is a fair request.
Answer it in the conversation, labelled as your suggestion, grounded in what
you read.

A suggestion enters the playbook only when the user adopts it. It is then
theirs, and the fact sheet says so: "User adopted the agent's suggestion, in
this conversation". Steps you drew from a third party's documentation still
need the user to confirm that their team does it that way.

Do not volunteer suggestions in place of asking. Ask the open question first.
Offer ideas only if the user asks for them, or says they do not know and
would like some.

### When a fact is missing

Ask the user. If they do not know, or prefer not to say, do one of these, and
let them choose:

* **Mark it as an open question** in the playbook, where the reader will meet
  it: "Open question: the reason code for this refund is not yet agreed." The
  team can settle it before the playbook is agreed.
* **Find it out**: say who could be asked, or what could be read.
* **Leave it out**, where the template does not need it.

Never fill the gap with something plausible.

Three gaps cannot be filled this way. If the user cannot say **what the
situation is**, **what has to be decided**, or **what can be done about it**,
there is no plan to write yet. Tell them so.

## The sections

The built-in template has a title, a date, a status, the Scenario, and one
section for each step. A project's own template may differ. Fill in each of
its sections by the guidance written in it, with the same standard of
evidence. Follow the house style of the existing playbooks.

### Title

A short phrase naming the situation, not a decision. For example: "Unexpected
charge after account deletion", "Customer cannot cancel subscription",
"Crossing VAT registration threshold". Run `pse list`, and follow the style
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

### Information

The facts the person handling it must collect before deciding. This is a
list of things to find out, not a list of findings.

* Group them, numbered, as the examples do: what the reporter says, the
  internal evidence, the third party's evidence, the policy that constrains
  the decision.
* Under each, name exactly what to collect, and where: the service and the
  log, the storage path and the fields, the dashboard page, the document.
* Name who owns each source, or can get at it, if the user said.
* Mark what is assumed rather than checked.
* Name the critical unknowns, and how to settle them quickly.

### Options

The realistic courses of action, numbered, each a distinct action.

* For each: a short name, when it applies, and its tradeoffs, in the user's
  words. Keep their weight: "usually" stays "usually".
* Include the conservative fallback the user named, such as escalating,
  pausing or asking for approval.
* Mark the default only if the user named one.
* One option is enough when that is all the team has. Its fallback is then
  what to do when it does not apply.

### Selection

How to choose an option when the scenario arises. This is a rule, not a
choice already made.

* The criteria, as the user gave them.
* Which findings from Information lead to which option: "Select **Option 1**
  if either is true: …". Check that the rules leave no case uncovered and do
  not conflict. Raise any gap with the user, rather than closing it yourself.
* The default, if there is one, and when to depart from it.
* The confidence in the rule, and the key assumptions behind it, if the user
  gave them.

### Execution

The steps to carry out each option, in order, as instructions.

* Number the steps, with the details under each. Say who takes each step, by
  role.
* Be exact where the reader acts: the page, the menu path, the setting, the
  command. Write menu paths and identifiers in backticks: `Billing` →
  `Subscriptions`.
* What to record, and where: who acted, when (in UTC), the identifiers, and a
  reason code if the team uses one.
* Whom to tell, and what, with the wording if the user gave it. Long wording
  can go in its own section, such as **Communication Templates**.
* Dependencies, timings and checks, such as "T-1 business day" or "verify the
  resulting state".
* The rollback path, if a step goes wrong.
* What to do so that it does not happen again, such as opening a bug with
  given details, if the user wants that.

### Evaluation

How the team will know, after the playbook has been used, whether it works.

* What to measure, and the targets, as the user set them.
* When it is reviewed, and by whom, by role.
* What happens if the targets are missed.

If there is nothing to say under a heading the template does not require,
leave the heading out. Add extra sections, such as **Communication
Templates** or **Reference Links**, only with the user's agreement and from
evidence.

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
* Instead of the Information step "Check the logs", write the one the user
  confirmed: "`CheckoutService` logs: the charge attempts for the order ID,
  with the provider's charge ID for each."
* Instead of adding "Do nothing" as an option because playbooks have one, ask
  what the team would do. Write theirs.
* Instead of the selection rule "Refund if it seems like a duplicate", ask
  what evidence shows a duplicate, and write that: "Select Option 1 if both
  charges carry the same order ID."
* Instead of "Resolve 95% within one business day" because it sounds right,
  ask what the target is. If there is none yet, mark it as an open question.
* Instead of a dashboard path from memory, read the provider's documentation,
  or ask the user, and write the path they confirm.

## Before you finish

* Every sentence traces to a line of the fact sheet.
* The playbook is a plan for every occurrence of the situation. It records no
  incident, and no decision made at the time.
* The scenario, the decision, the options, the selection rules, the policies
  and the targets all come from the user.
* Every identifier, step and link was read from a source, or given by the
  user.
* There is a conservative fallback, named by the user.
* The selection rules cover every case, or the gaps are open questions.
* Every "Unknown — ask" line was asked about. Its answer is used, or it is an
  open question, or it was left out, by the user's choice.
* Every section of the template in use is filled in, and no guidance text is
  left.
* The status is `Draft`.
* Links between playbooks were made by the tool, not by hand.
* The user has seen the draft with its sources, and approved it.
