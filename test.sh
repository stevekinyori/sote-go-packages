#!/bin/zsh
#
# This is a testing script for Sote Packages
#
 set -uo pipefail
sons!
#
echo -n "Refreshing sotetest schema ."
`PGPASSWORD=$DB_PWSD psql -h $DB_HOST -d $DB_NAME -U $DB_USER -f db/migration/testdbstructures.sql 1> /dev/null 2> /tmp/packages_$$.out`
if [[ -s /tmp/packages_$$.out ]]; then
    echo "Done"
    echo "sotetest schema has been refreshed"
else
    echo "PSQL has FAILED."
    echo "    Investigate APP_ENVIRONMENT and AWS_REGION are correct."
    echo "    Check it the database is hosted on postgres and port 5432 is being used."
    exit 1
fi
#
# Display processing message
echo -n 'processing internal tests .'
#
# The following will generate coverage.out files and new packages should be added.
echo -n 'sLogger ' 1> /tmp/tmp_$$.out
rm sLogger/coverage.out 2> /dev/null; go test sLogger/*.go -coverprofile sLogger/coverage.out 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sError ' 1>> /tmp/tmp_$$.out
rm sError/coverage.out 2> /dev/null; go test sError/*.go -coverprofile sError/coverage.out 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sDatabase ' 1>> /tmp/tmp_$$.out
rm sDatabase/coverage.out 2> /dev/null; go test sDatabase/*.go -coverprofile sDatabase/coverage.out 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sConfigParams ' 1>> /tmp/tmp_$$.out
rm sConfigParams/coverage.out 2> /dev/null; go test sConfigParams/*.go -coverprofile sConfigParams/coverage.out 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sAuthorize ' 1>> /tmp/tmp_$$.out
rm sAuthorize/coverage.out 2> /dev/null; go test sAuthorize/*.go -coverprofile sAuthorize/coverage.out 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sMessage ' 1>> /tmp/tmp_$$.out
rm sMessage/coverage.out 2> /dev/null; go test sMessage/*.go -coverprofile sMessage/coverage.out 1>> /tmp/tmp_$$.out
echo "Done"
cat /tmp/tmp_$$.out
read RC <<< "$( grep '^FAIL' /tmp/tmp_$$.out | awk '/[F][A][I][L]./ {print 1}' )"
if [[ "$RC" == "1" ]]; then
    echo "FAILED: At least one test didn't pass."
    exit 1
else
    echo "PASSED: All tests passed"
fi
#
# Display processing message
echo -n 'processing external tests .'
#
# # The following will run tests that use the packages as an external package or program.
# # There is no coverage because the tests are outside the directory with the source file.
echo -n 'sLogger_external_test ' 1> /tmp/tmp_$$.out
go test sLogger_external_test.go 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sError_external_test ' 1>> /tmp/tmp_$$.out
go test sError_external_test.go 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sDatabase_external_test ' 1>> /tmp/tmp_$$.out
go test sDatabase_external_test.go 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sConfigParams_external_test ' 1>> /tmp/tmp_$$.out
go test sConfigParams_external_test.go 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sAuthorize_external_test ' 1>> /tmp/tmp_$$.out
go test sAuthorize_external_test.go 1>> /tmp/tmp_$$.out
echo -n '.'
echo -n 'sMessage_external_test ' 1>> /tmp/tmp_$$.out
go test sMessages_external_test.go 1>> /tmp/tmp_$$.out
echo "Done"
cat /tmp/tmp_$$.out
read RC <<< "$( grep '^FAIL' /tmp/tmp_$$.out | awk '/[F][A][I][L]./ {print 1}' )"
if [[ "$RC" == "1" ]]; then
    echo "FAILED: At least one external test didn't pass."
    exit 1
else
    echo "PASSED: All external tests passed"
fi
#
# Display processing message
echo -n 'processing test coverage .'
#
# # To output the results from the -coverprofile, using the following command
go tool cover -func=sLogger/coverage.out > coverage_review.out
go tool cover -func=sError/coverage.out >> coverage_review.out
go tool cover -func=sDatabase/coverage.out >> coverage_review.out
go tool cover -func=sConfigParams/coverage.out >> coverage_review.out
go tool cover -func=sAuthorize/coverage.out >> coverage_review.out
echo "Done"
#
# # Review the coverage totals for 70%+ compliance
read RC <<< "$( grep '^total' coverage_review.out | awk '/[0-6][0-9]./ {print 1}' )"
if [[ "$RC" == "1" ]]; then
    echo "FAILED: At least one components coverage is less than 70%"
    exit 1
else
    echo "PASSED: All components coverage is 70% or higher"
    exit 0
fi
echo "For details, look at coverage_review.out"