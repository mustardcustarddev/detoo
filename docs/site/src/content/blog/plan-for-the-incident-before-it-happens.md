---
title: 'Plan for the incident before it happens: PIOSEE for software engineers'
description: Pilots decide under pressure with PIOSEE. Software teams can do better still, by working through it before the incident, and writing the plan down.
pubDate: 2026-10-17
tags: [piosee, incidents, decision-making]
tool: pse
draft: false
---

A customer emails support at ten to six on a Friday. They have been charged
twice for their subscription, and they would like their money back, please.

Somebody has to decide what happens next. Is it really two charges, or one
charge and a pending authorisation that the bank will release on its own? Is
it one customer, or the first of a hundred? Who is allowed to issue a refund,
and how? Does anyone need to be told?

None of these is a hard question. But they will be answered by whoever picks
up the ticket, on whatever facts they can find before they go home, and
answered differently by the next person who picks up the next one. The real
cost of an unplanned incident is that run of hurried, inconsistent decisions,
each of them defensible on its own.

## What pilots do

Pilots have a name for a way of making a decision when something goes wrong:
PIOSEE. It stands for six steps, taken in order.

1. **Problem.** What exactly has gone wrong, and what has to be decided?
2. **Information.** What do we know, and what do we need to find out?
3. **Options.** What could we do?
4. **Select.** Which option, and by what rule?
5. **Execute.** Do it.
6. **Evaluate.** Did it work? If not, go round again.

There is nothing clever in that list. Its value is that it stops you jumping
from the problem straight to the first action that comes to mind. It makes
you look at the facts, consider more than one course, and check the result.

A pilot works through it in the cockpit, in the moment. Software teams have
an advantage pilots lack in the moment: most of our incidents are not
surprises. We know customers will be charged twice, because it has happened
before. We know the certificate will expire, the queue will back up, and the
third-party API will go down on a bank holiday. We can do the thinking now,
while nobody is waiting.

## Playbooks and runbooks

You may already keep runbooks. A runbook says how to do something: how to
restart the worker, how to rotate the key, how to issue a refund. It assumes
you have already decided to do it.

A PIOSEE playbook comes before that. It says how to decide *what* to do. It
names the situation it is for, the decision the situation calls for, the
facts to gather first, the options, and the rules for choosing between them.
Only then does it get to the steps, which may well point to your runbooks.

The rules for choosing are the part that matters most, and the part most
easily left unwritten. "Refund if appropriate" leaves the decision to whoever
is on duty. A rule reads more like "Refund the second charge if both charges
have settled and the customer has one subscription; otherwise escalate to the
billing owner." Somebody new
can follow it at ten to six on a Friday and make the same call your most
experienced engineer would.

## The double charge, planned in advance

Here is what a playbook for that Friday email might hold. Each heading is a
section of the playbook.

**Scenario.** A customer reports being charged more than once for the same
thing.

**Problem.** Should we refund, and how much, and is this one customer or a
fault affecting many?

**Information.** Gather these before deciding:

* The charges in the payment provider's dashboard: whether each one has
  settled, or is still a pending authorisation.
* Whether the customer has one subscription or two, perhaps from signing up
  twice.
* Our own logs for the payment: whether the request was retried, and whether
  it carried the same idempotency key each time.
* How many other customers were charged twice in the same window.

**Options.**

1. Explain that one charge is a pending authorisation, which the bank will
   release, and refund nothing.
2. Refund the duplicate charge.
3. Treat it as an incident: stop the job that takes payments, and find the
   fault before refunding anyone.

**Selection.**

* If one of the charges is still pending, choose option 1.
* If both have settled, and only this customer is affected, choose option 2.
* If more customers than the threshold your team sets are affected, choose
  option 3, whatever else is true.

**Execution.** For option 2: issue the refund in the provider's dashboard,
reply to the customer with the time it takes to reach their account, and
record what was done where the team records its incidents. For option 3:
page the on-call engineer, and follow the incident process from there.

**Evaluation.** Did the customer have their money back within the time we
promise? Did we find out why it happened? Would the playbook have led
someone new to the same decision?

That is perhaps a page. It took an hour or two to write, mostly spent finding
out where the evidence lives and who is allowed to do what. Every one of
those answers would otherwise have been worked out under pressure.

## When not to bother

A playbook costs effort to write, and more to keep up to date. A plan for a
situation that never arises is waste, and an out-of-date plan can be worse
than none.

Be choosy, then. Write a playbook for a situation that recurs, or one that is
costly to get wrong: to customers, to money, or to a regulator. For a
situation that is rare and cheap to get wrong, decide when it arises. For a
situation close to one you have already planned for, extend that playbook
rather than starting another.

Then, after each use, ask the playbook's own Evaluation questions, and change
it. A playbook that has been used and revised a few times is worth far more
than one written once and filed.

## Keeping playbooks with the code

Playbooks belong where the team will find them, and where they change
alongside the systems they describe. For most software teams, that is the
repository. A playbook that mentions the payment job should sit in the same
pull request as the change that renames it.

detoo's `pse` tool keeps a log of PIOSEE playbooks as numbered Markdown files
in `docs/piosee`, the same way `adr` keeps architecture decisions:

```bash
pse init
pse new "Customer charged twice"
```

`pse init` starts the log, with a first playbook about when to write
playbooks at all. `pse new` creates the next one from a template with a
section for each step, ready for you to fill in. When a playbook is rewritten,
`pse new --replace` links the new one to the one it replaces, so the history
stays readable.

If you use a coding agent, the `create-piosee` skill will write the first
draft with you. It asks you about the problem, then researches the rest: the
code that handles the situation, the logs and dashboards where the evidence
lives, and the documentation of the systems involved, such as your payment
provider's. It asks you only for what it could not find, and policy, such as
who may issue a refund and up to how much, always comes from you. Like every
document the skills create, the playbook starts as a Draft, for the people
who will act on it to review.

The next time the email arrives at ten to six on a Friday, whoever picks it
up will know what to look at, what the choices are, and how to choose.
