#!/usr/bin/env bash
#
# This is a testing script for Sote Packages
#
# TODO add step to install the sotetest data structures
####  Add here
#
# Display processing message
echo 'processing...'
#
# The following will generate coverage.out files and new packages should be added.
echo -n 'sLogger ' 1> temp.out
rm sLogger/coverage.out 2> /dev/null; go test sLogger/* -coverprofile sLogger/coverage.out 1>> temp.out
echo -n 'sError ' 1>> temp.out
rm sError/coverage.out 2> /dev/null; go test sError/* -coverprofile sError/coverage.out 1>> temp.out
echo -n 'sDatabase ' 1>> temp.out
rm sDatabase/coverage.out 2> /dev/null; go test sDatabase/* -coverprofile sDatabase/coverage.out 1>> temp.out
echo -n 'sConfigParams ' 1>> temp.out
rm sConfigParams/coverage.out 2> /dev/null; go test sConfigParams/* -coverprofile sConfigParams/coverage.out 1>> temp.out
echo -n 'sAuthorize ' 1>> temp.out
rm sAuthorize/coverage.out 2> /dev/null; go test sAuthorize/* -coverprofile sAuthorize/coverage.out 1>> temp.out
cat temp.out
read RC <<< $( grep '^FAIL' temp.out | awk '/[F][A][I][L]./ {print 1}' )
if [[ "$RC" = "1" ]]; then
    echo "FAILED: At least one test didn't pass."
    exit 1
else
    echo "PASSED: All tests passed"
fi
#
# # The following will run tests that use the packages as an external package or program.
# # There is no coverage because the tests are outside the directory with the source file.
echo -n 'sLogger_external_test ' 1> temp.out
go test sLogger_external_test.go 1>> temp.out
echo -n 'sError_external_test ' 1>> temp.out
go test sError_external_test.go 1>> temp.out
echo -n 'sDatabase_external_test ' 1>> temp.out
go test sDatabase_external_test.go 1>> temp.out
echo -n 'sConfigParams_external_test ' 1>> temp.out
go test sConfigParams_external_test.go 1>> temp.out
echo -n 'sAuthorize_external_test ' 1>> temp.out
go test sAuthorize_external_test.go 1>> temp.out
cat temp.out
read RC <<< $( grep '^FAIL' temp.out | awk '/[F][A][I][L]./ {print 1}' )
if [[ "$RC" = "1" ]]; then
    echo "FAILED: At least one external test didn't pass."
    exit 1
else
    echo "PASSED: All external tests passed"
fi
#
# # To output the results from the -coverprofile, using the following command
go tool cover -func=sLogger/coverage.out > coverage_review.out
go tool cover -func=sError/coverage.out >> coverage_review.out
go tool cover -func=sDatabase/coverage.out >> coverage_review.out
go tool cover -func=sConfigParams/coverage.out >> coverage_review.out
go tool cover -func=sAuthorize/coverage.out >> coverage_review.out
#
# # Review the coverage totals for 70%+ compliance
read RC <<< $( grep '^total' coverage_review.out | awk '/[0-6][0-9]./ {print 1}' )
if [[ "$RC" = "1" ]]; then
    echo "FAILED: At least one components coverage is less than 70%"
    exit 1
else
    echo "PASSED: All components coverage is 70% or higher"
    exit 0
fi
echo "For details, look at coverage_review.out"