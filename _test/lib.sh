function get-changed-charts {
  git diff --name-only HEAD origin/HEAD -- charts \
    | cut -d/ -f 1-2 \
    | sort \
    | uniq
}

function ensure-dirs-exist {
  local dirs="${@}"
  local pruned

  # ensure directories exist and just output directory name
  for dir in ${dirs}; do
    if [ -d ${dir} ]; then
      pruned+="$(basename "${dir}")\n"
    fi
  done

  echo -e "${pruned}"
}

function generate-test-plan {
  ensure-dirs-exist "$(get-changed-charts)"
}

function run-test-plan {
  local test_plan

  if is-pull-request; then
    test_plan="$(generate-test-plan)"
  else
    test_plan="${TEST_CHARTS}"
  fi

  log-lifecycle "Running test plan"
  log-info "Charts to be tested:"
  echo "${test_plan}"
  echo "${test_plan}" | xargs -I {} -P 1 echo \# helm install {} \&\& helm uninstall {}
}

function is-pull-request {
  if [ ! -z ${TRAVIS} ] && \
     [ ${TRAVIS_PULL_REQUEST} != false ]; then
    log-warn "This is a pull request."
    return 0
  else
    return 1
  fi
}

function log-lifecycle {
  echo "${bldblu}==> ${@}...${txtrst}"
}

function log-info {
  echo "${wht}--> ${@}${txtrst}"
}

function log-warn {
  echo "${bldred}--> ${@}${txtrst}"
}
