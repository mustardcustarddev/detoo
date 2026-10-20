# 3. Find the documentation on a website

Date: 2026-10-17

## Status

Accepted

## Story

As a member of the public using detoo, I want a documentation website for its
tools and skills, so that I can search for an answer instead of running the
help option in each tool.

## Acceptance criteria

* The site is published on GitHub Pages.
* It documents `adr`, `rfc`, `sty`, `pse` and `detoo`.
* It documents the agent skills.
* It explains how to install and set up the tools.
* It has a quick start with a few example commands.
* Each command's page shows the same text as the tool's built-in help, such
  as `adr help new`.

## Notes

Today the documentation is the README and each tool's built-in help, so a
user has to run the help option in each tool. Someone used to searching
documentation may struggle.

The built-in help is Markdown and example files under `internal/*/help/`,
built into each binary.

Out of scope: searching the site, which is a separate story.
