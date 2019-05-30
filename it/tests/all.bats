#!/usr/bin/env bats

# All tests

load test_helper

@test "Run without arguments" {
	run it
	assert_success
}
