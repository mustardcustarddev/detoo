# Interviewing, researching and confirming

A request such as "write a playbook for when a Stripe payment is handled
incorrectly" names an area. A playbook needs more:

* how the situation shows itself;
* the decision it calls for, and the risks of getting that decision wrong;
* what must be found out, where, and from whom;
* the options, and when each applies;
* the rule for choosing between them;
* the steps for each option;
* how the team will know it worked.

The first two are the user's: only they know which situation they mean, and
what they have to decide. The rest can usually be found in the code, the
project's documents and the documentation of the systems involved. So you
interview for the problem, research the rest, and ask the user to confirm
what you found. You interview step by step only for what research leaves
open.

## Principles

* **It is a plan for next time.** Ask about the situation in general: "When a
  customer reports this, what do they usually tell you?" If the user answers
  with one incident, take it as an example of the situation. Ask what would be
  the same next time, and what might differ.
* **The template drives the work.** Read the template in use afresh for
  every playbook. Each section is a step the person following the playbook
  will expect to find, and each needs research, a confirmation or both.
* **Interview for the problem first.** What happens, how it shows itself,
  which systems are involved, and what has to be decided. Research aimed at
  the wrong situation wastes the user's time later.
* **Research before you ask.** Learn what the code, the configuration, the
  project's documents and the systems' documentation can tell you. Then the
  user is asked only to confirm, and to fill the gaps.
* **Propose with sources.** Every proposal says where it came from: a file
  and line, a document, a page of the provider's documentation. The user is
  confirming your reading of those sources, and needs to see them to judge it.
* **Confirm, don't quiz.** A proposal is put as a statement to confirm, with
  the options Yes and No, and a correction typed in. Do not turn what you
  found into an open question the user has to answer from scratch.
* **Policy is the team's.** A proposal may state a limit, an owner, a default
  or a target only when a policy, a contract or the system's documentation
  states it. Otherwise ask, openly.
* **Evidence may prompt; it may not answer.** You may put an observation to
  the user as a question: "Payments are retried on a timeout. Could a retry
  cause a second charge?" Their answer is the evidence. A cause you inferred
  is not a fact until they confirm it.
* **Ask for the fallback.** Every playbook needs a safe course for when the
  evidence is unclear, or time runs out. Propose one if a source gives one.
  Otherwise ask for it. Do not supply it.
* **Small batches.** Up to four related questions in one call, then wait.
* **Follow up vague answers.** "We'd check the logs" is not an Information
  step. Ask which logs, for what, and who can read them, or find out. Do not
  sharpen it yourself.

Ask every question through your question tool, as [asking.md](asking.md)
describes. Ask in the chat only if you have no such tool.

## Interviewing for the problem

Ask in one or two rounds, openly, citing anything you already know:

* What happens, or is reported, that calls for this playbook?
* How does the team come to know of it: a customer contact, an alert, a
  report?
* Which systems are involved, as far as you know: the provider, the services,
  the tools support uses?
* When it happens, what exactly has to be decided? What goes wrong if it is
  decided wrongly, either way?
* Is there an incident, a ticket or a document I should read first?

Then play the Scenario and the Problem back, and confirm them, as below.

## Researching

Work through the playbook's concerns in order, and record each finding in the
fact sheet with its source. The workflow in
[SKILL.md](../SKILL.md#5-research) says what to look for under each.

1. **The context of the problem**: the code and documents behind the
   situation.
2. **The course of investigation**: where the evidence lives, and how to
   reach it.
3. **The options and the execution**: what the systems' documentation and the
   project's own tools allow, and how.
4. **Selection and Evaluation**: what decides between the options, and how
   success is measured, where any source says.

For the documentation of a third-party system:

* Read the official documentation: the provider's own guides and API
  reference, not a blog post about them.
* Match it to what the project uses: the API version pinned in configuration
  or a client library's version, the products it uses, the dashboard or the
  API, whichever the team will act through.
* Note each page you read, with its link, for the fact sheet and the
  **Reference Links** section.
* Note the constraints the documentation sets, such as time limits, actions
  that cannot be undone, and conditions under which an action is refused.

## Confirming a proposal

Write the proposal in the chat first: the section as it would read in the
playbook, then its sources. Then ask, through the question tool:

* **The question**: "Is this the Information section?", or a play-back such
  as "So: support cannot see why a payment failed, and must decide whether to
  retry it or refund. Is that right?"
* **Options**: **Yes** ("Use it as it is") and **No** ("Something is wrong
  with it"). Say in the question that the user can choose Other to correct it
  in their own words.

Confirm up to four sections in one call, one question each. Where a section
holds a real choice, such as which of four documented options the team would
use, ask that as one multiple-choice question, with an option for each and
Other for one not listed.

Then:

* **Yes**: mark the entry confirmed in the fact sheet.
* **Other**: use the correction. Play it back if it changes more than a
  detail.
* **No**: ask what is wrong, openly, then research again or interview for the
  section.

## Question bank

For what research did not settle. Ask these step by step, specific to the
situation, and cite what you did find. This bank is for the built-in
template; adapt it to the template in use.

### Scenario

* What happens, or is reported, that calls for this playbook?
* How often does it happen, or when is it expected? Has it happened before?
  Where can I read about it?
* What standing conditions frame it, such as which customers, which markets
  or which systems are involved?
* Is there more than one situation here? Should they be one playbook or
  several?

### Problem

* When this happens, what exactly has to be decided?
* What goes wrong if that decision is made wrongly? And if it is made wrongly
  the other way?
* What goes wrong if it is made late, or inconsistently from one case to the
  next?

### Information

* I could not find where this is recorded. Where would the person handling it
  look?
* Who can get at each source?
* What do we usually assume, rather than check? Should it be checked?
* What might not be knowable in time? How could it be settled quickly?
* Is there a policy, a limit or an agreement that constrains the decision?

### Options

* What can the team do in this situation, beyond what I found?
* For each: what does it cost, and what does it risk? Can it be undone?
* What is the safe course when the evidence is unclear, or time runs out?

### Selection

* How should the person handling it choose? What matters most?
* Which findings point to which option? Is there a case the rules leave
  uncovered, or one where two rules both apply?
* Is there a default? When should it not be used?
* Up to what amount, or in what time, may they act without approval?

### Execution

* For each option: who takes each step, by role?
* What does each step need: access, approval, another team?
* What must be recorded, and where? Do you use reason codes?
* Who must be told, and what do they need to hear? Is there wording to reuse?
* What should be done so that this does not happen again, and who owns that?

### Evaluation

* How will you know the playbook is working? What would you measure?
* What targets would count as success?
* Who reviews it, and how often, or after how many uses?
* If the targets are missed, what happens next?

When the user returns after the playbook has been used:

* What happened, in outline? Where can I read about it?
* Did the playbook lead to a sound decision, in time?
* Was any information missing, hard to get, or wrong?
* Was there an option it did not list, or a case its rules did not cover?
* Were any steps unclear, out of date, or in the wrong order?
* Should it be corrected in place, replaced, or updated by a new playbook?

## Worked example

The user says: "Write a PIOSEE for when a Stripe payment is handled
incorrectly."

### Interview for the problem

"What goes wrong with the payment, and how does the team find out? What does
support have to decide when it happens?" Suppose the user says: a customer
writes in that they were charged, but their order still shows as unpaid.
Support must decide whether to complete the order or refund the payment. The
risks are refunding a customer who should have received the order, and
shipping an order that was never paid.

Play back, and confirm: "So: a customer reports that they were charged, but
their order shows as unpaid, and support must decide whether to complete the
order or refund. Is that right?" Yes, No, or Other.

### Research

Suppose reading the code finds:

* checkout creates a Stripe PaymentIntent, with the order ID in its
  metadata, in `checkout/payments.go`;
* the order is marked paid only by the webhook handler for
  `payment_intent.succeeded`, in `webhooks/stripe.go`, which logs each event
  ID it handles;
* the Stripe API version is pinned in `config/stripe.yaml`;
* an admin endpoint marks an order paid, given a PaymentIntent ID, in
  `admin/orders.go`;
* there is no refund policy in the repository.

Suppose reading Stripe's documentation, for that API version, covers how to
find a PaymentIntent and its events, how a webhook event that failed can be
sent again, and how to create a refund. Note the link of each page you read.

### Confirm

Write each section as a proposal, with its sources, then confirm up to four
in one call:

* **Information**: the PaymentIntent for the order, found by the order ID in
  its metadata; whether it succeeded; whether the `payment_intent.succeeded`
  event reached the webhook handler, by its log; the order's status. Sources:
  the files above, and Stripe's pages on PaymentIntents and events.
* **Options**: 1. send the event again, so the handler marks the order paid;
  2. mark the order paid with the admin endpoint; 3. refund the payment.
  Sources: Stripe's pages on resending events and refunds, and
  `admin/orders.go`.
* **Selection**: send the event again when the PaymentIntent succeeded and the
  handler never received its event; mark the order paid when the handler
  received it but failed. Source: derived from the code's handling of the
  event. When to refund instead is policy, which no source states, so it is
  left for the interview.
* **Execution**: the steps for each option, from the documentation and the
  code.

Questions not to ask: "What should the Information section say?", or "How do
refunds work?". Research answers those.

### Interview for the rest

Research found no refund policy, so ask: "When should support refund rather
than complete the order? Is there a limit on the refund they may give without
approval, and who approves one above it?" Ask for the
fallback, the owners by role and the targets the same way, unless a source
gave them.

### Then draft

Draft from the confirmed fact sheet, show the sources, and ask for approval.
The playbook is for every report of a payment taken for an unpaid order. It
does not mention the customer whose report prompted it.
