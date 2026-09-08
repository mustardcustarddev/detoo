# Bash completion for the adr command.
#
# Install this file into your bash completion directory, or source it from your
# shell's startup file.

_adr() {
  # Autocomplete only when entering the last term
  if [ "${COMP_CWORD}" -eq $((${#COMP_WORDS[@]} - 1)) ]
  then
    COMPREPLY=( $(adr autocomplete "${COMP_WORDS[@]}") )
  fi
}

complete -F _adr -o default adr
