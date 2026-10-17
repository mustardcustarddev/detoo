# Interviewing for a PIOSEE playbook

A request such as "write a playbook for when a customer can't cancel" names a
situation. A playbook needs more:

* how the situation shows itself;
* the decision it calls for, and the risks of getting that decision wrong;
* what must be found out, where, and from whom;
* the options, and when each applies;
* the rule for choosing between them;
* the steps for each option;
* how the team will know it worked.

The repository and the documents can tell you some of this, often the most
concrete parts. The rest is the user's, and their team's. This is how to draw
it out without leading them.

## Principles

* **It is a plan for next time.** Ask about the situation in general: "When a
  customer reports this, what do they usually tell you?" If the user answers
  with one incident, take it as an example of the situation. Ask what would be
  the same next time, and what might differ.
* **The template drives the questions.** Read the template in use afresh for
  every playbook. Each section is a step the person following the playbook
  will expect to find.
* **Read before you ask.** First learn what the code, the configuration and
  the documents can tell you. Then you ask only what you could not read, and
  can ask precisely: "Account deletion is handled by `UserDeletionService`,
  which writes an audit record to `deleted/{userId}/…`. Is that where support
  would look to confirm when an account was deleted?"
* **Facts are yours to find; the plan is the user's.** Show the user the
  sources you found. Which of them matter, what the options are, and how to
  choose are for them to say.
* **Open before closed.** Ask about the situation, the risks, the options and
  the criteria with open questions. Offer options only where the answers are
  known: the sources you found, the options already named, whether to record a
  point as an open question. An open question still goes through the question
  tool, with its answer typed in.
* **Never offer guessed answers as choices.** "Should the default be to
  refund, to refuse or to escalate?" asked before the user has named any
  options puts words in their mouth. "What can support do in this situation?"
  does not.
* **Evidence may prompt; it may not answer.** You may put an observation to
  the user and ask whether it matters: "Payments are retried on a timeout.
  Could a retry cause a second charge?" Their answer is the evidence.
* **Ask for the fallback.** Every playbook needs a safe course for when the
  evidence is unclear, or time runs out. If the user has not named one, ask
  for it. Do not supply it.
* **Ask for the numbers.** Thresholds, deadlines, limits and targets make a
  playbook decidable and measurable, and only the team can set them. Ask for
  them, such as "within what time?" or "refunds up to what amount?". Do not
  fill them in.
* **Small batches.** Three or four related questions, then wait.
* **Play back.** After each section, say it back in a sentence or two, so the
  user can correct it before the next section builds on it.
* **Follow up vague answers.** "We'd check the logs" is not an Information
  step. Ask which logs, for what, and who can read them. Do not sharpen it
  yourself.

Ask every question through your question tool, as [asking.md](asking.md)
describes:

* the open question in the question itself;
* options only for answers with evidence, and for the choices when a fact is
  missing;
* the answer typed in as free text.

Ask in the chat only if you have no such tool.

## Planning the interview

Before the first question, write a plan in your working notes:

1. List the sections of the template in use, with the guidance under each.
2. Under each, write the questions this situation raises, using the question
   bank below as prompts, not as a script.
3. Cross off what the fact sheet already answers. Turn what remains into
   specific questions, citing what you found where it helps.

## Question bank

This bank is for the built-in template. Adapt it to the template in use, and
to the situation.

### Scenario

* What happens, or is reported, that calls for this playbook?
* How does the team come to know of it: a customer contact, an alert, a date
  in the calendar?
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

* Here are the sources I found, and what each holds. Is any of it wrong, or
  out of date?
* What does the person handling this need to find out before deciding?
* Where is each fact found, and who can get at it?
* What do we usually assume, rather than check? Should it be checked?
* What might not be knowable in time? How could it be settled quickly?
* Is there a policy, a limit or an agreement that constrains the decision?

### Options

* What can the team do in this situation? What else? (Ask again before moving
  on.)
* When does each option apply?
* For each: what does it cost, and what does it risk? Can it be undone?
* What is the safe course when the evidence is unclear, or time runs out?
* Do you want me to suggest options too? (Only if they want; label them.)

### Selection

* How should the person handling it choose? What matters most?
* Which findings point to which option? Is there a case the rules leave
  uncovered, or one where two rules both apply?
* Is there a default? When should it not be used?
* How confident are you in this rule? What does it assume that could change?
* Do you want my view first? (Only if they ask; label it.)

### Execution

* For each option: what are the steps, in order? Who takes each, by role?
* What does each step need: access, approval, another team?
* Exactly where is it done? I found this in the documentation; is this how
  your team does it?
* What must be recorded, and where? Do you use reason codes?
* Who must be told, and what do they need to hear? Is there wording to reuse?
* If a step goes wrong, how is it undone?
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

The user says: "Write a PIOSEE for when a customer says they've been charged
twice."

### Facts to gather first

Read before asking:

* `pse config` and `pse list`, and any playbook about charges, refunds or
  payments.
* The code that takes payments: where a charge is created, whether it is
  retried, and what it logs.
* How the payment provider is used: the identifiers stored for a customer,
  and the provider documentation the project links to.
* Any policy on refunds, in the repository or wherever the user points you.

Suppose this finds the following:

* charges are created by `CheckoutService`, which retries once after a
  timeout, with an idempotency key built from the order ID;
* each charge is logged with the order ID and the provider's charge ID;
* the provider's documentation on refunds is linked from
  `docs/payments.md`;
* there is no refund policy in the repository.

### Questions not to ask

* "What should the playbook say?" (a form field; ask about the situation)
* "This is caused by the retry, isn't it?" (a guessed cause, as a leading
  question)
* "Should support refund, refuse or escalate?" before the user has named any
  options (offering your options in place of theirs)
* "Is the default to refund because it keeps customers happy?" (a guessed
  reason)

### The interview

**Scenario.** "When a customer says they've been charged twice, what do they
usually tell you, and how does it reach the team?" Then play back: "So: a
customer contacts support, saying their bank shows two charges for one order.
Is that the situation?"

**Problem.** "What has support to decide when this comes in? What goes wrong
if they decide it wrongly either way?" Suppose the user says: whether to
refund one charge; the risks are refunding a real second order, and leaving a
customer out of pocket.

**Information.** "Each charge is logged by `CheckoutService` with the order
ID and the provider's charge ID. Is that where support would look? What else
do they need, such as the bank statement, or whether there were two orders?"

**Options.** "What can support do?" Suppose they say: refund the duplicate,
or explain that there were two orders. "Anything else? What should they do
when they can't tell?" They add: escalate to engineering.

**Selection.** "How should support tell which applies? Is one the default?"
Write their rules, in their words. Ask what the rules miss, such as one charge
that succeeded and one still pending.

**Execution.** "Who refunds, and where? The provider's refund guide is linked
from `docs/payments.md`. Is that the flow you use? What should be recorded in
the ticket? Is there a reason code?"

**Evaluation.** "How will you know the playbook works? What would you count,
and who reviews it?"

### Then draft

Draft with what the user agreed in each section, show the sources, and ask
for approval. The playbook is for every duplicate-charge report. It does not
mention the customer whose report prompted it.
