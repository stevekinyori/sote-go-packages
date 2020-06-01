#!/usr/bin/env bash
#
# This is a testing script for the seeding utility
#
#go test Tests/*
#
# Test run.go and lookupValue.go
go test -coverprofile packages/coverage.out
#
# To output the results from the -coverprofile, using the following command
#go tool cover -func=packages/coverage.out