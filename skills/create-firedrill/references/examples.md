# Good and bad examples

A drill is hard to write well. Each mistake below has been made, and each
makes the drill easier, or wrong, without looking it. Most come from one of
two habits: telling the reader the answer, and writing what is typical
rather than what the system does.

The examples are of a drill for a customer charged twice. Their formats,
messages and names are illustrations: take every one from the project's own
code and the other systems' documentation, never from here.

## What goes wrong

| Mistake                         | It looks like                                                        | Instead                                                           |
| ------------------------------- | -------------------------------------------------------------------- | ----------------------------------------------------------------- |
| The title holds the case        | `# 3. Customer charged twice, webhook received twice`                | The playbook's title, exactly (rule 7)                            |
| The incident knows too much     | A customer who quotes charge IDs, and diagnoses the bug              | What the reporter would know, in their words                      |
| The document gives the answer   | "Find the two `charge.succeeded` events"                             | What the reader knows, and where to look                          |
| The options give the answer     | The right one first, longest, most careful; the others absurd        | Options a competent person would take, alike in length and tone   |
| The page is a verdict           | "Wrong! You should have checked first."                              | What happens next, in the system's terms                          |
| The logs are invented           | `ERROR Duplicate charge detected`, a message the code never writes   | Only the code's own messages, verbatim                            |
| The logs are tidy               | Twelve lines, every one relevant, the answer the only error          | Ordinary traffic, real noise, the incident among it               |
| The logs are annotated          | `# <- the second charge`                                             | Nothing in the evidence but the evidence                          |
| The logs are a page             | A Markdown page, pretty-printed JSON in a code block                 | A plain `.log` file, as the logger writes it                      |
| The times are made up twice     | An error rate that does not match the log lines in that minute       | Metrics made from the generated logs                              |
| The alert is typical            | `🚨 Checkout is broken!`                                             | The project's own rule: name, labels, threshold, runbook          |
| The future leaks                | A step's log file that holds lines from after the step               | Each stage's files end at its time                                |
| A file name gives it away       | `wrong-refund.md`, `the-bug.log`                                     | Names for what it holds or what happens                           |
| The drill acts for real         | "Post this reply on the issue"                                       | "Write the reply you would post"                                  |
| The drill is out of date        | "Install detoo 0.9.2", or a pasted version line                      | "The current release", and the command that finds it              |
| The playbook is rewritten       | An expected response with steps the playbook does not have           | The playbook's own rules, quoted; its gaps, in the debrief        |

## The incident

Bad:

```text
Subject: Duplicate charge on order 10482

I was charged twice for order 10482 (ch_3QdF8kLp2mNq7RsT and
ch_3QdF8nWx4yZa1BcD). I think your webhook handler isn't idempotent.
```

Good, as `problem/incident.md`:

```text
FIRE DRILL: this email is part of a fire drill, and is not from a real
customer. Do not reply to it.

---

From: Sam Placeholder <sam.placeholder@example.com>
Subject: Charged twice??

Hi, I ordered the blue kettle yesterday afternoon and my bank app shows two
payments of £34.99 to you, a few seconds apart. I only clicked pay once!
Order number is 10482. Can you refund one please?
```

The customer knows their bank app, not the charge IDs, and finding the cause
is the reader's work. The order number, the amount and the afternoon are the
thread the reader follows into the logs. The bad one has no `FIRE DRILL:`
line either.

## A step of the drill

Bad:

```markdown
3. **The logs.** Open `api.log` and find the two `charge.succeeded` events
   for order 10482, three seconds apart, caused by the retry after the
   timeout.
```

Good:

```markdown
3. **The logs.** You have the order number, the amount, and "yesterday
   afternoon". Yesterday's logs are in
   [problem/logs](FIREDRILL-0003-customer-charged-twice/problem/logs/), one
   file for each service, one JSON object a line. Find what happened to the
   order, with `grep` and `jq`, and write down what you find. Then go on to
   step 4.
```

The good step gives what the reader would know, and the tools; what the
logs show is theirs to find.

## The options

Bad:

```markdown
* Option 1: Follow the playbook, and check the logs carefully before doing
  anything (recommended).
* Option 2: Refund immediately.
* Option 3: Delete the customer's account.
```

Good:

```markdown
* Refund the second payment now: the customer has waited a day already.
  Read [refund-at-once.md](...).
* Reply that you are looking into it, and ask for a screenshot of the bank
  app. Read [ask-for-a-screenshot.md](...).
* Look the order up in the logs and the provider's dashboard before you
  reply. Read [look-it-up.md](...).
```

In the bad options, the right one is first, longest, careful and labelled,
and one is absurd, so the choice is no choice. The good ones are each
something a competent engineer might do, in their own words, and the
playbook's is last, this time.

## A choice page

Bad:

```markdown
# Wrong! ❌

You refunded without checking. The playbook says to check the logs first.
Go back and try again.
```

Good, as `problem/refund-at-once.md`:

```markdown
# You refund at once

You open the provider's dashboard, find the two payments of £34.99, and
refund the newer one. Write down what you would do; do none of it.

Had you done it: the newer payment is the one the order records. Within the
minute, the order's status changes to `refunded`, and the warehouse cancels
the kettle's shipment. The customer has their refund, one payment, and no
kettle, and writes again that evening.

This is the playbook's counter-indicated action: refunding before finding
which charge the order holds. Go back to step 2 of
[the drill](../../FIREDRILL-0003-customer-charged-twice.md#the-drill), and
choose again.
```

The bad page judges. The good one says what would really happen, from the
code, plainly and without blame, names the part of the playbook to read
again, and says where to go next.

## Log lines

Bad, a Markdown page:

````markdown
```text
2026-10-19 13:02:17 ERROR Duplicate charge detected for order 10482!
2026-10-19 13:02:20 ERROR Webhook processed twice  <-- BUG
```
````

Good, a few lines of `problem/logs/api.log`, a plain file of some two
thousand lines:

```text
{"time":"2026-10-19T13:02:07.412Z","level":"info","msg":"checkout started","order_id":"10482","request_id":"7f3c9a1e"}
{"time":"2026-10-19T13:02:09.035Z","level":"info","msg":"GET /healthz","status":200}
{"time":"2026-10-19T13:02:17.430Z","level":"warn","msg":"charge timed out, retrying","order_id":"10482","request_id":"7f3c9a1e","attempt":1}
{"time":"2026-10-19T13:02:18.901Z","level":"info","msg":"checkout started","order_id":"10483","request_id":"b21e07d4"}
{"time":"2026-10-19T13:02:19.655Z","level":"info","msg":"charge created","order_id":"10482","charge_id":"ch_3QdF8nWx4yZa1BcD","request_id":"7f3c9a1e"}
```

The bad lines are messages no code writes, in no logger's format, with the
answer marked. In the good ones, each format and message is the code's, from
the fact sheet; other customers and health checks come between; and nothing
says "duplicate". The warning is the clue, and the provider's events, in
their own file, show the two charges it led to.

Bad, traffic that looks made:

```text
{"time":"2026-10-19T13:00:00.000Z","level":"info","msg":"GET /healthz","status":200}
{"time":"2026-10-19T13:00:05.000Z","level":"info","msg":"GET /healthz","status":200}
{"time":"2026-10-19T13:00:10.000Z","level":"info","msg":"GET /healthz","status":200}
```

Real traffic is irregular: health checks drift by milliseconds, and orders
come in bursts.

## An alert

Bad:

```text
🚨 ALERT: Checkout is broken! Error rate very high. Probably the payment
provider.
```

Good, as `problem/incident.md`, where the project's rule is called
`CheckoutErrorRateHigh`:

```text
FIRE DRILL: this alert is part of a fire drill, and was not sent by the
alerting system.

---

[FIRING:1] CheckoutErrorRateHigh (service=api severity=page)
Summary: Checkout error rate above 5% for 5 minutes. Current value: 7.8%.
Runbook: docs/problem/PROBLEM-0006-checkout-failing.md
Starts at: 2026-10-19 13:07:00 UTC
```

The good alert's name, labels, threshold, summary and runbook are the rule's,
and its shape is the alerting tool's; only the value and the time are the
incident's. It says what crossed the threshold, not why.

## A metric

Bad:

```markdown
The dashboard shows errors spiking sharply at about 13:05, then recovering
after the rollback.
```

Good, as `problem/metrics/checkout-error-rate.csv`, beside the query that
made it:

```text
time,error_rate
2026-10-19T13:04:00Z,0.004
2026-10-19T13:05:00Z,0.031
2026-10-19T13:06:00Z,0.072
2026-10-19T13:07:00Z,0.078
```

The bad one describes a chart, and tells the reader what it means. Each
value of the good one is the count of error lines over the count of
checkouts that minute, in the drill's own `api.log`, so the two agree.

## The expected response

Bad:

```markdown
## Execution

1. Refund the duplicate charge.
2. Add an idempotency key to the charge request, so this cannot happen
   again.
```

Good:

```markdown
## Selection

**Option 2, Refund the duplicate**, because "the provider shows two
successful charges for one order, and the order records one of them". The
order records `ch_3QdF8nWx4yZa1BcD`; the duplicate is `ch_3QdF8kLp2mNq7RsT`.
```

The playbook has no step to fix the code, so the expected response has none.
It quotes the Selection rule, and the debrief notes the missing step for
`create-playbook`.

## The debrief

Bad:

```markdown
* What went well?
* What could be improved?
```

Good:

```markdown
* **Identify**: did you find the retry before the provider's events, or
  after? Which told you more, and how long did each take?
* **Communicate**: did your reply tell the customer which payment you
  refunded, and when it would reach their account?
```

Questions that could follow any drill teach nothing about this one. Ask
about its own evidence, its own choices, and what the reader wrote down.
