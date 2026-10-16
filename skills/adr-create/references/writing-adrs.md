# Writing an ADR

An ADR is written as if it were a conversation with a future developer: one
who was not in the room, who finds the decision in the code, and who has to
decide whether to keep it. They can only judge it well if everything in the
record is true. A plausible reason that nobody actually had is worse than no
reason, because it will be believed.

## Facts first

Before writing any prose, build a fact sheet: every claim the record will make,
and where it came from. Keep it in your working notes. It is shown to the user
with the draft, and is not written into the record.

| # | Claim | Source |
| --- | --- | --- |
| 1 | Records are stored as Markdown in `docs/adr` | `adr config` output |
| 2 | The team needs to run on air-gapped machines | User, in this conversation |
| 3 | `internal/store` is imported by 14 packages | `grep -rl` count, run today |
| 4 | Why SQLite was preferred to PostgreSQL | **Unknown — ask** |

### Acceptable sources

* **The user**, in this conversation: what they stated, or explicitly agreed to.
  A suggestion of yours that they did not take up is not a fact.
* **The repository**: code, configuration, documents and existing ADRs that
  you have read. Note the path, and the line where it matters.
* **Commands you ran**, and their output: counts, dependency lists, test
  results, `git log` history.
* **A document the user pointed you to.**

### Not acceptable

Do not write any of these unless a source above supports it:

* a motivation, a problem, or a pressure the team was under;
* an alternative that was considered, or why it was rejected;
* a benchmark, a measurement, a size, a count or a date;
* a stakeholder, a team, an incident or a requirement;
* a prediction: that something "will improve", "scales better", "reduces risk";
* a general claim such as "industry standard" or "best practice".

Code shows what was done, not why. An inferred reason is a guess, however
likely it seems. Ask.

### When a fact is missing

Ask the user. If they do not know, or prefer not to say, do one of these, and
let them choose:

* **Say so plainly.** "The alternatives considered at the time were not
  recorded." An honest gap tells the future reader where the record's knowledge
  ends.
* **Leave it out**, where the template does not need it.

Never fill the gap with something plausible.

## The sections

The built-in template has a title, a date, a status, and three sections. A
project's own template may differ: fill in each of its sections by the
guidance written in it, with the same standard of evidence.

### Title

A short phrase naming the decision, such as "Use PostgreSQL for persistence"
or "Share one engine between the record tools". Run `adr list` and follow the
style of the existing titles. The title becomes the file name, so keep it
short.

### Date

The tool writes today's date. When recording a decision that was made earlier,
ask the user which date they want: the day it is recorded, or the day it was
made. Set `ADR_DATE` for the latter. A date taken from `git log` is evidence of
when the code changed, which is not necessarily when the decision was made; say
which it is.

### Status

The tool writes `Accepted`. If the decision has not been agreed yet, ask
whether it should be `Proposed`, and edit the line. Do not decide this
yourself: whether something is agreed is a fact only the user knows.

### Context

The forces at play: technical, political, social, and local to the project,
which are probably in tension, and should be called out as such. The language
is value-neutral. It describes facts, and argues for nothing.

* State the problem that made a decision necessary, and why then.
* State the constraints that limit the options.
* If alternatives were considered, list each with the trade-offs that were
  actually weighed. A short bulleted list suits alternatives. Do not add
  alternatives nobody considered, to make the decision look thorough.

### Decision

The response to those forces, in full sentences and the active voice: "We
will …".

* Name concrete things: packages, directories, technologies, commands, rules.
  A decision precise enough to check against the code is one that can be kept.
* One decision per record. If you find two, suggest two records.
* Write only what was decided. Do not extend it with rules the user did not
  agree to.

### Consequences

The resulting context, after the decision is applied: what becomes easier,
what becomes harder, and what risks need mitigating. All of them, not just the
positive ones.

* A consequence that follows mechanically, and can be checked, may be stated
  from the evidence: "This adds a dependency on X" when `go.mod` shows it.
* A consequence that is a judgement or a prediction needs the user's
  agreement: "Builds will be slower", "This will be hard to change later".
* Where it is known, say what would be the signal to revisit the decision.

## Style

* Full sentences, organised in paragraphs. Bullets suit lists, such as
  alternatives or risks, but are not an excuse for fragments.
* One or two pages at most. Short records get read and kept up to date.
* Plain words. Prefer "we chose X because Y" to "X was leveraged to enable Y".
* No marketing: "robust", "seamless", "modern", "clean", "powerful",
  "significantly", "simply", "obviously" add nothing a reader can check.
* Keep the user's meaning. Do not turn "it seemed easier" into "it reduces
  complexity", or "we might need it" into "it is required".
* No placeholders and no template guidance left behind.

## Examples

| Instead of | Write, if the evidence supports it |
| --- | --- |
| PostgreSQL is the industry-standard choice for scalable systems. | We already run PostgreSQL for the billing service, and the team knows how to operate it. |
| This decision greatly improves maintainability. | Each tool is now a definition of about fifty lines; a fix to the engine fixes both. |
| We evaluated several options. | We considered two alternatives: a full copy of each tool, and a shared domain with a stack for each tool. |
| Performance will be better. | *(Ask: is there a measurement? If not, leave it out.)* |
| The team decided after careful analysis. | *(Leave it out: it says nothing the reader can use.)* |

## Before you finish

* Every sentence traces to a line of the fact sheet.
* Every "Unknown — ask" line was asked about, and the answer is used, stated
  as a gap, or left out by the user's choice.
* Every section of the template in use is filled in, and no guidance text is
  left.
* Links between records were made by the tool, not by hand.
* The user has seen the draft with its sources, and approved it.
