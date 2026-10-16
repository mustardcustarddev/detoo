# Interviewing for a PIOSEE record

A request such as "help us decide whether to roll back the release" names a
situation. A record needs more: what exactly the problem is and how long
there is, what is known and what is not, the options and how the user weighs
them, which they choose and why, how it will be carried out, and how they will
know it worked. The systems can tell you some of the facts. The rest is the
user's, and their team's. This is how to get it out, without leading them,
and without slowing them down.

## Principles

* **Time first.** The first question is how long there is. It sets the pace of
  everything after it. If harm is happening and the action is obvious, the
  interview waits.
* **One step at a time, in order.** Each step depends on the one before. Do
  not ask for options before the problem is agreed, or for a choice before
  the options are laid out. Going back when something new comes to light is
  part of the method; jumping ahead is not.
* **The template drives the questions.** Read the template in use afresh for
  every record. Each section is a step the team will expect to find.
* **Read before you ask.** Learn what the repository and the systems can tell
  you first, so that you ask only what you could not read, and can ask
  precisely: "Errors rose at 14:05, three minutes after release 2.3 was
  deployed. Is anything else known to have changed?"
* **Facts are shared; judgements are the user's.** Show the user the facts
  you found, with their sources. Whether those facts explain the problem, and
  what to do about it, is for them to say.
* **Open before closed.** Ask about the problem, the options and the reasons
  with open questions. Offer options only where the answers are known: which
  of the options already named, whether to record a fact as unknown. An open
  question still goes through the question tool, with its answer typed in.
* **Never offer guessed reasons as choices.** "Did you choose it because it
  was fastest, safest or cheapest?" puts words in their mouth. "Why that one
  over the others?" does not.
* **Evidence may prompt; it may not answer.** You may put an observation to
  the user and ask whether it matters: "The release includes a migration. Does
  that affect rolling back?" Their answer is the evidence.
* **Small batches.** Under time pressure, one question at a time. Otherwise,
  three or four related questions, then wait.
* **Play back.** After each step, say it back in a sentence, so the user can
  correct it before the next step builds on it.
* **Follow up vague answers.** "It's bad" is not a problem statement. Ask
  what is happening, to whom, and how many. Do not sharpen it yourself.

Ask every question through your question tool, as [asking.md](asking.md)
describes: the open question in the question itself, options only for answers
with evidence and for the choices when a fact is missing, and the answer typed
in as free text. Ask in the chat only if you have no such tool.

## Planning the interview

Before the first question, unless time is short, write a plan in your working
notes:

1. List the sections of the template in use, with the guidance under each.
2. Under each, write the questions this situation raises, using the question
   bank below as prompts, not as a script.
3. Cross off what the fact sheet already answers. Turn what remains into
   specific questions, citing what you found where it helps.

When time is short, skip the written plan, and ask the first unanswered
question of the current step.

## Question bank

For the built-in template, for a live decision. For a decision already made,
ask the same questions in the past tense, about what was known and thought at
the time. Adapt to the template in use, and to the situation.

### Before anything: time

* How long do we have before this must be decided? What happens then?
* Is anything being harmed right now? Is there an obvious action that stops
  it? (If so: act first, record later.)

### Problem

* What is happening? Since when? How do you know?
* Who or what does it affect, and how many?
* How bad is it, and is it getting worse?
* Is this the problem, or a symptom of it? What would solving it look like?

### Information

* Here is what I found, and where. Is any of it wrong, or out of date?
* What else do you know that I could not read?
* What do we not know? Can anyone find out in the time we have?
* Who else should be asked: another team, a supplier, someone who has seen
  this before?
* Is there a runbook, an earlier record or an agreement that applies?

### Options

* What could we do? What else? (Ask again before moving on.)
* Is doing nothing, or waiting, an option? Can we buy time?
* For each option: what could go wrong? What would it gain us? Can it be
  undone?
* Do you want me to suggest options too? (Only if they want; label them.)

### Select

* Which option do you choose?
* Why that one over the others?
* Who is making the decision, and should the record name them, by role?
* Do you want my view first? (Only if they ask; label it.)

### Execute

* What are the steps, in order? Who does each?
* What will show it is working? What will show it is not?
* Is there a point at which you would stop, and try something else?

### Evaluate

* When, and how, will you judge whether it worked?
* If it has not worked by then, what happens next?

When the user returns to evaluate a record:

* What happened? Where can I see it?
* Did it solve the problem? Fully, partly, not at all?
* What did you learn? Is there anything you would do differently?
* Is there anything you see now that you could not see at the time?
* Does the problem remain, or has it changed? If so, is that a new decision?

## Worked example

The user says: "Checkout is erroring after the release. Help us decide what
to do, with PIOSEE."

### Time

"How long do you have, and is anything being harmed right now?" Suppose the
user says customers cannot pay, a sale starts at 16:00, and it is 14:20. Harm
is happening, but the action is not obvious: they do not know whether to roll
back. So the steps go ahead, briskly: one question at a time, drafting as you
go.

### Facts to gather first

Only those that could change the choice, and quickly:

* `pse config` and `pse list`, and any earlier record about checkout or
  releases.
* What changed in the release: `git log` between the tags, and any migration.
* What the error dashboards and logs show, and since when.
* Whether there is a runbook for rolling back, or a feature flag for the new
  code.

Suppose this finds: release 2.3 went out at 14:02; errors rose at 14:05; the
release adds a tax-code feature behind a flag, and a migration that adds a
column; the runbook says to roll back with `deploy rollback`.

### Questions not to ask

* "What is the problem?" (a form field; ask what is happening, and play back)
* "The release caused this, didn't it?" (a guessed cause, as a leading
  question)
* "Should we roll back, fix forward or turn off the flag?" before the user has
  named any options (offering your options in place of theirs)
* "Was turning off the flag chosen because it is lowest risk?" (a guessed
  reason)

### The steps

**Problem.** "What are customers seeing, and how many of them?" Then play
back: "So: since 14:05, about one in twenty-five checkouts fails, customers
cannot pay, and the sale starts at 16:00. Is that the problem?"

**Information.** "Release 2.3 went out at 14:02, with a tax-code feature
behind a flag and a migration adding a column. Errors began three minutes
later. Is anything else known to have changed? Do you know the cause yet?"

**Options.** "What could you do?" Suppose they say: roll back, or fix
forward. "Anything else?" They add: turn off the tax-code flag. "For each:
what could go wrong, and can it be undone?" If they ask whether the migration
matters to a roll back, that is the point to share what you read.

**Select.** "Which do you choose, and why that one?" Record their answer, in
their words.

**Execute.** "Who turns off the flag, and what will show it worked?"

**Evaluate.** "When will you check, and what would count as fixed?" Write it
into Evaluate, with "Not yet evaluated."

### Then draft

Draft with what the user agreed at each step, show the sources, and get a
quick approval. Afterwards, when the user comes back, evaluate in the same
record. If the flag did not help, the next decision, perhaps to roll back, is
a new record that updates this one, or reverses it, as the user judges.
