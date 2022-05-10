#!/bin/bash

set -ue

CMD=unityguid
TEST_ROOT=/tmp/test

run_test() {
  local expect_file
  local subcommand
  expect_file=$1
  subcommand=$2
  shift;
  shift;
  fullcommand="$CMD $subcommand $@"

  if ! $fullcommand > "$TEST_ROOT/result"; then
    echo "❌ $fullcommand"
    exit 1
  fi

  if ! diff "$TEST_ROOT/result" "$expect_file"; then
    echo "❌ $fullcommand"
    exit 1
  fi

  echo "✅ $fullcommand"
}

run_setup() {
  mkdir -p $TEST_ROOT
}

run_teardown() {
  rm -rf $TEST_ROOT
}

create_meta_file() {
  local target_file
  local target_guid
  target_file=$1
  target_guid=$2

  cat << __GUID__ > "$target_file"
fileFormatVersion: 2
guid: $target_guid
folderAsset: yes
DefaultImporter:
  externalObjects: {}
  userData:
  assetBundleName:
  assetBundleVariant:
__GUID__
}

run_tests() {
  run_teardown
  run_setup

  local expect_file
  expect_file="$TEST_ROOT/expect"

  # list
  printf "1234\t/tmp/test/list/1.meta\n" > "$expect_file"
  printf "5678\t/tmp/test/list/2.meta\n" >> "$expect_file"

  mkdir "$TEST_ROOT/list"
  create_meta_file "$TEST_ROOT/list/1.meta" "1234"
  create_meta_file "$TEST_ROOT/list/2.meta" "5678"
  run_test "$expect_file" list "$TEST_ROOT/list"

  # conflict
  printf "0d3e014d4fe4741d3bb198eeaf4037a8	list1	list2	1b.cs.meta	2b.cs.meta\n" > "$expect_file"
  printf "0d3e014d4fe4741d3bb198eeaf4037a8	list1	list3	1b.cs.meta	3c.cs.meta\n" >> "$expect_file"

  mkdir "$TEST_ROOT/conflict"
  list1="$TEST_ROOT/conflict/list1.tsv"
  list2="$TEST_ROOT/conflict/list2.tsv"
  list3="$TEST_ROOT/conflict/list3.tsv"
  printf "0d3e014d4fe4741d3bb198eeaf4037a8\t1a.cs.meta\n" >  "$list1"
  printf "0d3e014d4fe4741d3bb198eeaf4037a8\t1b.cs.meta\n" >> "$list1"
  printf "0d3e014d4fe4741d3bb198eeaf4037a8\t2a.cs.meta\n" >  "$list2"
  printf "0d3e014d4fe4741d3bb198eeaf4037a8\t2b.cs.meta\n" >> "$list2"
  printf "0d3e014d4fe4741d3bb198eeaf4037a8\t3a.cs.meta\n" >  "$list3"
  printf "0d3e014d4fe4741d3bb198eeaf4037a8\t3c.cs.meta\n" >> "$list3"
  run_test "$expect_file" conflict "$list1" "$list2" "$list3"

  # replace
  mkdir -p "$TEST_ROOT/replace/Assets/Dir"{1,2,3}
  mkdir -p "$TEST_ROOT/replace/ProjectSettings"
  replace1="$TEST_ROOT/replace/Assets/Dir1/1.meta"
  replace2="$TEST_ROOT/replace/Assets/Dir2/2.meta"
  replace3="$TEST_ROOT/replace/Assets/Dir3/3.meta"
  guid1="0d3e014d4fe4741d3bb198eeaf4037a8"
  guid2="ad3e014d4fe4741d3bb198eeaf4037a8"
  guid3="bd3e014d4fe4741d3bb198eeaf4037a8"

  create_meta_file "$replace1" "$guid1"
  create_meta_file "$replace2" "$guid2"
  create_meta_file "$replace3" "$guid3"
  printf "$guid1 => eff4fd773896412687d92b61b19ad0ef	Assets/Dir1/1.meta\n" > $expect_file
  {
    fullcommand="$CMD replace $TEST_ROOT/replace"
    if ! $fullcommand "$guid1" | grep "$guid1" > /dev/null; then
      echo "❌ $fullcommand"
      exit 1
    fi
    echo "✅ $fullcommand"
  }

  run_teardown
}

if ! command -v $CMD > /dev/null; then
  echo "ERROR: not found command $CMD"
  exit 1
fi

run_tests
