---
title: Custom templates
description: Give a log its own template for new documents, or choose one for a single document, and see which template a tool is using.
section: guides
order: 3
---

Each tool ships with a template for new documents. A project can replace it
with its own.

## A template for the log

Put the template in `templates/template.md` inside the log directory, such as
`docs/adr/templates/template.md`. Every new record in that log is created from
it.

The tool fills in four placeholders, the first time each appears on a line:

| Placeholder | Becomes |
| --- | --- |
| `NUMBER` | The record's number, such as `12`. |
| `TITLE` | The title given to `new`. |
| `DATE` | Today's date, such as `2026-10-17`, or the tool's `DATE` variable, such as `ADR_DATE`. |
| `STATUS` | The status a new document starts with, `Draft`. |

`adr`'s built-in template is a good place to start:

```markdown
# NUMBER. TITLE

Date: DATE

## Status

STATUS

## Context

The issue motivating this decision, and any context that influences or constrains the decision.

## Decision

The change that we're proposing or have agreed to implement.

## Consequences

What becomes easier or more difficult to do and any risks introduced by the change that will need to be mitigated.
```

Keep the `## Status` section: `new --supersede` and `link` rewrite it.

## A template for one document

The tool's `TEMPLATE` variable, such as `ADR_TEMPLATE`, names a template for
the documents created while it is set, in place of the log's own:

```bash
ADR_TEMPLATE=templates/short.md adr new "Use UTC throughout"
```

## Which template is in use

```bash
adr config
```

`config` prints the log directory and the project's template, or nothing for
the template if the built-in one is in use. It does not report the `TEMPLATE`
variable, which overrides both.
