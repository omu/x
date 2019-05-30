#!/usr/bin/env bats

# All tests

load test_helper

@test "Run without arguments" {
	run true
	file_create foo
	see ls -al
	assert_success
}
