# Survey checklist

Examine every area below, in every part of the repository it applies to. For
each, record what you looked at and what you found: a candidate decision,
"nothing significant", or "not applicable" with the reason. The coverage table
in the report is built from these notes.

The questions are prompts, not a script. Follow what the code shows you.

## 1. Shape and boundaries

* What are the entry points: executables, services, functions, jobs, libraries?
* What are the deployable or distributable units, and how many are there?
* How is the code divided into modules, packages or layers? Is there a rule
  about which may depend on which, and does the code actually follow it?
  Check the imports, not only the directory names.
* Are there extension points: plugins, hooks, interfaces meant for others to
  implement?

## 2. Languages, runtimes and frameworks

* Which languages and versions, and why more than one, if so?
* Which frameworks shape the code: web, CLI, UI, dependency injection, ORM?
* Is any code generated? From what, and is the output committed?

## 3. Dependencies

* Read every manifest and lock file: `go.mod`, `package.json`,
  `pyproject.toml`, `requirements*.txt`, `Cargo.toml`, `pom.xml`,
  `build.gradle`, `Gemfile`, `*.csproj`, and the like.
* For each direct dependency, measure how far it reaches: count the packages
  or files that import it. A dependency imported everywhere is a candidate; one
  wrapped behind a single adapter usually is not.
* Is there a notable absence, such as no framework where one is usual, or a
  policy of keeping dependencies few? Look for evidence that it is deliberate.

## 4. Interfaces and contracts

* What does the system expose to users or other systems: a CLI, an HTTP or
  RPC API, a library API, events, files?
* What formats does it read and write: file formats, schemas, configuration
  files, environment variables?
* Is there a compatibility or versioning policy, stated or evident?

## 5. Data

* Where is state kept: databases, files, caches, object storage, memory?
* Who owns each schema or file format, and how does it change: migrations,
  versioned formats?
* What are the consistency and transaction boundaries?

## 6. Integration and communication

* What external services, APIs or systems does it call?
* Synchronous or asynchronous; which protocols; which message brokers?
* How are failures handled across those boundaries: retries, timeouts,
  idempotency?

## 7. Cross-cutting concerns

* Error handling: is there one strategy across the system, such as error
  types, wrapping, or mapping errors to exit codes or responses?
* Logging, metrics and tracing.
* Configuration: where it is read from, and in what order of precedence.
* Security: authentication, authorisation, secrets, input validation,
  sandboxing.
* Concurrency: threads, goroutines, async, actors, locking strategy.
* Time, locale, encoding.

## 8. Quality attributes

* Performance or resource strategies evident in the code.
* Portability: which operating systems, architectures, runtimes.
* Reliability, availability, recovery.
* Accessibility, internationalisation.

## 9. Construction

* Testing strategy: which layers are tested and how; harnesses, fixtures,
  golden files, fakes or mocks; end-to-end tests.
* Build: build tools, Makefiles, scripts, reproducibility.
* CI and CD: what runs, and what gates a merge or a release.
* Release and versioning: tags, changelogs, distribution channels.
* Repository layout: one module or many, monorepo conventions.
* Linting and formatting that enforce more than style.

## 10. Deployment and operations

* Packaging: binaries, containers, packages, installers.
* Infrastructure as code, environments, topology.
* How the system is installed, upgraded and rolled back.

## 11. Decisions written down elsewhere

Decisions often exist already, outside the log. Read:

* the README, especially any "Architecture" or "Design" section;
* design documents, specifications, RFCs, wikis checked into the repository;
* contributor guides;
* instructions for agents: `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, skills,
  editor rules;
* comments that explain why. Search for words such as "deliberately",
  "instead of", "rather than", "we chose", "decided", "trade-off", "because",
  "do not", "must not", "never".

## 12. History

Use read-only version control commands only.

* `git log --oneline | wc -l`, and the first commits, to understand the
  project's age and origins.
* Search commit messages for "migrate", "replace", "switch", "adopt",
  "introduce", "remove", "rewrite", "move to".
* Look for large structural commits: many files renamed or moved at once.
* `git log --follow` on a key file, to see when an approach arrived.

History is evidence of when code changed, not of why.
