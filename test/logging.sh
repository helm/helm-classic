# Text color variables
txtund=$(tput sgr 0 1)          # Underline
txtbld=$(tput bold)             # Bold
bldred=${txtbld}$(tput setaf 1) #  red
bldblu=${txtbld}$(tput setaf 4) #  blue
bldwht=${txtbld}$(tput setaf 7) #  white
txtrst=$(tput sgr0)             # Reset

pass="${bldblu}-->${txtrst}"
warn="${bldred}-->${txtrst}"
ques="${bldblu}???${txtrst}"

function log-lifecycle {
  echo "${bldblu}==> ${@}...${txtrst}"
}

function log-info {
  echo "${wht}--> ${@}${txtrst}"
}

function log-warn {
  echo "${bldred}--> ${@}${txtrst}"
}
