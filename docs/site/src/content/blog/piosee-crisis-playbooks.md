---
title: "Fly the aircraft: what pilots know about outages that we don't"
pubDate: 2026-10-17
author: "tom-regan"
draft: false
tags: [ "incident management", "SRE", "decision-making", "software development" ]
tool: pse
description: "When an engine fails after take-off, pilots know the engine is no longer the problem. I use PIOSEE, a model from the flight deck, to write incident playbooks."
heroImage: ../../assets/piosee-crisis-playbooks.jpg
---

A twin-engined airliner loses an engine just after take-off. The crew's first job is to keep it flying: straight ahead,
climbing, cleaned up, checklist done. Only then do they get to the problem they actually have to solve, which is landing
on the one engine they've got left. They'll need a longer landing distance than they planned for, and with the extra
drag they may not have the fuel to reach the airport they wanted.

That's the example Petter Hörnfeldt uses in a 2020 video about handling a crisis. Hörnfeldt is a Boeing 737 captain and
instructor who runs the Mentour Pilot YouTube channel, and the video walks through a decision-making model called
PIOSEE. He points out that pilots fly something like 99.99% of their flights the same way, by the same standard
procedures, and nobody judges them on those. They're judged on the rare flight where something goes wrong, which might
come many years into a career or never come at all.

Here is an engine failure in software. The primary database falls over at ten past three in the morning. That's the
thing that broke. The problem the engineer on call has to solve could be that checkout is down, or that writes are being
lost, or that a customer contract's SLA will be breached in forty minutes. A lost write wants a different fix from a
slow checkout, and the question of which one they're dealing with often doesn't come up until they're half an hour in.

## The meeting three days later

A few days after a bad outage there's usually a meeting. Somebody has built a timeline, the contributing factors are up
on a screen, and nobody is being blamed. The Google SRE book, which has a chapter called "Postmortem Culture: Learning
from Failure", says a postmortem is "not punishment—it is a learning opportunity for the entire company." ITIL's problem
management practice does a similar job with more forms, finding the causes of incidents and managing known errors so
there are fewer of them next time. I'm in favour of all of it. None of it is any use at ten past three.

At that hour the decisions fall to whoever is on call for the database. They know the system better than anyone else
awake, and they're the right person to be deciding, though nobody has told them how. In the hour between the page and the fix,
they have to work out what the problem really is, what they know about it, what they could do, which of those things to
try, and whether it's working. On most teams I've seen, they do that by phoning the most senior person they can wake up.

## Experience is doing all the work

Relying on the senior engineer isn't as daft as it sounds. Gary Klein, a psychologist who studied fireground commanders
and other experts, found that they rarely sit and weigh options against each other. They recognise the situation as one
they've met before, a workable course of action comes to mind, they run it through in their heads, and they act. He
called it recognition-primed decision making, and it's why a good senior engineer can glance at a graph and say
"connection pool" while everyone else is still scrolling through logs.

Recognition needs something to recognise, though, and the incidents that really hurt tend to be the ones nobody has seen
before. The expert has to be available, which means on the rota, awake, not on a beach and still working for you. If the
only way to get good at incidents is to live through lots of them, the skill ends up with whoever happened to be on call
on the worst nights. And at 3am the expert is as tired as everyone else. Stress usually narrows attention, and
experienced people latch onto the first hypothesis, or the loudest symptom, or the fix that worked last time.

When one of them gets through a bad night well, we say they're good in a crisis, and sooner or later we promote them.
Airlines couldn't wait for that kind of person to turn up in every cockpit. A 2016 paper on cockpit decision-making
begins: "Many case studies show that unstructured decision-making processes in teams are contributing factors to
accidents."

## What PIOSEE is

PIOSEE is one of several decision-making models airline pilots use. The letters stand for Problem, Information, Options,
Select, Execute and Evaluate. You define the problem you actually need to solve, gather the information that bears on
it, lay out the options, pick the one that best solves the problem, carry it out without dithering, and then check
whether it's working. If anything has changed, you go back to the problem. Hörnfeldt draws it as a circle.

T-DODAR (Time, Diagnose, Options, Decide, Assign, Review) covers almost the same ground in different words. FOR-DEC
(Facts, Options, Risks and benefits, Decision, Execution, Check) came out of a Lufthansa and German Aerospace Center
working group in the early 1990s, and it is the one that caught on most widely in Europe. Its designers put the hyphen
between the R and the D deliberately, to make the crew stop and ask whether they'd missed anything before committing.
Hörnfeldt reckons airlines use different tools but they're all much the same, and I agree. I use PIOSEE because the
words are plain and the order explains itself.

Before any of it comes a rule every pilot learns: fly the aircraft. Whatever else is going on, somebody keeps the plane
in the air. Hörnfeldt is blunt about panic. It came with us from the caves, he says, and wasn't much use even then;
running from a sabre-toothed tiger, it would have helped more to know which tree to climb.

In software, flying the aircraft means stopping the bleeding. Roll back, flip the flag, fail over, shed load or put up
the maintenance page, and diagnose afterwards. The 2016 paper concluded that FOR-DEC is useful "in complex situations
when there is enough time", and PIOSEE is the same, so it starts once the bleeding has stopped.

## Writing the playbook

A pilot has to run PIOSEE in the cockpit with the warning lights on. We can do most of ours weeks earlier, at a desk,
with a cup of tea.

A lot of what gets called a playbook is really a runbook: step-by-step instructions for one job, like failing over the
database. Runbooks are useful, but they assume somebody has already decided that failing over is the right move. Pilots
have checklists for the failures somebody foresaw, and the 2016 paper describes decision models as what airlines built
for "situations without any preconfigured solutions". A PIOSEE playbook is where the deciding goes. For each class of
incident, I write one with a section for each step.

**Fly the aircraft.** The stabilising actions that are safe to take before anyone knows the cause. Keep the list short
and keep everything on it reversible.

**Problem.** How this kind of incident tends to show up, written in terms of who is hurt rather than which component
failed. "Redis is down" describes a component. "Sessions are being dropped and logged-in users are being signed out"
describes a problem. At the top goes the first question to ask: what's the harm, who's affected, and how fast is it
growing?

**Information.** Where to look, and what to ignore. The specific dashboards, queries and log searches that answer the
questions the problem raises, and the people or teams who can see things you can't. It saves the usual scramble where
everyone on the call opens a different dashboard.

**Options.** The menu, worked out while everyone is calm: roll back, fail over, scale out, switch off a feature,
rate-limit a noisy client, change provider, restore from backup. Next to each one goes what it costs, how long it takes,
whether it can be undone and who is allowed to do it. Under pressure, people reach for whatever they used last time, and a
written list reminds them of the rest.

**Select.** The rules for choosing, agreed in advance. Reversible beats irreversible, and data integrity beats
availability. Service comes back before anyone goes hunting for the root cause. If I could keep only one section, it
would be this one, because choosing is where pressure does the most damage. With the rules agreed beforehand, the
on-call engineer is applying a policy the team has already signed off, which is much easier than making a judgement call
from nothing.

**Execute.** Who does it, how they announce it, and how the teams that depend on the system find out it's under way. The
relevant runbooks are linked here. Committing also means not starting a second fix while the first is still going in.

**Evaluate.** What success looks like and how long to wait for it: the metric that should move, the time by which it
should have moved, and what happens if it hasn't. A team without one can spend forty minutes watching a mitigation that
isn't working, hoping it will.

The on-call engineer still has to think. They get to spend that thinking on whatever is new about tonight, because the
predictable decisions were made weeks ago by people who weren't tired.

## One loop per system

Plenty of companies handle a big outage by putting one person in charge of all of it. The incident commander role comes
from the Incident Command System, which grew out of wildfire response in 1970s California, and Google's own incident
process is based on it. It answers a real problem, which is five engineers making five uncoordinated changes to
production at once.

I'd put the decisions somewhere else. A PIOSEE playbook belongs to a system, and the engineer on call for that system is
the person competent to run it. They know what the system can do, what each option costs and which of its dashboards
can't be trusted. An incident commander from another team, however experienced, has to ask them all of that before
deciding anything. Ownership deals with the five engineers as well. If only the payments on-call is allowed to change
payments, there's no one else in there to collide with.

Air traffic control can see every aircraft in its airspace and tells each crew plenty, but it doesn't fly any of them.
Under US rules, the pilot in command "is directly responsible for, and is the final authority as to, the operation of
that aircraft."

Say the checkout team's on-call sees card payments failing. Their problem is that customers can't pay, and their
Information step includes a question for the team upstream: is the payments service healthy? That report is what starts
the payments on-call's own PIOSEE, and they frame the problem in their own terms, perhaps as authorisations timing out
for one card provider. Each loop keeps its own Evaluate. When payments ships a fix, checkout's success rate is the
signal that tells both teams whether it worked.

That means the playbook needs two more things written into it. One is its triggers: the alerts that start it, and the
reports from downstream teams that should start it too. The other, in the Information section, is the list of upstream
systems it depends on and how to reach their on-call. With those written down, the checkout on-call knows exactly who to
message, and the payments on-call knows that a message from checkout is a reason to open their own playbook.

## Where the analogy bends

A cockpit has two pilots trained on the same procedures, sitting side by side. A system's on-call is usually one person,
on their own in the middle of the night, with no one to cross-check them. The playbook is the closest thing they have to
a first officer. The big incident call, with twenty people from teams that have never worked together and half of them
on mute, is what you get when nobody knows whose loop it is.

An aircraft is much the same machine next week as it is today. A software system isn't, and a playbook written for last
quarter's architecture will send people to dashboards that no longer exist. The meeting three days later is the
place to catch it. Add one question to the postmortem template, "Did the playbook help?", and change the playbook
whatever the answer, so the meeting fixes the next bad night's decisions as well as the system.

Most of a good options list is a senior engineer's pattern recognition, written down where the person on the rota can
find it. Klein's experts are usually right. There just aren't enough of them to go round.

## How I'd start

[detoo](../../) includes [`pse`](../../docs/pse/), a command-line tool that keeps PIOSEE
playbooks as Markdown files in your repository, and an agent skill, [create-piosee](../../docs/skills/create-piosee/),
that writes them with you.

1. Pick the three bad days you think are most likely for a system you're on call for. Your postmortems will tell you
   which.
2. [Install `pse` and the skill](../../docs/skills/install/), then ask your agent for the first playbook in a sentence,
   such as "Write a playbook for when the primary database falls over."
3. Answer its questions about the problem. Before it researches anything, it asks what happens, how you come to know of
   it and what has to be decided. Describe the harm rather than the component, and name the reports from downstream
   teams as well as the alerts, so the playbook says what starts it.
4. Ask it for a fly-the-aircraft section. Its template doesn't have one, but it adds any section you ask for. Keep
   everything on it reversible.
5. Check its research. It reads your code, your runbooks and the documentation of the systems involved, then proposes
   the information to gather, the options and the steps, citing a source for each and asking you to confirm it. Make
   sure it lists the upstream systems you depend on and how to reach their on-call, and gives each option its cost, how
   long it takes, whether it can be undone, who can approve it and the signal that will show it worked.
6. Agree the selection rules with the team, and with whoever has to live with the trade-offs. The skill won't set them
   for you: it proposes a rule only where a policy or a document already states one, and otherwise asks.
7. Have the people who'll follow the playbook review it. The skill creates every playbook as a draft, and only the team
   can say it's agreed.
8. Run a game day with someone junior on call, and see where the playbook falls apart.
9. After every real incident, ask whether the playbook helped, and change it. Once a playbook has been used, the skill
   can review it against its own Evaluate section.

Hörnfeldt says a decision model takes you from feeling like a victim of circumstances to sitting in the driver's seat.
The person who needs that most is the engineer on their first on-call shift, reading a message from the checkout team at
ten past three in the morning. Fly the aircraft first. It helps if somebody has already worked out where
you're going to land.

---

### Further reading

- Petter Hörnfeldt (Mentour Pilot), "HOW to deal with CRISIS.", YouTube (2020).
- Betsy Beyer, Chris Jones, Jennifer Petoff and Niall Richard Murphy (eds), *Site Reliability Engineering* (2016),
  chapters 14, "Managing Incidents", and 15, "Postmortem Culture: Learning from Failure".
- Henning Soll, Solveig Proske, Gesine Hofinger and Gunnar Steinhardt, "Decision-Making Tools for Aeronautical Teams:
  FOR-DEC and Beyond", *Aviation Psychology and Applied Human Factors* (2016).
- Gary Klein, *Sources of Power: How People Make Decisions* (1998).
- US Code of Federal Regulations, 14 CFR § 91.3, "Responsibility and authority of the pilot in command".
- AXELOS, *ITIL 4 Problem Management* practice guide.
