#!/usr/bin/env bash
set -e

UNVER=`jq -r '.modules[0].resources["aws_s3_bucket.unversioned"].primary.id' terraform.tfstate`
VER=`jq -r '.modules[0].resources["aws_s3_bucket.versioned"].primary.id' terraform.tfstate`

runtest() {
  local test=${1}
  local bucket=${2}
  local expected=${3}
  local age=${4}
  local sleep=${5}

  printf "\nTEST: $test\n"

  aws s3 rm --quiet --recursive s3://$bucket/
  sleep $sleep

  ../s3undelete -bucket $bucket -age $age
  local actual=`aws s3 ls --recursive s3://$bucket/ | wc -l`

  if [[ $actual == $expected ]]; then
    echo "PASSED: Found $actual files in bucket $bucket and expected $expected."
  else
    echo "FAILED: Found $actual files in bucket $bucket but expected $expected."
    exit 1
  fi
}

runtest "Unversioned, no restores" $UNVER 0 1h 0
runtest "Versioned, all restored" $VER 5 1h 0
runtest "Versioned, no restores" $VER 0 1s 2
