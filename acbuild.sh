#!/usr/bin/env bash
set -e

name=redisservice
os=linux
version=0.0.3
arch=amd64

acbuildend () {
    export EXIT=$?;
    acbuild --debug end && exit $EXIT;
}

acbuild --debug begin
trap acbuildend EXIT

GOOS="${os}"
GOARCH="${arch}"

CGO_ENABLED=0 go build

acbuild set-name s-urbaniak.github.io/images/redisservice
acbuild copy "${name}" /"${name}"
acbuild set-exec /"${name}"
acbuild port add www tcp 8080
acbuild label add version "${version}"
acbuild label add arch "${arch}"
acbuild label add os "${os}"
acbuild annotation add authors "Sergiusz Urbaniak <sergiusz.urbaniak@gmail.com>"
acbuild write --overwrite "${name}"-"${version}"-"${os}"-"${arch}".aci

gpg --yes --batch \
    -u sergiusz.urbaniak@gmail.com \
    --armor \
    --output "${name}"-"${version}"-"${os}"-"${arch}".aci.asc \
    --detach-sign "${name}"-"${version}"-"${os}"-"${arch}".aci
