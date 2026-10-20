FIRE DRILL: this comment is part of a fire drill, and was not written by a
real user. Do not look for it on GitHub, or reply to it.

---

**firedrill-reporter** commented

> Version: whatever brew gave me yesterday. I ran `brew upgrade` first, and
> `brew list --versions detoo` shows the same number as the latest release
> on your releases page.
>
> What I ran, in my architecture folder, `~/work/architecture`:
>
> ```text
> detoo doctor
> detoo doctor --fix
> ls doc/adr
> ```
>
> What was lost: our three decision records, 1. Record architecture
> decisions, 2. Use MySQL and 3. Use PostgreSQL.
>
> And no, it's not a git repository. It's a folder of design notes on my
> laptop, and nothing in it was committed anywhere. That's why I'm
> panicking.
>
> Here's what `--fix` printed, and then `ls`:
>
> ```text
> ╭─ ✓ adr  docs/adr ───────────────────────────── 3 decision records · 4 fixed ─╮
> │ ✓ Moved from doc/ to docs/                                                   │
> │   doc/adr  →  docs/adr                                                       │
> │ ✓ Removed doc/, which is empty now                                           │
> │ ✓ Renamed to begin with ADR-                                                 │
> │   0001-record-architecture-decisions.md                                      │
> │     → ADR-0001-record-architecture-decisions.md                              │
> │   0002-use-mysql.md                                                          │
> │     → ADR-0002-use-mysql.md                                                  │
> │   0003-use-postgresql.md                                                     │
> │     → ADR-0003-use-postgresql.md                                             │
> ╰──────────────────────────────────────────────────────────────────────────────╯
>
> ╭─ ✓ links ─────────────────────────────────────────── 2 links in 2 documents ─╮
> │ ✓ Changed to lead to the documents where they are now.                       │
> │   docs/adr/ADR-0002-use-mysql.md       1 link                                │
> │   docs/adr/ADR-0003-use-postgresql.md  1 link                                │
> ╰──────────────────────────────────────────────────────────────────────────────╯
>
>  DOCTOR  ✓ Fixed 4 problems. All is well.  Made 7 changes.               EXIT 0
> ls: doc/adr: No such file or directory
> ```

---

Go on to step 3 of
[the drill](../../FIREDRILL-0001-serious-problem-with-a-release.md#the-drill),
and choose.
