---
title: "AI DevEx: if your agent only writes code, it's writing slop"
pubDate: 2026-10-22
author: "tom-regan"
tags: [ "AI DevEx", "developer experience", "coding agents", "architecture decisions" ]
categories: [ "Technology" ]
tool: adr
heroImage: ../../assets/ai-devex-make-your-agent-the-librarian.png
description: "An agent that only writes code leaves code nobody can maintain. Put it to work keeping the reasons too, for the developers and agents who come next."
---

Anybody who has worked with me will likely have experienced my 'software engineering is librarianship' monologue - it has
worn out more developers than any career-grind at FAANG ever could. But I've recently wondered what librarianship looks like in
an age where the librarians are only skimming the books.

A pull request arrives on a Thursday afternoon. Fourteen hundred lines, tests passing, written that morning with a
coding agent. The reviewer gets halfway down and asks a reasonable question: why does this service keep its own cache,
when there's a shared one two directories away?

The author doesn't know. The agent chose it. So they open a fresh session and ask, and the agent gives them a good
answer: something about isolation and failure domains. It sounds right. It's also a guess (dare we even call it a lie?), because the session that made it has gone.

The cache gets merged, because the tests pass and the reason sounds fine. A year later, someone else's agent spots two
caches doing the same job and helpfully consolidates them.

That's what I mean by slop. Merriam-Webster made "slop" its word of the year for 2025, for low-quality content made in
bulk by AI, and most people think of it as the stuff filling their social feeds. I think most of it is in our
repositories. The code works. Nobody can explain it, and the next person to touch it can only accept it or change it
blindly.

## The code was never the hard part

In 1985, Peter Naur wrote an essay called *Programming as Theory Building*. His argument was that a program isn't really
its source code. It's a theory, held in the heads of the people who built it, of how the code maps onto the world and
why it's shaped the way it is. The code is what that theory leaves behind. "The death of a program," he wrote, "happens
when the programmer team possessing its theory is dissolved."

Read that with a coding agent in mind. When an agent writes most of a change, the theory behind it lives in a context
window, and the context window is dissolved at the end of the session. If you didn't make the decisions yourself, and
nobody wrote them down, the program is dead on arrival. It runs. Nobody holds its theory.

Michael Nygard made the practical version of this point in 2011, in the post that started architecture decision records.
Without the rationale, he said, a new person on a project has two choices: blindly accept a past decision, or blindly
change it.

A coding agent is that new person, every time. It joins your project from scratch in every session. A new human
colleague at least feels baffled or annoyed by a strange decision and goes to ask someone. An agent just picks one of
Nygard's two options, quickly and confidently, and moves on.

## What the numbers suggest

I don't want to lean too hard on statistics here. The research on AI and software quality is young, and the tools change
faster than anyone can study them. But the early signs point the same way.

GitClear analysed 211 million changed lines of code written between 2020 and 2024. From 2021 to 2024, the share of
changed lines that were cloned code rose from 8.3% to 12.3%, and the share associated with refactoring fell from a
quarter to under a tenth. In 2024, for the first time in their data, copied-and-pasted code outnumbered moved code. More code is
being added. Less of it is being reorganised into something someone could understand.

Google's DORA research found, in its 2024 report, that as AI adoption increased, it was accompanied by an estimated
1.5% drop in delivery throughput and a 7.2% drop in delivery stability. When Google announced the 2025 report, the
conclusion had sharpened: "AI doesn't fix a team; it amplifies what's already there." A team that already understands its system gets faster. A team
that doesn't gets to not understanding it faster.

## Give your agent the library

When Andrej Karpathy coined "vibe coding" in early 2025, he described giving in to the vibes and forgetting that the code
even exists. He said it was fine for throwaway weekend projects, and I agree with him. The trouble starts when the
weekend project has users, a team and a pager rota, and everyone is still working that way.

Use AI as much as you like. Just stop using it only as a typist. The most useful thing an
agent can do for a project that has to last is librarianship: cataloguing what was decided and why, noticing what's
missing from the catalogue, checking the shelves against it, and weeding out what no longer belongs. That's the work
that keeps a project open to the next developer, and it's exactly the work most teams never find time for.

There's a catch, and it matters more than anything else in this post. If you ask an agent to "document the
architecture", it will happily write you a reason for every decision in the codebase. Most of them will be invented,
like the reason for that cache. A made-up rationale is worse than none at all, because people believe it. Documentation
slop is still slop.

So the agent has to work like a librarian rather than a ghostwriter. It can find the decisions in the code, and gather
the facts around them, but the reasons have to come from the people who had them. Where nobody remembers, the record
should say so.

That's how I built detoo's [agent skills](../../docs/skills/). They write only what they have a source for, the code,
the documents in the repository or you, and ask you for the rest:

- [find-missing-adrs](../../docs/skills/find-missing-adrs/) reads the whole codebase and proposes records for the
  architectural decisions it embodies that nobody has written down. For each, it separates what the code shows from the
  reasons only you can supply.
- [create-adr](../../docs/skills/create-adr/) writes the record, and interviews you for why you chose this option over
  the others, rather than inventing an answer.
- [check-adrs](../../docs/skills/check-adrs/) checks the code against the decisions in force, rule by rule, with the
  evidence. Where they disagree, it offers two ways forward: change the code, or record a new decision that supersedes
  the old one. It never quietly edits the record to match the code.
- [review-adrs](../../docs/skills/review-adrs/) does the weeding. It goes through the log and asks of each record whether
  it was really an architectural decision, and offers to deprecate or supersede the ones that weren't.

The same goes beyond architecture. There are skills for RFCs, user stories, incident playbooks and SOPs, and
[`detoo doctor`](../../docs/guides/doctor/) keeps the shelves in order. Everything an agent writes starts as a draft.
Accepting it is for people to do.

## Some of your readers will be agents

The next developer to open your repository might be a new hire next spring. It might also be an agent in ten minutes'
time, yours or a colleague's. They need the same thing: the reasons, in plain text, next to the code, where they'll find
them without being told.

The industry has already noticed half of this. AGENTS.md, a file that tells coding agents how to work in a repository,
was one of the founding projects of the Linux Foundation's Agentic AI Foundation in December 2025, and OpenAI said at the
time that more than 60,000 open-source projects and agent frameworks had adopted it. DORA's 2025 AI Capabilities Model lists AI-accessible
internal data as one of seven capabilities that amplify the benefits of AI.

But an AGENTS.md tells an agent how to build and test your code. It doesn't tell it why the system is shaped the way it
is, and that's the part an agent is most likely to undo. Decision records do. detoo keeps them as numbered Markdown
files in `docs/adr`, linked to the records they supersede, which a person can read on GitHub and an agent can find with
`grep`. Nothing about them is specific to agents, which is the point. Write for the humans and the agents can read it
too.

## How I'd start

1. [Install detoo](../../docs/install/) and [the agent skills](../../docs/skills/install/).
2. Ask your agent: "Find the architecture decisions in this codebase that have no ADR." It'll come back with a list,
   and the evidence for each.
3. Pick the three you'd least like a newcomer to undo, and record them with `create-adr`. Answer its questions
   honestly. If nobody remembers why, say that, and let the record say it too.
4. Add a line to your AGENTS.md, or whatever your agent reads, telling it to read `docs/adr` before changing the
   architecture.
5. Before the next big change, ask: "Are we still following our architecture decisions?" Where the code has drifted,
   decide on purpose whether the code or the decision is wrong.
6. Run `detoo doctor` in CI, so the log stays where the tools and the agents expect it.

Code is now the cheapest thing your agent produces. The reasons behind it are the expensive part, and they're the part
that's gone when the session ends. Have the agent help you keep those too.

---

### Further reading

- Peter Naur, "Programming as Theory Building" (1985).
- Michael Nygard, "[Documenting Architecture Decisions](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions)"
  (2011).
- GitClear, "[AI Copilot Code Quality: 2025 Data](https://www.gitclear.com/ai_assistant_code_quality_2025_research)" (2025).
- Google Cloud, "[Announcing the 2024 DORA report](https://cloud.google.com/blog/products/devops-sre/announcing-the-2024-dora-report)" (2024).
- Google Cloud, "[Announcing the 2025 DORA report](https://cloud.google.com/blog/products/ai-machine-learning/announcing-the-2025-dora-report)" (2025).
- Google Cloud, "[Introducing DORA's inaugural AI Capabilities Model](https://cloud.google.com/blog/products/ai-machine-learning/introducing-doras-inaugural-ai-capabilities-model)"
  (2025).
