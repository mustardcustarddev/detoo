---
title: "Test your incident playbook with a fire drill before an incident tests it"
pubDate: 2026-10-20
author: "tom-regan"
tags: [ "incident management", "fire drills", "SRE" ]
categories: [ "Technology" ]
tool: fyr
description: "A playbook nobody has practised is a plan nobody has tested. Write a self-led fire drill for it, and find the gaps on a quiet afternoon."
heroImage: ../../assets/firefighters-tackling-blaze.png
---

A new issue arrives on a weekday morning:

> **doctor --fix deleted all our ADRs**
>
> Hi. We've kept our architecture decisions with adr-tools for years, and I wanted to try detoo, so I installed it
> with brew and ran `detoo doctor --fix` on our architecture folder, like the README says to.
>
> Now `doc/adr` is GONE. All of our decision records, just deleted. It even printed that it removed doc/.
>
> This is really bad. Please pull this release before it does this to someone else.

I maintain detoo, so this one is mine. What do I do first? Pull the release, as they ask? Ask them for more? Try it
myself? Or close it, because nothing in detoo deletes records?

This issue isn't real: it's the first page of a fire drill, and each of those is one of its choices. Two of them, done
for real, would do harm.

## The playbook behind it

In [my last post](../senior-to-staff-single-point-of-failure/), I argued for writing your incident judgement down as
PIOSEE playbooks, so the engineer on call can make the first decisions without waiting for the most senior person on
the team. detoo has one of its own, for this: *Serious problem with a release*. It says what to find out, what the
options are (revoke the release, fix forward, reply and close when detoo isn't at fault, or ask and wait), and the rules for
choosing between them.

I ended that post by saying you should practise a playbook and see where it falls apart. A fire drill is a way to do
that.

## A drill you can run on your own

A fire drill is one person meeting a simulated incident and responding to it as the playbook says. Nobody plays the
reporter or hands out the evidence, so it's written like a Choose Your Own Adventure book for engineers.

The drill's document sets the scene, then tells you what happens, a step at a time. Where the playbook asks for a
decision, the drill asks you to choose. You write down your choice and why, then turn to the page it names, which says
what happens next. When the drill is over, a `solution` directory holds the response the playbook leads to, and a
debrief to look back with.

Here's the page you get if you close that issue straight away:

> Had you done it: the reporter, who believes years of decisions have just been destroyed, with no copy anywhere, gets
> a reply that says it cannot have happened, and a closed issue. They do not reply. They never learn where their
> records are.

The page doesn't tell you that you were wrong. It tells you what it cost, and leaves the why for the debrief.

A drill does no harm, by design. It runs on a local copy of the system, run the way you'd run it for development, and
nothing is broken to make the incident. Any step that would reach outside your machine, such as posting a comment or
revoking a release, is written down, not taken. Every page that looks like a real message begins `FIRE DRILL:`.

## The answer is in the evidence

If you ask the reporter what happened, they reply, and paste what `detoo doctor --fix` printed:

```text
╭─ ✓ adr  docs/adr ───────────────────────────── 3 decision records · 4 fixed ─╮
│ ✓ Moved from doc/ to docs/                                                   │
│   doc/adr  →  docs/adr                                                       │
│ ✓ Removed doc/, which is empty now                                           │
│ ✓ Renamed to begin with ADR-                                                 │
```

Nothing was deleted. Doctor moved the records from `doc/adr` to `docs/adr`, renamed them, fixed the links between
them, and then removed the empty `doc/` folder. The reporter read "Removed doc/" and stopped reading.

The reporter knows what they saw, and says "deleted" when it was moved. Their words are the pressure. What they paste
is the evidence.

The drill doesn't ask you to take its word for it, either. It has you build the release the reporter has, from its
tag, in a scratch folder, copy in three records written the way adr-tools writes them, and run what they ran. It takes
about 45 minutes, most of it spent finding your way around and deciding.

## Writing one is the hard part

A drill is only worth doing if its evidence is real. A log line the code never writes, or an error a payment provider
never sends, makes a drill that practises a system you don't have. So the evidence has to come from the code, and from
the documentation of the other systems involved, not from memory.

detoo has a tool and a skill for it.

[`fyr`](../../docs/fyr/) keeps fire drills as Markdown in your repository, next to the playbooks `plm` keeps and the
SOPs `sop` keeps. `fyr new "Serious problem with a release"` creates a numbered drill with the playbook's title, and
the `problem` and `solution` directories beside it.

[create-firedrill](../../docs/skills/create-firedrill/) is an agent skill that writes the rest. It reads the playbook
and the SOPs it links to, then researches the evidence rather than inventing it: how the project runs locally, what the
code logs on the problem's path, often straight from the tests' golden files, and what the other systems involved
report, from their documentation. The doctor output above is the output of one of detoo's own tests. Then it plots the
choices, drawn from what a person under pressure would plausibly do, and asks you once before it writes anything.

It writes three kinds of drill. In a *reproduced* drill, like this one, you run the real system on the reporter's
inputs. In a *log-driven* one, you search the logs the drill holds for one identifier, across several services. In an
*alert-driven* one, an alert starts it and the evidence arrives as time passes.

## How I'd start

1. Pick the playbook that would hurt most to get wrong, or the one written most recently and read least.
2. [Install `fyr` and the skill](../../docs/skills/install/), and ask your agent for the drill in a sentence: "Write a
   fire drill for the release playbook." If there's no log of drills yet, it offers to start one with `fyr init`.
3. Check its plan when it asks. The expected response should be the playbook's, as it stands, even where you can see
   the playbook is wrong. That's what the debrief is for.
4. Give the drill to someone who didn't write the playbook, and keep them away from the `solution` directory.
5. Go through the debrief together, then take what it found back to the playbook with
   [create-playbook](../../docs/skills/create-playbook/).
6. Next time, write a drill for a different case of the same problem, such as the one where the evidence is unclear.

The [whole drill](https://github.com/mustardcustarddev/detoo/blob/main/docs/firedrill/FIREDRILL-0001-serious-problem-with-a-release.md)
is in detoo's repository, if you want to try being its maintainer for 45 minutes.

If a real issue like this one comes, I'd rather have found the playbook's gaps on a quiet afternoon, when getting it
wrong cost nothing.

---

### Further reading

- Open Practice Library, "[Fire Drills](https://openpracticelibrary.com/practice/fire-drills/)".
