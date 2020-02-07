#!/bin/bash

VERBOSE=${VERBOSE:-"0"}
V=""
if [[ "${VERBOSE}" == "1" ]];then
    V="-x"
    set -x
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

NAME=${1:?"app name"}
OUT=${2:?"output path"}
MUSES_SYSTEM=${3:?"version go package"}

set -e

GOOS=${GOOS:-linux}
GOARCH=${GOARCH:-amd64}
GOBINARY=${GOBINARY:-go}
GOPKG="$GOPATH/pkg"
BUILDINFO=${BUILDINFO:-""}
STATIC=${STATIC:-1}
LDFLAGS="-extldflags -static"
GOBUILDFLAGS=${GOBUILDFLAGS:-""}
GCFLAGS=${GCFLAGS:-}
export CGO_ENABLED=0

if [[ "${STATIC}" !=  "1" ]];then
    LDFLAGS=""
fi

# gather buildinfo if not already provided
# For a release build BUILDINFO should be produced
# at the beginning of the build and used throughout
if [[ -z ${BUILDINFO} ]];then
    BUILDINFO=$(mktemp)
    ${ROOT}/tool/version.sh ${NAME}> ${BUILDINFO}
fi

echo ROOT       "${ROOT}"
echo GOBINARY          "${GOBINARY}"
echo BUILDINFO          "${BUILDINFO}"
echo V          "${V}"
echo GOPKG          "${GOPKG}"
echo GOOS          "${GOOS}"
echo GOARCH          "${GOARCH}"
echo MUSES_SYSTEM  "${MUSES_SYSTEM}"

# BUILD LD_VERSIONFLAGS
LD_VERSIONFLAGS=""
while read line; do
    read SYMBOL VALUE < <(echo $line)
    LD_VERSIONFLAGS=${LD_VERSIONFLAGS}" -X ${MUSES_SYSTEM}.${SYMBOL}=${VALUE}"
done < "${BUILDINFO}"

# forgoing -i (incremental build) because it will be deprecated by tool chain.
time GOOS=${GOOS} GOARCH=${GOARCH} ${GOBINARY} build ${V} ${GOBUILDFLAGS} ${GCFLAGS:+-gcflags "${GCFLAGS}"} -o ${OUT} \
       -pkgdir=${GOPKG}/${GOOS}_${GOARCH} -ldflags "${LDFLAGS} ${LD_VERSIONFLAGS}"
