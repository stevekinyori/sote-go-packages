#!/usr/bin/env bash
#
# This is a testing script for the seeding utility
#
# The following will generate coverage.out files and new packages should be added.
echo -n 'sLogger '
rm sLogger/coverage.out 2> /dev/null; go test sLogger/* -coverprofile sLogger/coverage.out
echo -n 'sError '
rm sError/coverage.out 2> /dev/null; go test sError/* -coverprofile sError/coverage.out
echo -n 'sDatabase '
rm sDatabase/coverage.out 2> /dev/null; go test sDatabase/* -coverprofile sDatabase/coverage.out
echo -n 'sConfigParams '
rm sConfigParams/coverage.out 2> /dev/null; go test sConfigParams/* -coverprofile sConfigParams/coverage.out
echo -n 'sAuthorize '
rm sAuthorize/coverage.out 2> /dev/null; go test sAuthorize/* -coverprofile sAuthorize/coverage.out

# The following will run tests that use the packages as an external package or program.
# There is no coverage because the tests are outside the directory with the source file.
go test sLogger_external_test.go
go test sError_external_test.go
go test sDatabase_external_test.go
go test sConfigParams_external_test.go
go test sAuthorize_external_test.go

# To output the results from the -coverprofile, using the following command
go tool cover -func=sLogger/coverage.out > coverage_review.out
go tool cover -func=sError/coverage.out >> coverage_review.out
go tool cover -func=sDatabase/coverage.out >> coverage_review.out
go tool cover -func=sConfigParams/coverage.out >> coverage_review.out
go tool cover -func=sAuthorize/coverage.out >> coverage_review.out

# Review the coverage totals for 70% compliance
grep '^total' coverage_review.out | awk '/[0-6][0-9]./ { print "FAILED: Coverage must be 70% or higher" }'; exit 1
grep '^total' coverage_review.out | awk '/[1,7-9][0-9]./ { print "PASSED: Coverage over 70%" }'
echo "For details, look at coverage_review.out"
