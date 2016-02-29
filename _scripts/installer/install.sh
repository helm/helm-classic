#!/usr/bin/env bash

# Invoking this script:
#
# curl https://get.helm.sh | sh
#
# - download helm.zip file
# - extract zip file (into current directory)
# - making sure helm is executable
# - explain what was done
#

set -eo pipefail -o nounset

function check_platform_arch {
  local supported="linux-amd64 linux-i386 darwin-amd64 darwin-i386"

  if ! echo "${supported}" | tr ' ' '\n' | grep -q "${PLATFORM}-${ARCH}"; then
    cat <<EOF

${PROGRAM} is not currently supported on ${PLATFORM}-${ARCH}.

See https://github.com/helm/helm for more information.

EOF
  fi
}

function get_latest_version {
  local url="${1}"
  local version
  version="$(curl -sI "${url}" | grep "Location:" | sed -n 's%.*helm/%%;s%/view.*%%p' )"
  
  if [ -z "${version}" ]; then
    echo "There doesn't seem to be a version of ${PROGRAM} avaiable at ${url}." 1>&2
    return 1
  fi

  url_decode "${version}"
}

function url_decode {
  local url_encoded="${1//+/ }"
  printf '%b' "${url_encoded//%/\\x}"
}

PROGRAM="helm"
PLATFORM="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
HELM_ARTIFACT_REPO="${HELM_ARTIFACT_REPO:-"helm"}"
HELM_VERSION_URL="https://bintray.com/deis/${HELM_ARTIFACT_REPO}/helm/_latestVersion"
HELM_BIN_URL_BASE="https://dl.bintray.com/deis/${HELM_ARTIFACT_REPO}"

if [ "${ARCH}" == "x86_64" ]; then
  ARCH="amd64"
fi

check_platform_arch

VERSION="$(get_latest_version "${HELM_VERSION_URL}")"
HELM_ZIP="helm-${VERSION}-${PLATFORM}-${ARCH}.zip"

echo "Downloading ${HELM_ZIP} from Bintray..."
curl -Lsk "${HELM_BIN_URL_BASE}/${HELM_ZIP}" -O

echo "Extracting..."
unzip -qo "${HELM_ZIP}"
rm "${HELM_ZIP}"

chmod +x "${PROGRAM}"

cat <<EOF

${PROGRAM} is now available in your current directory.

To learn more about helm, execute:

    $ ./helm

EOF
