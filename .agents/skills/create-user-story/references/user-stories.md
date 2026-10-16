# User stories

A summary of Mike Cohn's guide,
[User Stories](https://www.mountaingoatsoftware.com/agile/user-stories), at
Mountain Goat Software (last updated 12 July 2026). It is kept here so that the
skill works offline, and so that its advice rests on a stated source. Read the
original for the full text, its examples and its links to further articles.
The summary is in our words; phrases in quotation marks are the guide's.

## What a story is

A user story is "a short description of desired functionality told from the
perspective of someone who wants or needs that functionality". That someone is
usually a user or a customer, but may be a stakeholder, another system, or
anyone else who benefits. A story points at a meaningful outcome, not at an
internal task.

The common template is:

> As a [type of user], I [need/want/am required] to [do something], so that
> [reason or benefit].

It makes the writer think about three things: who the story is for, what they
need, and why it matters. It is a thinking aid, not a ritual: some items read
better without it, some teams use job stories, and some technical work is
clearer as a technical backlog item.

The written sentence is "a reminder to have a conversation". The requirement
itself emerges from that conversation, through examples, rules, sketches,
tests and decisions. A story represents a requirement, but the sentence is
rarely the whole of it.

## Card, conversation, confirmation

Ron Jeffries' three words for a story:

* **Card**: the short written description, wherever it is kept.
* **Conversation**: where understanding grows. The team discusses goals,
  trade-offs, edge cases, business rules, design options and implementation
  concerns.
* **Confirmation**: how the team knows the story is complete: acceptance
  criteria, examples, tests, sketches, rules or notes that say what must be
  true when it is done.

A useful story has all three.

## A good story

* **Has clear value.** Someone cares whether it is done, and the team knows
  why it matters. The value may be to a user, a customer, the business or the
  product; some stories reduce risk or meet a compliance need.
* **Is small enough to finish.** One that cannot be finished in a sprint is
  usually too large. Large stories look like progress, but end the sprint
  "almost done".
* **Is testable.** The team can tell whether it is done. Not every criterion
  has to be known when the story is written, but they must be before it enters
  a sprint.
* **Leaves room for conversation.** It does not prescribe every screen, field,
  table and design choice before the team has discussed the work.

A weak story leaves the user, the goal or the reason vague. The guide's
example: "As a user, I want reports so that I can see information" names none
of the three. "As a regional sales manager, I need to compare monthly revenue
by territory so that I can identify where coaching or support is needed" names
a user, a goal, and the decision the user is trying to make, without
specifying a single field or graph.

Naming the decision the user is trying to make leaves the team free to find
the best solution. "I need a dashboard so that I can see metrics" prescribes a
solution; "I need to see which recently released features are used least so
that I can decide whether to improve, promote, or remove them" does not.

## Acceptance criteria

The details that tell the team whether a story was completed correctly. They
may be examples, rules, tests or brief notes; what matters is the shared
understanding of what must be true for the story to be complete. The guide's
conference example:

> As a conference attendee, I can save sessions to a personal agenda so that I
> can plan which sessions to attend.
>
> * A logged-in attendee can add a session to a personal agenda.
> * A session already on the attendee's agenda is marked as saved.
> * The attendee can remove a saved session.
> * If two saved sessions occur at the same time, the attendee is warned of
>   the conflict.

## How much detail

"Detailed enough for the next decision." A story far down the backlog may be
deliberately vague; one near the top needs enough clarity that the team can
believe it will fit in a sprint. The test the guide offers: "Does the team
understand this story well enough to believe it can be completed in a
sprint?" If not, add detail by refinement, conversation, examples, acceptance
criteria or splitting.

## Epics, splitting, refinement and mapping

* **An epic is a large story.** It keeps an important capability in mind
  before its details are known, and is split into smaller stories as it nears
  development.
* **Splitting** turns a large story into smaller ones that each still
  represent meaningful progress, so the team can finish work and get feedback.
  Do not split only by technical layer.
* **Refinement** adds the right understanding at the right time: vague ideas
  become clearer stories, large ones are split, risks surface, and acceptance
  criteria emerge.
* **Story mapping** lays stories out across a user journey, workflow or
  product area, to reveal missing steps, alternatives and release options.

## Stories, tasks and technical work

A story describes a useful outcome, from the perspective of someone who
benefits. A task describes work the team does to deliver it. "As a customer, I
can save a payment method so that checkout is faster next time" is a story;
"Create payment-token database table" is probably a task.

Technical work can sometimes be described from a meaningful user, business,
risk or operational perspective. Otherwise it is better kept as a technical
backlog item, spike, chore or task. Do not force everything into story form.

## Who writes them

Anyone. The product owner is accountable for the backlog, but developers,
testers, designers, analysts, stakeholders, customers and users all
contribute. Who takes part in the conversation matters more than who typed the
first sentence.

## Common mistakes

* Treating the written story as the whole requirement.
* Writing stories that are too large.
* Splitting stories into tasks.
* Writing every story from the product owner's perspective, when the product
  owner is often not the user.
* Adding too much detail too soon.
* Using the template without thinking.

## Before a story enters a sprint

These should have reasonably clear answers:

* Who benefits from it?
* What outcome, capability, decision or behaviour does it support?
* Why does that outcome matter now?
* Is it small enough to finish in a sprint?
* Has it been split by something other than technical layer?
* Does it leave room for conversation and better design choices?
* Does the team know what must be true for it to be complete?
