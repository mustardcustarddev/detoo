# Asking the user

## Use the question tool

If your tools include one for asking the user questions, every question you
put to the user goes through it: each round of an interview, each play-back,
each request for approval, and each choice the workflow offers. The tool is
named `AskUserQuestion` in Claude Code, `ask_user` in Gemini CLI, and
`request_user_input` in some other agents; use whichever you have.

These questions are decisions only the user can make, which is what such a
tool is for, even if its own description says to use it sparingly.

Write questions in the chat only when you have no such tool, or a call to it
fails. Then ask in the chat, and wait for the answer.

## Ask an open question through it

Most questions in an interview are open: who, why, what happens today. A tool
that needs options can still ask them. Build each question from these parts,
and only these:

* **The question**: the open question, in full, with any fact you found that
  prompts it. "Three roles can see invoices: `finance`, `admin` and
  `account-owner`. Who would export them?"
* **Options with evidence**: answers you can source, such as roles, records
  or options found in the repository, or already named by the user. In the
  example: `finance`, `admin`, `account-owner`.
* **Options about the question**, which guess nothing: the choices the skill
  gives for a fact the user cannot supply, such as "I don't know: leave it as
  an open question", or "Leave it out". A skill with a guide to writing lists
  them under "When a fact is missing".
* **The free-text answer.** Claude Code adds "Other" to every question, where
  the user types their own answer; say so in the question, as in "Choose Other
  to answer in your own words." If your tool has a free-text kind of question,
  as `ask_user` does, use that for an open question instead of options.

An open question with no evidenced answers therefore has two options about
the question, and its answer arrives as free text. That is the normal case,
not a workaround.

Each option's description says what choosing it will do, in the same terms
as its label. Like the label, it names nothing you have no evidence for: not
a person, a team or a source you have not found. "Say who could be asked" is
a description; "Ask the payments team" is a guess.

Never make up the number of options with guessed answers: guessed reasons,
guessed users, guessed causes. A user picking the nearest option is not
evidence.

## A round

* One call is one round: up to four related questions, in the order the
  interview plan gives. Under time pressure, one question a call.
* Put the facts you found, the draft so far, or anything long in the chat
  first, then make the call. The call holds the questions alone.
* Keep each header to a word or two; Claude Code allows twelve characters.
* Use a multiple choice only where several answers can be true together, such
  as which records to change.

## Play-backs and approval

* **Play-back**: the question is the sentence played back, "So: as a finance
  clerk, I need to export a month's invoices. Is that right?", with the
  options "Yes" and "Not quite", and the correction typed in as free text.
* **Approval**: write the draft and its sources in the chat, then ask "Shall
  I create it as it is?", with the options "Create it" and "Change it", and
  the changes typed in as free text.
* **A choice the workflow offers**, such as what to do with a record: one
  option for each choice the skill lists, your recommendation first, marked
  as yours.
