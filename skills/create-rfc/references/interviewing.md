# Interviewing for an RFC

A request such as "create an RFC to migrate the system to hexagonal
architecture" names a theme. A draft needs far more: what exactly would
change, why, how far it goes, in what order, at what cost, what else was
weighed, and what the author wants reviewers to settle. Almost all of that is
in the user's head. This is how to get it out, without leading them.

## Principles

* **The template drives the questions.** Read the template in use afresh for
  every RFC. Each section, and the guidance under it, is something a reviewer
  will expect to find; the interview exists to fill those sections.
* **The theme shapes the questions.** "What is the motivation?" is a form
  field. "What is hard to change today that hexagonal architecture would make
  easier?" is a question. Ask what a reviewer of *this* change would need to
  know.
* **Read before you ask.** Learn what the repository can tell you first, so
  that you ask only what you could not read, and can ask precisely: "Nine of
  the twelve modules call the database directly. Are all nine in scope?"
* **Open before closed.** Ask about intent and reasons with open questions.
  Offer choices only where the possible answers are known: which RFC is
  obsoleted, which of the modules you found are in scope, whether to defer a
  point to reviewers.
* **Never offer guessed reasons as choices.** A user picking the nearest
  option from a list you invented is not evidence. "Is it for testability,
  performance or team autonomy?" puts words in their mouth; "What would this
  make possible that is hard now?" does not.
* **Evidence may prompt; it may not answer.** You may put an observation from
  the code to the user, and ask whether it matters: "The handler tests need a
  running PostgreSQL. Is that part of the problem, or not relevant?" Their
  answer is the evidence, not your observation.
* **Small batches.** Three or four related questions at a time, then wait.
  Long questionnaires get short answers.
* **Play back.** After each round, restate in two or three sentences what you
  have understood, so the user can correct it before you build on it.
* **Follow up vague answers.** "It's a mess" is a feeling, not a motivation.
  Ask what happened, where, and how often. Do not sharpen it yourself.
* **Stop when the template is covered.** When every section can be written
  from the fact sheet, with each gap stated or made an unresolved question by
  the user's choice, stop asking and draft. A draft is for comment; it does not
  need every answer.

Use your environment's tool for asking the user structured questions, if it
has one (such as `ask_question` or `AskUserQuestion`), with a free-text answer
always available. Otherwise ask in the conversation, and wait.

## Planning the interview

Before the first question, write a plan in your working notes:

1. List the sections of the template in use, with the guidance under each.
2. Under each section, write the questions this theme raises, using the
   question bank and the theme lenses below as prompts, not as a script.
3. Cross off what the fact sheet already answers. Turn what remains into
   specific questions, citing what you found where it helps.
4. Order the rounds so that each builds on the last: intent and scope first,
   because every other answer depends on them.

Revise the plan as answers come in. An answer often makes later questions
unnecessary, or raises new ones.

## Question bank

For the built-in template. Adapt to the template in use, and to the theme.

### Round 1: intent and scope (Title, Summary)

* In a sentence or two, what are you proposing?
* What do the theme's terms mean in this proposal? ("Hexagonal", "event
  sourcing", "monorepo" mean different things to different people.)
* What does the end state look like, concretely?
* What is in scope, and what is deliberately out?
* Who needs to comment on it: which teams or people does it affect?

### Round 2: motivation

* What is the problem today? Can you point to an example: a change, an
  incident, a request that was hard?
* Who feels the problem, and how?
* Why now, rather than earlier or later?
* What happens if nothing is done?
* How would you know, afterwards, that the change worked? (Only if the user
  has a measure; do not invent one.)

### Round 3: design

* Which parts of the system change, and how? Confirm against what you found.
* What are the new boundaries, interfaces, components or rules?
* Is the change done at once, or incrementally? In what order? What coexists
  with what during the transition, and for how long?
* What stays the same, that a reviewer might assume changes?
* Are there compatibility, migration or rollback concerns: data, APIs,
  deployments, users?
* What, if anything, is undecided? Should reviewers settle it, or will you?

### Round 4: drawbacks

* What does it cost: effort, time, disruption, new things to learn or run?
* What could go wrong, and what would you do about it?
* What objections do you expect from reviewers?
* Who loses something by it?

### Round 5: alternatives

* What else did you consider? Why not that?
* Did you weigh doing nothing, or doing less? What would that cost?
* Is there a smaller change that would solve most of the problem?

### Round 6: unresolved questions, and the record

* What do you most want reviewers' views on?
* Which of those must be settled before the RFC can be accepted, and which can
  wait until the work is under way?
* Does this obsolete or update an existing RFC? Should the link be made now,
  or when this one is accepted? (See the warning in
  [rfc-tool.md](rfc-tool.md#links).)

## Theme lenses

Different kinds of change raise different questions. Use these as prompts
when planning; a proposal may fit more than one.

* **Restructuring the architecture** (layers, hexagonal, modular monolith,
  splitting services): what the boundaries are and who owns each; which way
  dependencies point; what moves where; whether it is incremental; how old and
  new structures coexist; how the rule is kept once the migration ends.
* **Adopting a technology** (a database, a framework, a protocol, a platform):
  what it replaces; who will operate it and whether they know how; cost and
  licence; how data or code moves to it; what leaving it would take.
* **Changing an interface** (a public API, a file format, a CLI, a schema):
  who consumes it; whether it breaks them; versioning and deprecation; how long
  the old form is kept; how consumers are told.
* **Migrating data**: volume; downtime; backfill; verification; rollback; what
  happens to writes during the move.
* **Changing a process** (review, release, on-call, the RFC process itself):
  who must do something differently; what is enforced and how; the transition;
  how to tell whether it is working.
* **A cross-cutting concern** (security, observability, error handling,
  configuration): where it applies, and where it does not; how consistency is
  kept; what existing code must change; what a new component must do.

## Worked example

The user says: "Create an RFC to migrate the system to hexagonal
architecture."

### Facts to gather first

Read in proportion to the change; here, the structure of the whole system.

* `rfc config`, `rfc list`, RFC 1 for the project's process, and any RFC or
  ADR on architecture or structure.
* The top-level layout, and the main modules or packages.
* Where the system does input and output: HTTP handlers, database access,
  message queues, files, third-party clients. Which modules contain it.
* Which way dependencies point: does business logic import database or HTTP
  packages? Count, with a command.
* How the tests reach external systems: real services, containers, fakes.
* Any existing interfaces that already separate logic from infrastructure.

Suppose this finds: twelve modules; nine call the database directly from the
HTTP handlers; the handler tests need a running PostgreSQL; the payment
provider's client is called from four modules; no RFC or ADR concerns
structure.

### Questions not to ask

* "What is the motivation for this RFC?" (a form field, not a question)
* "Is the motivation testability, maintainability or flexibility?" (guessed
  reasons as choices)
* "Shall I describe the standard benefits of hexagonal architecture?" (an
  invitation to invent)
* "Which modules are there?" (you can read that)

### The plan, round by round

**Round 1: intent and scope.**

1. "What does hexagonal architecture mean for this system, in your words: what
   would be a port, and what an adapter?"
2. "I found twelve modules, of which nine call the database directly from the
   handlers: [list]. Is the whole system in scope, or some of those?"
3. "Who would need to review this?"

Then play back what the user said. If, say, they had answered that the
database and the payment provider are the adapters, and that orders and
payments come first while reporting is out:

> "So you're proposing to put the database and the payment provider behind
> interfaces owned by the business logic, starting with orders and payments,
> and leaving reporting out. Is that right?"

**Round 2: motivation.**

1. "What is hard today that this would make easier? Is there a recent example?"
2. "Why now?"
3. "The handler tests need a running PostgreSQL. Is that part of the problem,
   or not relevant to you?"

**Round 3: design.**

1. "Would you migrate one module at a time, or all at once? If one at a time,
   which first, and why?"
2. "While the migration is under way, will migrated and unmigrated modules
   coexist? Is there anything they share that makes that hard?"
3. "Once it's done, how should the boundary be kept: a rule in review, a lint,
   an architecture test?"
4. "Is any of this undecided, and better left to reviewers?"

**Round 4: drawbacks.** "What will it cost, and what objections do you
expect?" Then, if the user wants a prompt: "Every module in scope changes,
which is nine of the twelve. Should the RFC say what that means for other work
during the migration?"

**Round 5: alternatives.** "Did you consider anything else? Why not that?",
then "Did you weigh doing nothing, or doing less?". Do not list alternatives
the user has not mentioned, such as "wrapping only the payment provider", as
if they had weighed them. If they ask for ideas, offer some as your
suggestions, and include only those they say they actually weighed, or now
want reviewers to weigh.

**Round 6: unresolved questions and the record.** "What do you most want
reviewers to settle? Must any of it be settled before acceptance?"

### Then draft

Draft only when the fact sheet covers every section. The proposal, its scope,
its reasons and its sequence are the user's; the counts and the current
structure are yours, with their sources; anything still open is an unresolved
question, because the user chose to make it one.
