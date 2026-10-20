# Kinds of drill

Drills differ in where their evidence comes from. Choose the kind from the
playbook: its Scenario says how the problem arrives, and its Information
says where the evidence is.

| Kind          | It starts with                                | The evidence                                                 | The reader's work                          |
| ------------- | --------------------------------------------- | ------------------------------------------------------------ | ------------------------------------------ |
| Reproduced    | a report from a person, such as an issue      | the real system, run locally on the reporter's inputs        | run what the reporter ran, and look        |
| Log-driven    | a report, or a question from support          | log files the drill holds, from each source                  | search, filter, and follow one identifier  |
| Alert-driven  | an alert, or a metric crossing a threshold    | the alert, then metrics, events and logs, as time passes     | act on part of the evidence, under time    |

* Choose **reproduced** when a local copy can show the problem without
  breaking anything: a command-line tool, a library, a migration, a build.
  The reader runs the real thing, on inputs the drill supplies.
* Choose **log-driven** when the problem happened in a running service and
  cannot be made locally without breaking something: a race, a third party's
  retry, a customer's own data.
* Choose **alert-driven** when the Scenario starts with an alert or a
  dashboard, and the problem goes on while the reader works.

Most drills are mainly one kind, with something of another: a customer's
email can start a drill whose work is in the logs. Name the kind in the plot
you play back in step 7, with why.

## Reproduced

detoo's own
[drill 1](https://github.com/mustardcustarddev/detoo/blob/main/docs/firedrill/FIREDRILL-0001-serious-problem-with-a-release.md)
is one: a user reports that `detoo doctor --fix` deleted their records, and
the reader builds the current release and runs it on records like theirs.

* **The inputs are files in the drill.** Copy them from the project's test
  fixtures, where there are some, into a directory of their own under
  `problem`, such as `problem/architecture/`. The step that uses them copies
  them to a scratch folder outside the working copy first.
* **The version is built, not installed.** The reader builds the release in
  question from its tag, in a git worktree outside the working copy, and runs
  it by its path. Installing it would change their own tools, and building in
  the working copy would change their work. Give the command that finds the
  tag, not the tag.
* **Say where a command that changes files may run.** Name the scratch
  folder in the ground rules, and say why: a tool that acts on the nearest
  project above where it runs can reach the working copy from a folder
  inside it.
* **What the reporter pastes is the real output**, for the commands they
  ran: take it from the project's golden files or tests, and leave out any
  line that names a version. Do not run the commands to get it unless the
  user asks you to (rule 1).
* **The reporter is wrong in a human way.** A report of what was lost or
  broken is often a misreading: "deleted" for moved, "broken" for changed.
  The reporter's words carry the misreading, and the output they paste
  carries the truth, for the reader who reads it. This is the case to choose
  when the playbook has an option for a problem the system did not cause.
* **The step says what to do, not what will be seen**: "Run what they ran,
  with the release you built. Look at what is in the folder now." The reader
  sees the result on their own machine, and the expected response says what
  it should be.

## Log-driven

The reader spends most of the drill in log files, as in a real incident:
finding the customer's request, following it from service to service, and
noticing what is not there. The drill holds the logs, and you make them.

* **Plain files, in the real format.** Write each source's logs as the file
  the system would write, such as `problem/logs/api.log`: one JSON object a
  line, if that is what the logger writes, never pretty-printed, and never a
  Markdown page with a code block. The reader works on them with the tools
  they would use, such as `grep`, `jq` and `less`. The document names the
  tools, not what to search for.
* **Enough to search.** About an hour of ordinary traffic around the
  incident: hundreds to a few thousand lines a source, of the messages the
  code writes on its ordinary paths, such as health checks and other
  customers' orders. Enough that reading top to bottom is not the way; not
  so much that the drill is about the tools.
* **Generate the traffic, and place the incident by hand.** Write a small
  script, in your scratchpad, not in the drill. It writes the ordinary lines
  from the real messages, with identifiers made up in the real formats,
  times with the irregular gaps and bursts of real traffic, and a fixed
  random seed. Then add the incident's lines yourself, at their times, among
  them. The script stays out of the drill: it holds the answer.
* **Noise from the same code.** Unrelated warnings and errors the system
  really writes, such as a timeout to another service or a client's bad
  request, appear as they would on any day. The incident's lines are not the
  only errors, and are often not errors at all: a double charge is two
  successful charges.
* **The evidence is what is there, and what is not.** Incidents rarely
  announce themselves: no line says "duplicate charge". The reader finds a
  retry, then two `charge.succeeded` events for one order, seconds apart.
  Where the code writes nothing, as when a webhook is received twice, the
  logs show nothing, and the reader has to notice it.
* **Each source keeps its own clock.** If the web server logs local time and
  the application UTC, so do the files, and the reader lines them up, as
  they would in a real incident. The times themselves agree.
* **One identifier leads through.** The identifiers on the incident page
  lead from file to file: the customer's order number to the request ID, the
  request ID to the charge IDs, the charge IDs to the provider's events.

## Alert-driven

The drill starts with a machine saying that something is wrong, before
anyone knows what, and the problem goes on while the reader works. It is a
timeline: each step has a time, and says what has changed since the last.

* **The alert is the project's own.** Find the alerting rules in the
  repository, such as Prometheus rules, cloud alarms in Terraform, or
  monitors kept as code. Write the alert as the alerting tool delivers it:
  its name, severity, labels, summary and runbook link, from the rule, and
  the value that crossed the threshold, from the incident. If the alert the
  playbook relies on is defined nowhere, it does not exist: that is a
  finding for the debrief, and the drill starts the way the problem would
  really be found, such as by a customer.
* **Metrics as a query returns them.** A metric is a file of values, such as
  a CSV of the query's result at one-minute steps, with the query beside it:
  never a description of a chart. Use the metric names and labels the code
  exports.
* **Metrics agree with the logs.** Make each metric from the generated logs,
  such as the error rate from the count of error lines and requests each
  minute, rather than writing both by hand. A reader who checks one against
  the other finds that they agree.
* **It unfolds.** Reveal the evidence in stages: at 13:07 the alert fires;
  at 13:12 it is still firing, the rate has doubled, and support has two
  tickets; at 13:20 the provider posts an incident. Each stage's files hold
  only what exists by then, such as `logs/13-07/api.log`, then
  `logs/13-12/api.log`, going on from it.
* **What the reader does changes what happens.** After a choice, the page
  says how the system responds, in its metrics: a rollback that removes the
  cause, and the rate falls within minutes; a restart that does not, and it
  falls, then climbs again; more instances, for a fault in a dependency, and
  nothing changes. Each response follows from the cause the case chose.
* **Pressure, honestly.** No one is watching a self-led reader's clock, so
  ask them to keep one: "Give yourself five minutes for this step, then
  choose." The pressure the playbook expects, such as a manager asking for
  an update, arrives as a message on a page.
* **The immediate action is the first decision.** An alert is where a
  playbook's immediate action, such as failing over, is most often taken,
  and most often taken wrongly. Ask whether to take it at the first step,
  before the evidence is in.
