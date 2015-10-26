source "${TEST_DIR}/logging.sh"

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

function gke-install {
  export CLOUDSDK_INSTALL_DIR="${HOME}"
  export CLOUDSDK_CORE_DISABLE_PROMPTS=1
  curl https://sdk.cloud.google.com | bash

  local gcloud_install_path="${CLOUDSDK_INSTALL_DIR}/google-cloud-sdk/bin"
  export PATH="${gcloud_install_path}:$PATH"
  gcloud -q components update kubectl
}

function gke-login {
  gcloud -q auth activate-service-account --key-file "${GCLOUD_CREDENTIALS_FILE}"
}

function gke-config {
  gcloud -q config set project "${GCLOUD_PROJECT_ID}"
  gcloud -q config set compute/zone "${K8S_ZONE}"
}

function gke-create-cluster {
  gcloud -q container clusters create "${K8S_CLUSTER_NAME}"
  gcloud -q config set container/cluster "${K8S_CLUSTER_NAME}"
  gcloud -q container clusters get-credentials "${K8S_CLUSTER_NAME}"
}

function gke-destroy {
  log-lifecycle "Destroying cluster ${K8S_CLUSTER_NAME}"
  if command -v gcloud &>/dev/null; then
    gcloud -q container clusters delete "${K8S_CLUSTER_NAME}" --no-wait
  fi
}

function setup-gke {
  gke-install
  gke-login
  gke-config
  gke-create-cluster
}

function run-test-plan {
  local test_plan

  # if is-pull-request; then
  #   test_plan="$(generate-test-plan)"
  # else
    test_plan="$(echo ${TEST_CHARTS} | tr ' ' '\n')"
  # fi

  log-lifecycle "Running test plan"
  log-info "Charts to be tested:"
  echo "${test_plan}"

  if [ ! -z "${test_plan}" ]; then
    setup-gke

    # Run all the tests in order
    echo "${test_plan}" | xargs -I_name -P1 -- sh -c './test/test-chart "_name"'
  fi
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
