# Problem playbooks

A problem playbook is a plan, made in advance, for responding to a problem
the team expects to face, such as an incident, an outage or a customer
reporting an unexpected charge. It follows the PIOSEE method: the Problem to
decide, the Immediate action to take, the Counter-indicated actions to avoid,
the Information to gather, the Options, how to make the Selection, the steps
of Execution, and how to Evaluate the result. It refers to standard operating
procedures (SOPs) for its routine steps.

## Using a playbook

When the problem arises, find its playbook and follow it from the top. Take
the immediate action, if there is one, and none of the counter-indicated
actions. Record what was decided and done where the team records its
incidents. Afterwards, review the playbook: whether it led to a sound
decision in time, and what it lacked.

## Writing a playbook

Ask your agent, with the `create-playbook` skill, which `detoo skills
install` installs: "Write a playbook for a customer charged twice." It asks
you about the problem, and researches the code, the SOPs and the
documentation of the systems involved to propose the rest.

To write one yourself, run `plm new` with the problem, such as
`plm new "Customer charged twice"`. Write one for a problem that recurs, or
is costly to get wrong. A playbook is a draft until the people who will act
on it agree it. Practise it in a fire drill, with `fyr`.
