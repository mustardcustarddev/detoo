# Writing an RFC

An RFC has two audiences. Its reviewers read it now, and comment on it. Its
later readers find it in the log, and want to know what was proposed and why.
Both can only use it if it says what the author actually thinks. A plausible
reason that nobody holds draws comments on a position nobody holds, and, if
the RFC is accepted, is believed afterwards.

So the author's intent and reasoning are the substance of the RFC, and they
are the one thing you cannot read from the code. You gather them by asking.

## Facts first

Before writing any prose, build a fact sheet: every claim the RFC will make,
and where it came from. Keep it in your working notes. It is shown to the user
with the draft, and is not written into the RFC.

| #   | Claim                                                                       | Source                                                            |
| --- | --------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| 1   | The HTTP handlers call `database/sql` directly, in 23 places across 9 files | `grep -rn "sql\." internal/http`, run today                       |
| 2   | Tests of the handlers need a running PostgreSQL                             | `internal/http/main_test.go`, line 14                             |
| 3   | The problem: changing the payment provider took three weeks                 | User, in this conversation                                        |
| 4   | The migration is to be incremental, one module at a time                    | User, in this conversation                                        |
| 5   | Which module goes first                                                     | **Unknown — ask**                                                 |
| 6   | Whether the reporting module is in scope                                    | User: undecided, wants reviewers' views — **Unresolved question** |

### Acceptable sources

* **The user**, in this conversation: what they stated, or explicitly agreed
  to. A suggestion of yours that they did not take up is not a fact.
* **The repository**: code, configuration, documents, RFCs and ADRs that you
  have read. Note the path, and the line where it matters.
* **Commands you ran**, and their output: counts, dependency lists, test
  results, `git log` history.
* **A document the user pointed you to**, such as an issue, an incident report
  or a design note.

### Not acceptable

Do not write any of these unless a source above supports it:

* a motivation, a problem, or a pressure the team is under;
* a goal, a benefit, or a success measure;
* a scope: what is in or out of the proposal;
* a design choice: a boundary, a sequence, a technology, a rule;
* an alternative that was considered, or why it was rejected;
* a drawback, a cost, an estimate, a risk;
* a stakeholder, a team, an incident, a deadline or a requirement;
* a prediction: that something "will improve", "scales better", "reduces
  risk";
* a general claim such as "industry standard" or "best practice".

The code shows how things are now. It does not show what the user wants
instead, or why. An inferred reason is a guess, however likely it seems. Ask.

### Your own suggestions

The user may ask for your help with the proposal itself: "what would the ports
be?", "what are the risks?", "what else could we do?". That is a fair request.
Answer it in the conversation, labelled as your suggestion, grounded in what
you read in the code, and with any claim you cannot evidence marked as your
opinion.

A suggestion enters the RFC only when the user adopts it; it is then their
proposal, and the fact sheet says so ("User adopted the agent's suggestion, in
this conversation"). Adopting a design does not adopt a reason for it: the
reasons still come from the user, or from evidence.

Do not volunteer suggestions in place of asking. Ask the open question first;
offer help only if the user asks for it, or says they do not know and would
like ideas.

### When a fact is missing

Ask the user. If they do not know, or prefer not to say, do one of these, and
let them choose:

* **Make it an unresolved question.** An RFC is a request for comments, so an
  open point the author wants reviewers' views on belongs in Unresolved
  questions, stated as a question. This is the honest home for most gaps in
  the Design, the Drawbacks and the Alternatives.
* **Say so plainly** in the section: "Alternatives were not explored before
  this RFC; reviewers' suggestions are welcome."
* **Leave it out**, where the template does not need it.

Never fill the gap with something plausible.

Two gaps cannot be deferred. If the user cannot say **what** is proposed, or
**why**, the RFC is not ready to be drafted: tell them so, and suggest what
would settle it, rather than writing a proposal without a purpose.

## The sections

The built-in template has a title, a date, a status, and six sections. A
project's own template may differ: fill in each of its sections by the
guidance written in it, with the same standard of evidence.

### Title

A short phrase naming the proposed change, such as "Migrate the order service
to hexagonal architecture" or "Adopt gRPC for internal services". Run
`rfc list` and follow the style of the existing titles. The title becomes the
file name, so keep it short.

### Date

The tool writes today's date: the day the draft is created. Set `RFC_DATE`
only if the user asks for another.

### Status

Always `Draft`. The tool writes it; leave it. Accepting, rejecting or
withdrawing an RFC is the outcome of the comment process, and is recorded by
people, not by you, even if the user says the proposal is already agreed. If
it is already agreed and done, it is a decision to record, not a proposal to
comment on: suggest an ADR instead.

### Summary

The proposal in one paragraph, readable on its own: what would change, and,
in a clause, why. A reviewer who reads only this should know whether the RFC
concerns them. Write it last, from the other sections.

### Motivation

Why the change should be made: the problem it solves, and who it matters to.

* State the problem as the user described it, with the evidence they gave or
  you measured. Concrete beats general: "changing the payment provider took
  three weeks and touched 40 files" beats "the code is hard to change".
* Say who is affected: which teams, users or systems, and how.
* Say why now, if the user said.
* Do not inflate. If the user's reason is modest, the motivation is modest.

### Design

The change, in enough detail to comment on and to build.

* Describe the end state: what exists afterwards that does not now, and what
  goes away. Name concrete things: modules, directories, interfaces, services,
  commands, rules.
* Describe the path, if the user has one: the order of the work, how old and
  new coexist during it, and what is done when.
* Say what is out of scope, as the user drew it.
* Where the user has not decided a point, do not decide it for them: list it
  under Unresolved questions, and refer to it from here if that helps.
* Use the current state from the fact sheet to make the change concrete: "the
  23 direct calls to `database/sql` in `internal/http` would move behind a
  `Store` interface".

### Drawbacks

Why the change might not be made: its costs, its risks, and who loses
something by it. This is the case against, and it should be the honest one.

* Ask the user for the drawbacks they see, and the objections they expect.
* A cost that follows mechanically, and can be checked, may be stated from the
  evidence: "every module in scope would change; that is 9 of the 12".
* A judgement or a prediction needs the user's agreement: "the migration will
  slow feature work for a quarter".
* An RFC with no drawbacks is not credible. If the user can name none, say so
  plainly rather than inventing one, and consider asking reviewers in
  Unresolved questions.

### Alternatives

What else was considered, and why this design was chosen over it.

* List only what the user actually considered, with the reasons they gave.
* Doing nothing is an alternative; ask whether it was weighed, and what it
  would cost.
* Do not add alternatives nobody considered, to make the proposal look
  thorough. If the user wants ideas, see "Your own suggestions" above.

### Unresolved questions

What should be settled before the RFC is accepted, and what can wait. Keep the
two apart.

* These are the questions the author wants reviewers to answer: the main
  reason the RFC exists.
* Each is a real question, specific enough to answer: "Should the reporting
  module be in scope, given it is due to be replaced next year?" not "Scope?".
* If there are none, write "None."

## Tone

An RFC argues for a proposal, which an ADR's context does not. It may
persuade, but only with the author's actual reasons, and it must give the
case against honestly. Write as the author proposing: "We propose …", "This
RFC proposes …", or the style of the existing RFCs in the log.

## Style

* Full sentences, organised in paragraphs. Bullets suit lists, such as steps,
  alternatives or questions, but are not an excuse for fragments.
* As long as the proposal needs, and no longer. A reviewer should be able to
  read it in one sitting.
* Plain words. Prefer "we will move X behind an interface because Y" to "X will
  be leveraged to enable Y".
* No marketing: "robust", "seamless", "modern", "clean", "powerful",
  "significantly", "simply", "obviously" add nothing a reviewer can check.
* Keep the user's meaning. Do not turn "it seemed easier" into "it reduces
  complexity", or "we might need it" into "it is required".
* Define the theme's terms as the user defined them. "Hexagonal
  architecture", "microservices" or "event-driven" mean different things to
  different readers; say what they mean in this proposal.
* No placeholders and no template guidance left behind.

## Examples

Each pair shows what to write instead, if the evidence supports it.

* Instead of "Hexagonal architecture is the industry best practice for
  maintainable systems", write "Replacing the payment provider took three
  weeks, because provider calls are spread through the handlers."
* Instead of "This will greatly improve testability", write "The handlers'
  tests need a running PostgreSQL today; with the database behind an
  interface, they could use an in-memory fake", but only if the user gave this
  as a reason.
* Instead of "We evaluated several options", write "We considered one
  alternative: wrapping only the payment provider, and leaving the rest as it
  is."
* Instead of "The migration carries minimal risk", ask what risks the user
  sees, and write those.
* Instead of "Phase 2 will take about two sprints", ask whether there is an
  estimate. If not, leave it out, or ask reviewers.

## Before you finish

* Every sentence traces to a line of the fact sheet.
* Every "Unknown — ask" line was asked about, and the answer is used, made an
  unresolved question, stated as a gap, or left out, by the user's choice.
* What is proposed, and why, both come from the user.
* Every section of the template in use is filled in, and no guidance text is
  left.
* The status is `Draft`.
* Links between RFCs were made by the tool, not by hand, and only when the user
  chose to make them now.
* The user has seen the draft with its sources, and approved it.
