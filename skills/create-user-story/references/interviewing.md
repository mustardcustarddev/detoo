# Interviewing for a user story

A request such as "write a user story for exporting invoices" names a
subject. A story needs more: who would use it, what exactly they need to do,
why, what must be true when it is done, and whether it is small enough to
finish. The code can say what exists now. The rest is in the user's head, or
in the heads of the people they speak for. This is how to get it out, without
leading them.

## Principles

* **The template drives the questions.** Read the template in use afresh for
  every story. Each section, and the guidance under it, is something the team
  will expect to find; the interview exists to fill those sections.
* **The subject shapes the questions.** "Who is the user?" is a form field.
  "Who exports invoices today, and what do they do with them afterwards?" is a
  question. Ask what the people building and testing *this* story would need
  to know.
* **Read before you ask.** Learn what the repository can tell you first, so
  that you ask only what you could not read, and can ask precisely: "The
  software has three roles that can see invoices: `finance`, `admin` and
  `account-owner`. Which of them is this for?"
* **Who, what and why come first.** Every other answer depends on them, and
  none of them can be read from the code.
* **Outcomes, not solutions.** Ask what the user of the software is trying to
  get done, and what decision they are trying to make. If the user describes
  a screen or a button, ask what it is for; keep the design if they have
  decided it, but write the story around the outcome.
* **Open before closed.** Ask about users and reasons with open questions.
  Offer choices only where the possible answers are known: which of the roles
  you found, which existing story this replaces, whether to leave a point to
  refinement.
* **Never offer guessed reasons as choices.** A user picking the nearest
  option from a list you invented is not evidence. "Is it to save time, reduce
  errors or meet an audit?" puts words in their mouth; "What do they do with
  the invoices once they have them?" does not.
* **Evidence may prompt; it may not answer.** You may put an observation from
  the code to the user, and ask whether it matters: "Cancelled invoices are
  kept, with a `cancelled` flag. Should the export include them?" Their answer
  is the evidence, not your observation.
* **Ask for examples.** Concrete cases make the best acceptance criteria:
  "Walk me through the last time someone did this. What did they need?"
* **Small batches.** Three or four related questions at a time, then wait.
* **Play back.** After each round, say the story back in its one sentence,
  with the criteria so far, so the user can correct it before you build on it.
* **Follow up vague answers.** "It's clunky" is a feeling, not a reason. Ask
  what happens, to whom, and how often. Do not sharpen it yourself.
* **Stop at enough.** A story needs detail enough for the next decision, not
  every answer. When every section can be written from the fact sheet, with
  each gap left as an open question by the user's choice, stop asking and
  draft.

Use your environment's tool for asking the user structured questions, if it
has one (such as `ask_question` or `AskUserQuestion`), with a free-text answer
always available. Otherwise ask in the conversation, and wait.

## Planning the interview

Before the first question, write a plan in your working notes:

1. List the sections of the template in use, with the guidance under each.
2. Under each section, write the questions this subject raises, using the
   question bank and the subject lenses below as prompts, not as a script.
3. Cross off what the fact sheet already answers. Turn what remains into
   specific questions, citing what you found where it helps.
4. Order the rounds so that each builds on the last: who, what and why first.

Revise the plan as answers come in.

## Question bank

For the built-in template. Adapt to the template in use, and to the subject.

### Round 1: who and what (Title, Story)

* Who would use this? What do you call them? Is there more than one kind of
  person who would?
* What do they need to be able to do? What do they do instead today?
* Is the person asking for it the person who would use it? If not, who
  speaks for the users?
* Is this one story, or several? (See round 4.)

### Round 2: why (Story)

* What do they do with the result? What decision does it help them make?
* What happens today without it? Is there an example: a ticket, a complaint,
  a time it went wrong?
* Do they need it, want it, or are they required to have it? If required, by
  what?
* Why now?

### Round 3: confirmation (Acceptance criteria)

* How would you, or a tester, know it is done? What would you try first?
* Walk me through an example, with real values.
* What should happen when there is nothing to show, when something fails, or
  when the person is not allowed to?
* Are there limits that matter: sizes, times, formats, numbers of things?
* What should it deliberately not do?
* Which of these must be settled now, and which can wait for refinement?

### Round 4: size

* Could the team finish it in a sprint, do you think?
* If not, what is the smallest part that someone could use on its own? What
  could come after?
* Should the large story be recorded too, and replaced by the smaller ones?

### Round 5: notes and links

* Is anything out of scope that someone might assume is in?
* Does it depend on anything: another story, a decision, a system, a person?
* What is still open, that the team should settle in refinement?
* Does it replace or change an existing story? (See
  [sty-tool.md](sty-tool.md#links).)

## Subject lenses

Different kinds of story raise different questions. Use these as prompts when
planning; a story may fit more than one.

* **A new capability for an end user**: who they are and how they arrive at
  it; what they do instead today; the happy path; what they see when it fails;
  who is not allowed to.
* **A change to existing behaviour**: what it does now, which you can read;
  who relies on that; what changes for them; whether old data or old links
  must keep working; which earlier story it updates.
* **An internal or back-office user** (support, finance, operations): what
  part of their job it serves; what they use instead now, such as a
  spreadsheet or a request to engineering; how often; the volumes involved.
* **A rule from outside** (regulation, contract, policy): which rule, and
  where it is written; who is required to do what; what evidence an auditor
  would want; the deadline. "Am required to" usually belongs here.
* **Another system as the user** (an integration or an API consumer): which
  system, and who owns it; what it needs and in what form; how it handles
  failure; whether the contract is versioned.
* **A quality** (speed, accessibility, security, reliability): who notices it
  and when; the measure that would satisfy them. It may be better as
  acceptance criteria on the stories it affects than as a story of its own;
  ask.
* **Technical work** (an upgrade, a refactoring, a spike): whether someone who
  benefits can be named, such as an operator or a future developer, and what
  they get. If not, say so, and suggest it be kept as a technical backlog
  item rather than forced into story form. It is the user's call.
* **A large or vague subject** ("a reporting area"): it is an epic. Ask what
  the first useful slice is, and offer to record the epic and its first
  stories.

## Worked example

The user says: "Write a user story for exporting invoices."

### Facts to gather first

Read in proportion to the subject; here, the invoice feature and who can use
it.

* `sty config`, `sty list`, story 1 for the project's conventions, and any
  story, RFC or ADR about invoices or exports.
* Where invoices are shown and stored, and what fields they have.
* Which roles can see invoices, and any existing exports elsewhere in the
  software that the team might expect this to match.

Suppose this finds: invoices are shown one at a time on an invoice page;
three roles can see them; there is a CSV export of customers already, in
`internal/web/customers.go`; cancelled invoices are kept with a flag; no
story concerns invoices.

### Questions not to ask

* "Who is the user?" (a form field, not a question)
* "Is the export to save time, reduce errors or support audits?" (guessed
  reasons as choices)
* "Shall I write the standard acceptance criteria for a CSV export?" (an
  invitation to invent)
* "Where are invoices shown?" (you can read that)

### The plan, round by round

**Round 1: who and what.**

1. "Three roles can see invoices: `finance`, `admin` and `account-owner`.
   Who would export them?"
2. "What do they do today when they need many invoices at once?"
3. "Export which invoices: all of them, a month, a customer's?"

Then play back. If the user had answered that finance clerks need a month's
invoices, and copy them out by hand today:

> "So: as a finance clerk, I need to export a month's invoices, so that … —
> what do they do with them once they have them?"

**Round 2: why.**

1. "What do they do with the month's invoices?"
2. "Is there a recent example of it going wrong, or taking too long?"

**Round 3: confirmation.**

1. "Walk me through an export of last month: what would the clerk choose,
   and what would they get?"
2. "There is a customer export already, in CSV. Should this one match it, or
   is the format still open?"
3. "Cancelled invoices are kept, with a flag. Should they be in the export?"
4. "What should happen in a month with no invoices?"

**Round 4: size.** "Could the team finish this in a sprint?" If the user
wants a choice of formats and a schedule too, offer a split as your
suggestion, such as one month in one format first, then other formats, then a
scheduled export, and let them choose.

**Round 5: notes and links.** "Is anything out of scope? What is still open
for refinement?"

### Then draft

Draft only when the fact sheet covers every section. The user of the
software, their need and their reason are the user's; the roles, the existing
export and the cancelled flag are yours, with their sources; anything still
open is an open question in the Notes, because the user chose to leave it
there.
