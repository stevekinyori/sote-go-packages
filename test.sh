#!/usr/bin/env bash
#
# This is a testing script for the seeding utility
#
# The following will generate coverage.out files and new packages should be added.
rm slogger/coverage.out; go test slogger/* -coverprofile slogger/coverage.out
rm serror/coverage.out; go test serror/* -coverprofile serror/coverage.out

# The following will run tests that use the packages as an external package or program.
# There is no coverage because the tests are outside the directory with the source file.
go test serror_external_test.go
go test slogger_external_test.go

# To output the results from the -coverprofile, using the following command
go tool cover -func=serror/coverage.out > coverage_review.out
go tool cover -func=slogger/coverage.out >> coverage_review.out

# Review the coverage totals for 70% compliance
 grep '^total' coverage_review.out | awk '/[0-6][0-9]./ { print "FAILED: Coverage must be 70% or higher" }'
 grep '^total' coverage_review.out | awk '/[1,7-9][0-9]./ { print "PASSED: Coverage over 70%" }'
