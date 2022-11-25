#!/bin/zsh
#
# This is a testing script for Sote Packages
#
set -uo pipefail
#
# Display Checking system message
echo -n "Checking system config ."
#
# Checking to make sure environment variables are set correctly
read RC1 <<<"$(printenv | grep -i 'APP_ENVIRONMENT' | awk '/[development]./ || /[staging]./  || /[demo]./ || /[production]./ {print 0}')"
# REMOVED ~/.aws/config contains region
# read RC2 <<< "$( printenv | grep -i 'AWS_REGION' | awk '/[eu-west-1]./ {print 0}' )"
#if [[ "$RC1" == "" || "$RC2" == "" ]]; then
#    echo "SYSTEM IS NOT CONFIGURED: You must have APP_ENVIRONMENT and AWS_REGION set as environment variables."
#    echo "    APP_ENVIRONMENT must be 'development', 'staging', 'demo' or 'production'"
#    echo "    AWS_REGION must be 'eu-west-1'"
if [[ "$RC1" == "" ]]; then
  echo "SYSTEM IS NOT CONFIGURED: You must have APP_ENVIRONMENT set as environment variables."
  echo "    APP_ENVIRONMENT must be 'development', 'staging', 'demo' or 'production'"
  exit 1
fi
#
# Checking to make sure the sotetest schema is installed
DB_HOST=$(aws ssm get-parameters --with-decryption --names /sote/api/$APP_ENVIRONMENT/DB_HOST --query "Parameters[*].{Value:Value}")
DB_HOST=$(echo $DB_HOST | jq '.[] | .Value' | sed "s/^\([\"']\)\(.*\)\1$/\2/g")
echo -n "."
DB_PORT=$(aws ssm get-parameters --with-decryption --names /sote/api/$APP_ENVIRONMENT/DB_PORT --query "Parameters[*].{Value:Value}")
DB_PORT=$(echo $DB_PORT | jq '.[] | .Value' | sed "s/^\([\"']\)\(.*\)\1$/\2/g")
DB_NAME=$(aws ssm get-parameters --with-decryption --names /sote/api/$APP_ENVIRONMENT/DB_NAME --query "Parameters[*].{Value:Value}")
DB_NAME=$(echo $DB_NAME | jq '.[] | .Value' | sed "s/^\([\"']\)\(.*\)\1$/\2/g")
echo -n "."
DB_USER=$(aws ssm get-parameters --with-decryption --names /sote/api/$APP_ENVIRONMENT/DB_USERNAME --query "Parameters[*].{Value:Value}")
DB_USER=$(echo $DB_USER | jq '.[] | .Value' | sed "s/^\([\"']\)\(.*\)\1$/\2/g")
DB_PWSD=$(aws ssm get-parameters --with-decryption --names /sote/api/$APP_ENVIRONMENT/DATABASE_PASSWORD --query "Parameters[*].{Value:Value}")
DB_PWSD=$(echo $DB_PWSD | jq '.[] | .Value' | sed "s/^\([\"']\)\(.*\)\1$/\2/g")
echo "Done"
#
# Remove all coverage.out and coverage_review.out files
echo "removing coverage files"
find . -name "coverage*.out" -type f -delete 1>>/dev/null
#
# Display processing message
echo -n 'processing internal tests .'
#
# The following will generate coverage.out files and new packages should be added.
echo -n 'sLogger ' 1>/tmp/tmp_$$.out
go test sLogger/*.go -coverprofile sLogger/coverage.out 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sError ' 1>>/tmp/tmp_$$.out
go test sError/*.go -coverprofile sError/coverage.out 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sDatabase ' 1>>/tmp/tmp_$$.out
rm sDatabase/coverage.out 2>/dev/null
go test sDatabase/*.go -coverprofile sDatabase/coverage.out 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sConfigParams ' 1>>/tmp/tmp_$$.out
rm sConfigParams/coverage.out 2>/dev/null
go test sConfigParams/*.go -coverprofile sConfigParams/coverage.out 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sAuthentication ' 1>>/tmp/tmp_$$.out
rm sAuthentication/coverage.out 2>/dev/null
GOARCH=amd64 go test sAuthentication/*.go -coverprofile sAuthentication/coverage.out 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sMessage ' 1>>/tmp/tmp_$$.out
rm sMessage/coverage.out 2>/dev/null
go test sMessage/*.go -coverprofile sMessage/coverage.out 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sHTTPClient ' 1>>/tmp/tmp_$$.out
rm sHTTPClient/coverage.out 2>/dev/null
go test sHTTPClient/*.go -coverprofile sHTTPClient/coverage.out 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sCustom ' 1>>/tmp/tmp_$$.out
rm sCustom/coverage.out 2>/dev/null
go test sCustom/*.go -coverprofile sCustom/coverage.out 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sDocument ' 1>>/tmp/tmp_$$.out
rm sDocument/coverage.out 2>/dev/null
go test sDocument/*.go -coverprofile sDocument/coverage.out 1>>/tmp/tmp_$$.out
echo "Done"
cat /tmp/tmp_$$.out
read RC <<<"$(grep '^FAIL' /tmp/tmp_$$.out | awk '/[F][A][I][L]/ {print 1}')"
if [[ "$RC" == "1" ]]; then
  echo "FAILED: At least one test didn't pass or the code didn't build."
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
echo -n 'sLogger_external_test ' 1>/tmp/tmp_$$.out
go test tests/sLogger_external_test.go 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sError_external_test ' 1>>/tmp/tmp_$$.out
go test tests/sError_external_test.go 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sDatabase_external_test ' 1>>/tmp/tmp_$$.out
go test tests/sDatabase_external_test.go 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sConfigParams_external_test ' 1>>/tmp/tmp_$$.out
go test tests/sConfigParams_external_test.go 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sAuthentication ' 1>>/tmp/tmp_$$.out
go test tests/sAuthentication_external_test.go 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sMessage_external_test ' 1>>/tmp/tmp_$$.out
go test tests/sMessages_external_test.go 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sHTTPClient_external_test ' 1>>/tmp/tmp_$$.out
go test tests/sHTTPClient_external_test.go 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sCustom_external_test ' 1>>/tmp/tmp_$$.out
go test tests/sCustom_external_test.go 1>>/tmp/tmp_$$.out
echo -n '.'
echo -n 'sDocument_external_test ' 1>>/tmp/tmp_$$.out
go test tests/sDocument_external_test.go 1>>/tmp/tmp_$$.out
echo -n '.'
echo "Done"
cat /tmp/tmp_$$.out
read RC <<<"$(grep '^FAIL' /tmp/tmp_$$.out | awk '/[F][A][I][L]/ {print 1}')"
if [[ "$RC" == "1" ]]; then
  echo "FAILED: At least one test didn't pass or the code didn't build."
  exit 1
else
  echo "PASSED: All external tests passed"
fi
#
# Display processing message
echo -n 'processing test coverage .'
#
# # To output the results from the -coverprofile, using the following command
# shellcheck disable=SC2129
go tool cover -func=sLogger/coverage.out >>coverage_review.out
go tool cover -func=sError/coverage.out >>coverage_review.out
go tool cover -func=sDatabase/coverage.out >>coverage_review.out
go tool cover -func=sConfigParams/coverage.out >>coverage_review.out
go tool cover -func=sAuthentication/coverage.out >>coverage_review.out
go tool cover -func=sMessage/coverage.out >>coverage_review.out
go tool cover -func=sHTTPClient/coverage.out >>coverage_review.out
go tool cover -func=sDocument/coverage.out >>coverage_review.out
go tool cover -func=sCustom/coverage.out >>coverage_review.out
echo "Done"
#
# # Review the coverage totals for 70%+ compliance
# shellcheck disable=SC2162
read RC <<<"$(grep '^total' coverage_review.out | awk '/[0-6][0-9].&&![100]./ {print 1}')"
if [[ "$RC" == "1" ]]; then
  echo "FAILED: At least one components coverage is less than 70%"
  exit 1
else
  echo "PASSED: All components coverage is 70% or higher"
  exit 0
fi
echo "For details, look at coverage_review.out"
