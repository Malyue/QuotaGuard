#!/use/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
#SCRIPT_ROOT=$(dirname $(readlink -f "$0"))/..

cd ${SCRIPT_ROOT}


#SCRIPT_ROOT=/Users/malyue/GolandProjects/mine/QuotaGuard
CODEGEN_PKG=${CODEGEN_PKG:-$(go list -m -f '{{.Dir}}' k8s.io/code-generator)}

echo "Generating code with kube_codegen.sh"

source "${CODEGEN_PKG}/kube_codegen.sh"

THIS_PKG="github.com/Malyue/quotaguard"
#    --input-dir "${SCRIPT_ROOT}/pkg/apis" \


#  ~/go/bin/deepcopy-gen -v 6
#  --output-file zz_generated.deepcopy.go --go-header-file ~/QuotaGuard/hack/boilerplate.go.txt github.com/Malyue/quotaguard/pkg/apis/quota/v1
kube::codegen::gen_helpers \
    --boilerplate "$(pwd)/hack/boilerplate.go.txt" \
    "${SCRIPT_ROOT}/pkg/apis"


kube::codegen::gen_client \
    --with-watch \
    --output-dir "${SCRIPT_ROOT}/pkg/generated" \
    --output-pkg "${THIS_PKG}/pkg/generated" \
    --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt" \
    "${SCRIPT_ROOT}/pkg/apis"

#bash "${CODEGEN_PKG}"/kube_codegen.sh \
#  -v=5 \
#  "deepcopy,client,informer,lister" \
#  github.com/Malyue/quotaguard/pkg/generated \
#  github.com/Malyue/quotaguard/pkg/apis \
#  "quota:v1" \
#  --output-base "${SCRIPT_ROOT}" \
#  --boilerplate "${SCRIPT_ROOT}/hack/boilerplate.go.txt"

