#!/bin/bash

mkdir testdata || true
rm testdata/* || true

tcs=("TestDefaultValue" "TestSet1" "TestSet2" "TestSet3" "TestFlush1" "TestFlush2" "TestFlush3" "TestDelete1" "TestDelete2" "TestDelete3" "TestCheck1" "TestCleanupRoutine1")

for tc in ${tcs[@]}; do go test -v -count=1 -failfast -cpu=4 -run="$tc" >> "testdata/$tc" 2>&1; done;

# go test -v -benchtime=10s -failfast -cpu=4 -race -benchmem -bench='^BenchmarkSet1$' -run='notmatch'
