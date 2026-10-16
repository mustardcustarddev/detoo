# PIOSEE

A summary of what the sources say about PIOSEE, kept here so that the skill
works offline, and so that its advice rests on stated sources. It draws on the
Wikipedia article
[Pilot decision making](https://en.wikipedia.org/wiki/Pilot_decision_making),
which cites PIOSEE to Petr Smejkal, *The Command Handbook* (737 Publishing),
p. 40, and to Flight Deck Friend's article "How do Pilots Make Decisions?".
Read those for the full text. The summary is in our words.

## What it is

PIOSEE is one of several mnemonics that airline crews use to make a decision
in a structured way, often under time pressure and with more than one person
on the flight deck. The article names its steps but does not explain them;
the descriptions below follow its account of the five-step cycle (see below),
which takes the same course in other words.

| Step            | What the crew does                                                                                          |
| --------------- | ----------------------------------------------------------------------------------------------------------- |
| **P**roblem     | Recognise that there is a problem, and state it: what is happening, what it threatens, and how much time there is. |
| **I**nformation | Gather what bears on it, from instruments, documents, each other and people outside, and agree what is known. |
| **O**ptions     | Generate the courses of action open to them, without judging them too early.                                |
| **S**elect      | Weigh the risks and benefits of each option, and choose one.                                                 |
| **E**xecute     | Carry out the chosen option, safely and within the time there is, with each person clear on their part.     |
| **E**valuate    | Check whether the action is working, and whether the situation has changed. If it has not worked, start again from the Problem. |

The method is a cycle, not a line. Evaluation feeds the next round of
decision, and new information at any step can send the crew back to an
earlier one.

## The models around it

The same article describes a family of similar models. They differ mostly in
their words and in what they put first:

* **A five-step cycle**: Situation, Options, Choose, Act, Evaluate. It calls
  recognising the situation, and its dangers, the most important step, and
  asks for as many options as possible, however unlikely at first.
* **T-DODAR** (Time, Diagnose, Options, Decide, Assign, Review), used by
  British Airways, which added the initial T so that the time available is
  considered first. Flight Deck Friend describes PIOSEE as essentially the
  same model with different words.
* **FOR-DEC** (Facts, Options, Risks and benefits, Decide, Execute, Check),
  from Lufthansa and the German Aerospace Center. The hyphen is a pause, to
  check that every option has been considered before deciding.
* **DECIDE**, **DESIDE**, **GRADE**, **3P** and **SHOR**, which follow the
  same arc of perceive, consider, choose, act and check.
* **Aviate, Navigate, Communicate**, which is not a decision model but an
  order of priorities: first keep flying.

The skill borrows two things from these. From T-DODAR, it asks about time
first. From Aviate, Navigate, Communicate, it puts acting to stop harm before
recording anything.

## Why crews use it

The article gives these advantages of such techniques:

* they force the crew to name the facts;
* they stop the crew jumping to conclusions;
* they give a junior crew member a way to be heard;
* they let everyone on the flight deck take part in the decision;
* they let the captain reverse a wrong decision without losing authority.

The sources are about flight decks. Our reading, for a software team, is that
the same holds for an incident channel or a release meeting: the steps make the facts explicit, give everyone a place to raise an
option, and make it normal to revisit a decision when the evaluation says it
did not work.

## Its weaknesses

The article is as clear about these:

* **It can delay a quick, obvious action.** "It is important that the
  technique does not become an obstacle to solving problems." This is why the
  skill asks about time first, and tells the user to act first when harm is
  happening and the action is obvious.
* **It can be used to justify a decision rather than make one.** A record
  written after the event, with the options assembled to make the choice look
  inevitable, is the clearest case. This is why the skill records only the
  options considered at the time, and puts hindsight in Evaluate.
* **It does not capture intuition**, or "gut feeling". When a decision was
  made on instinct, the record says so, rather than inventing an analysis.

The article also says the technique should be standardised across an airline.
For a team, that means one template, used for every PIOSEE record, which is
what the tool's log provides.

## What the sources do not say

The sources describe PIOSEE; none of them report a trial of it. Do not
claim, in a record or to the user, that it is proven to improve decisions, or
that it is better than another model.
