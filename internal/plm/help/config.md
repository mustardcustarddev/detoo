Prints the configuration that the tools are using, as shell variable
assignments, so that it can be inspected or evaluated by other scripts:

* `plm_dir`: the directory holding the problem playbooks, relative to the
  current working directory;
* `plm_template`: the project's own template for new playbooks, or empty if
  the built-in template is used.
