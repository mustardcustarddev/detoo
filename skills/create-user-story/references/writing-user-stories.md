# Writing a user story

A story is read by the people who will refine it, build it and test it, and
later by whoever wants to know why the software does what it does. It is "a
reminder to have a conversation", so it is short; but each thing it does say
is taken as true. A user who does not exist, a reason nobody gave, or a
criterion nobody agreed to sends the conversation the wrong way, and gets
built.

So who the story is for, what they need, and why, are the user's. You gather
them by asking. See [user-stories.md](user-stories.md) for what makes a story
good.

## Facts first

Before writing any prose, build a fact sheet: every claim the story will make,
and where it came from. Keep it in your working notes. It is shown to the user
with the draft, and is not written into the story.

| #   | Claim                                                            | Source                              |
| --- | ---------------------------------------------------------------- | ----------------------------------- |
| 1   | Invoices can be viewed one at a time, and not exported           | `internal/web/invoices.go`, read    |
| 2   | The `finance` role can see every customer's invoices             | `internal/auth/roles.go`, line 31   |
| 3   | The story is for finance clerks reconciling the month's invoices | User, in this conversation          |
| 4   | Reason: reconciliation means retyping each invoice into a sheet  | User, in this conversation          |
| 5   | Cancelled invoices are left out of the export                    | User, in this conversation          |
| 6   | Whether credit notes belong in the export                        | User: undecided — **Open question** |
| 7   | Which columns the export has                                     | **Unknown — ask**                   |

### Acceptable sources

* **The user**, in this conversation: what they stated, or explicitly agreed
  to. A suggestion of yours that they did not take up is not a fact.
* **The repository**: code, configuration, documents, stories, RFCs and ADRs
  that you have read. Note the path, and the line where it matters.
* **Commands you ran**, and their output: counts, routes, roles, test results.
* **A document the user pointed you to**, such as an issue, a support ticket,
  a design or a note of a conversation with users.

### Not acceptable

Do not write any of these unless a source above supports it:

* a type of user, a persona, or a role the software has;
* a goal, a benefit, a reason, or a decision the user is trying to make;
* an acceptance criterion: a rule, an example, a limit, an edge case;
* a scope: what is in or out of the story;
* a design choice: a screen, a field, a format, a technology;
* a size, an estimate, a priority or a deadline;
* a prediction: that something "will save time", "improves conversion".

The code shows what the software does now. It does not show who wants it to
do more, or why. An inferred reason is a guess, however likely it seems. Ask.

### Your own suggestions

The user may ask for your help: "what would the acceptance criteria be?",
"how would you split it?", "what edge cases are there?". That is a fair
request, and often the most useful thing you can do. Answer it in the
conversation, labelled as your suggestion, grounded in what you read in the
code.

A suggestion enters the story only when the user adopts it; it is then theirs,
and the fact sheet says so ("User adopted the agent's suggestion, in this
conversation"). Adopting a criterion does not adopt a reason for the story:
the reason still comes from the user.

Do not volunteer suggestions in place of asking. Ask the open question first;
offer ideas only if the user asks for them, or says they do not know and would
like some.

### When a fact is missing

Ask the user. If they do not know, or prefer not to say, do one of these, and
let them choose:

* **Make it an open question** in the Notes. A story is a placeholder for a
  conversation, and much of its detail is meant to be settled later, in
  refinement. An honest open question is the right home for most gaps in the
  acceptance criteria.
* **Leave it out**, where the template does not need it.

Never fill the gap with something plausible.

Three gaps cannot be deferred. If the user cannot say **who** the story is
for, **what** they need, or **why**, the story is not ready to be written:
tell them so, and suggest what would settle it, such as asking the people who
will use it.

## The sections

The built-in template has a title, a date, a status, and three sections. A
project's own template may differ: fill in each of its sections by the
guidance written in it, with the same standard of evidence.

### Title

A short phrase naming what the user can do, from their side: "Export invoices
for a month", "Reorder a previous purchase", "Warn of clashing sessions". Not
a task ("Add CSV endpoint") or a component ("Invoice exporter"). Run
`sty list` and follow the style of the existing titles. The title becomes the
file name, so keep it short.

### Date

The tool writes today's date: the day the story is written. Set `STY_DATE`
only if the user asks for another.

### Status

Always `Draft`. The tool writes it; leave it. Whether a story is ready for a
sprint, in progress, done or dropped is for the team to decide and record,
even if the user says it is ready.

### Story

One sentence in the template's form, unless the user prefers another (see
[user-stories.md](user-stories.md#what-a-story-is)):

> As a [type of user], I [need/want/am required] to [do something], so that
> [reason or benefit].

* **The type of user** is a real role, as the user names it: "finance clerk",
  "returning customer", "on-call engineer". Not "user", and not the product
  owner unless they are the one who benefits. Use the software's own name for
  the role if it has one and the user agrees it is the same people.
* **Need, want or am required to** is the user's choice: "need" and "want"
  carry different weight to whoever prioritises, and "am required to" signals
  a rule from outside, such as a regulation.
* **The something** is what they can do, not how the software does it. "Export
  the month's invoices", not "click a CSV button on the invoices page", unless
  the user has decided the design and wants it fixed.
* **The reason or benefit** is the user's own, and is often the decision the
  user of the software is trying to make. If it only restates the goal ("so
  that I can export invoices"), ask what they do with the result.

If the user wants several sentences, or a story from another system's point of
view ("As the billing service, …"), write it their way.

### Acceptance criteria

What must be true for the story to be complete: the confirmation.

* A bulleted list. Each item is one checkable statement, in the present tense,
  of what the user of the software can do or will see: "Cancelled invoices are
  not included."
* Examples are welcome where a rule is easier to see in one: "An invoice for
  £1,200.50 appears as `1200.50`."
* Include the unhappy paths the user cares about: empty results, errors,
  permissions, limits.
* Criteria describe behaviour, not implementation. "The export opens in
  Excel" is a criterion; "uses the encoding/csv package" is not.
* Only criteria the user stated or agreed to. Your suggestions are offered in
  the conversation, and enter only if adopted.
* Not every criterion has to be known now. Write what is agreed, and put the
  rest in the Notes as open questions. If none are agreed yet, say so plainly:
  "None agreed yet; see the open questions."

### Notes

What the conversation has settled so far, and what it has still to settle.
Keep it short; a story is not a specification.

* **Context** the team needs that the sentence does not carry, such as who
  asked for it, or the example that prompted it.
* **Out of scope**, as the user drew it.
* **Constraints and dependencies** the user named, or that you found and they
  confirmed: "Needs the `finance` role, which exists today."
* **Open questions**, each specific enough to answer in refinement: "Should
  credit notes be in the export, or a separate one?" not "Credit notes?".
* **Links** to related stories, RFCs, ADRs, issues or designs, by name. Links
  between stories are made by the tool; see [sty-tool.md](sty-tool.md#links).

If there is nothing to say under a heading, leave the heading out. If there is
nothing at all, write "None."

## Size

A story should be small enough for the team to finish in a sprint. You cannot
judge that from the code alone, and you should not guess. Ask the user. If they
think it is too large, or the acceptance criteria run long, or the story joins
several goals with "and", it is probably an epic:

* Suggest splitting it into stories that each deliver something the user of
  the software can use: by step in a workflow, by business rule, by kind of
  data, by the simple case first and the variations after, or by user. Not by
  technical layer: "the API" and "the UI" are tasks, not stories.
* Offer your split as a suggestion, and let the user choose.
* Ask whether to record the large story too. If they want it, create it first,
  then each smaller story with `-r` on the large one; see
  [sty-tool.md](sty-tool.md#links).

## Style

* Plain words, in the user's terms and the domain's, not the code's. "Invoice",
  not "InvoiceDTO".
* Short. The story, a handful of criteria, and brief notes. Detail belongs in
  the conversation until it is needed.
* No marketing: "seamless", "intuitive", "robust", "easily", "quickly" add
  nothing anyone can test. If speed matters, ask how fast, and make it a
  criterion.
* Keep the user's meaning. Do not turn "it's a pain" into "it reduces errors",
  or "they might want" into "they need".
* No placeholders and no template guidance left behind.

## Examples

Each pair shows what to write instead, if the evidence supports it.

* Instead of "As a user, I want an export button so that I can export data",
  write "As a finance clerk, I need to export the month's invoices so that I
  can reconcile them against the bank statement without retyping them."
* Instead of "As the product owner, I want reorder functionality", ask who
  would reorder, and why, and write theirs: "As a returning customer, I want
  to reorder items from a previous purchase so that I can quickly buy things I
  regularly need."
* Instead of the criterion "The export is fast", ask whether there is a limit
  that matters. If there is, write it: "A month's invoices export in under ten
  seconds." If not, leave it out.
* Instead of the criterion "Create an `invoice_exports` table", leave it out:
  it is a task, not something the user of the software can check.
* Instead of inventing "Credit notes are excluded", ask. If the user is
  undecided, write the open question.

## Before you finish

* Every sentence traces to a line of the fact sheet.
* Who, what and why all come from the user.
* Every "Unknown — ask" line was asked about, and the answer is used, made an
  open question, or left out, by the user's choice.
* Every acceptance criterion is checkable, and was stated or agreed by the
  user.
* The story describes an outcome, not a task, and does not fix a design the
  user has not decided.
* Every section of the template in use is filled in, and no guidance text is
  left.
* The status is `Draft`.
* Links between stories were made by the tool, not by hand.
* The number in the title is the number in the file name.
* The user has seen the draft with its sources, and approved it.
