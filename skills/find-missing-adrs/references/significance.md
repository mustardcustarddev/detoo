# Architectural significance

> Decisions made at the right level, when they're needed.

An ADR log is only useful while it is small enough to read and every record in
it matters. A record of a local design choice, or of something trivial, hides
the records that matter. A record written before the decision was needed
records a guess, and constrains work nobody yet understands.

This guide is used to judge whether a decision deserves an ADR. Every ADR
skill that judges significance carries an identical copy, so that a record
created by one is judged the same way by another.

## Three levels

### Architectural

A decision is architectural when it shapes the system as a whole, so that
changing it later means changing many parts of the system, or changing a
contract that others depend on. The first ADR that `adr init` writes names the
categories: decisions that affect the **structure**, **non-functional
characteristics**, **dependencies**, **interfaces**, or **construction
techniques**.

Signs that a decision is architectural:

* It constrains how the parts of the system are built or how they relate:
  layering, the direction of dependencies, module or process boundaries.
* It fixes a contract consumed outside the component that owns it: a public
  API, the command-line surface, a wire protocol, a file or data format, a
  shared schema, an event format.
* It selects a technology that reaches throughout the system: the language, a
  framework, a datastore, a messaging system, or a dependency used everywhere
  rather than behind one adapter.
* It determines a quality attribute: the security model, consistency,
  performance or availability strategy, the concurrency model, how errors cross
  boundaries.
* It sets a technique the whole team follows: the testing strategy, the build
  and release process, code generation, the deployment topology.

### Design

A design decision is real, and may be important, but it is local. It lives
inside one component, and could be changed by refactoring that component
without changing any contract that other components or users rely on.

Examples: an internal data structure or algorithm; a pattern used inside one
package; a library wrapped behind an interface and used in one place; how a
package is split into files; naming within a module.

The right record for a design decision is usually next to the code: a doc
comment, a test, a package README, a commit message.

### Trivial

The choice follows from a convention, a default, or the only sensible option;
no viable alternatives were in tension; or it is a matter of taste.

Examples: formatting enforced by a formatter; a version bump; the language's
standard project layout; a tool that the ecosystem already mandates.

## The tests

Apply every test, and write down the answers: they are your reasoning, and the
user will want to see them. Wherever a test can be answered from the code,
answer it from the code, and cite what you found. Do not answer it from the
ADR's own description of itself.

1. **Reversal cost.** If the team changed its mind next quarter, what would it
   cost? Hours of work in one package suggests design. Work across many
   packages, a data migration, a breaking change for users, or retraining the
   team suggests architecture.
2. **Blast radius.** What would have to change: how many files, packages,
   services, teams, external consumers? Measure it: search for the imports,
   call sites and references.
3. **Refactoring.** Could the decision be reversed by a behaviour-preserving
   change confined to one component, behind an interface that stays the same?
   If so, it is not architectural.
4. **Real alternatives.** Were there at least two viable options, with
   trade-offs between them? If there was only one sensible choice, there was no
   decision to record.
5. **Future reader.** Would a new developer who found this in the code be
   puzzled by it, or tempted to undo it, without the record? If the code
   explains itself, the record adds little.
6. **Timing.** Was the decision needed when it was made? A decision about a
   feature, a scale, or an integration that does not exist yet is speculative.
   The principle cuts the other way too: a decision already embodied throughout
   the code, and still unrecorded, is overdue.

## Judgement

* **Scale is relative.** A small command-line tool has a small architecture.
  Judge a decision against this system, not against a large distributed one:
  ask whether it is among the decisions that shape *this* system.
* **The same technology can sit at either level.** A database chosen for the
  whole system is architectural; a JSON library used in one adapter is design.
  The difference is how far it reaches into the code, which you can check.
* **Levels change over time.** A choice that once reached everywhere may since
  have been isolated behind an abstraction, and become design. A local choice
  may have spread until it is architectural. Judge the decision as the code
  stands now, and say if that differs from when it was made.
* **Borderline is a legitimate answer.** Say which tests point which way, and
  let the user decide. Do not force a verdict in order to look decisive.
* **Project-wide process decisions count as construction techniques.** The
  decision to keep ADRs at all, a branching model, or a release process can be
  architectural if the whole project follows it.
