# You revoke the release

You open SOP 2, and follow it. Write down each step you would take, and the
notice and the README warning you would write, and do none of it.

Had you done it: the formula in the tap goes back to the release before, the
current release is deleted from GitHub with its tag kept, and a notice on the
documentation site, with a warning at the top of the README, tells every
user that the release is harmful, and to go back to the one before. Everyone
who follows it uninstalls the release, and installs the one before. The next
release, when there is one, has to take a later version.

The reporter's records are where `doctor --fix` left them. The release
before has the same `doctor --fix`, and does to an adr-tools folder exactly
what the current one does: everyone who follows the notice gets that
behaviour back, under an older version number.

Revoking a release for a problem detoo did not cause is the second of the
playbook's primary risks. Revoking cannot be taken back: SOP 2 keeps the tag,
and forbids publishing the release again.

The drill is over. Read [the solution](../solution/expected-response.md),
then [the debrief](../solution/debrief.md).
