TEST_CHARTS="${TEST_CHARTS:-alpine}"

GCLOUD_CREDENTIALS_FILE="${GCLOUD_CREDENTIALS_FILE:-"${ROOT_DIR}/.gcloud-helm-test-credentials.json"}"

GCLOUD_PROJECT_ID="${GCLOUD_PROJECT_ID:-helm-test-1106}"
K8S_ZONE="${K8S_ZONE:-us-central1-b}"
K8S_CLUSTER_NAME="${K8S_CLUSTER_NAME:-${GCLOUD_PROJECT_ID}-$(openssl rand -hex 2)}"

export HELM_BIN="${ROOT_DIR}/helm.bin"
export TEST_DIR="${ROOT_DIR}/test"

export HEALTHCHECK_TIMEOUT_SEC=300
