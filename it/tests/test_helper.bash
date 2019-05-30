#!/usr/bin/env bash

load lib

# ------------------------------------------------------------------------------
# Setup
# ------------------------------------------------------------------------------

export PATH=$PWD:$PATH

setup() {
	[[ ! -f ${BATS_PARENT_TMPNAME}.skip ]] || skip "previous test failed! skipping remaining tests..."

	fixture.enter
	file.create t ""
	git.create
}

teardown() {
	if [[ -z $BATS_TEST_COMPLETED ]]; then
		touch "${BATS_PARENT_TMPNAME}.skip"
		[[ ${#cleanup[@]} -eq 0 ]] || cry "Did not remove $(join ' ' "${cleanup[@]}") as test failed"
	else
		rm -rf -- "${cleanup[@]}"
	fi
}

# ------------------------------------------------------------------------------
# Helpers
# ------------------------------------------------------------------------------

cry() {
	printf "\\x1B[K" >/dev/tty

	local arg
	for arg; do
		echo -e "\\e[1;38;5;11m$arg\\e[0m" >/dev/tty
	done
}

die() {
	printf "\\x1B[K" >/dev/tty

	local arg
	for arg; do
		echo -e "\\e[1;38;5;198m$arg\\e[0m" >/dev/tty
	done

	exit 1
}

join() {
	local separator=$1
	shift
	echo "$(IFS="$separator"; echo "${*}")"
}

cr() {
	echo -en "$1\\r"
}

file.create() {
	local file=$1
	shift

	mkdir -p "$(dirname "$file")"

	if [[ $# -eq 0 ]]; then
		cat -
	else
		echo "$*"
	fi >"$file"
}

declare -ag cleanup=()

fixture.enter() {
	local tmpdir

	tmpdir=$(mktemp -d "${BATS_TMPDIR:-/tmp}/it.XXXXX")
	cleanup+=("$tmpdir")

	pushd "$tmpdir" &>/dev/null || exit 1
	mkdir -p fixture
	pushd fixture &>/dev/null || exit 1
}

git.create() {
	git init
	git config user.email "robot@example.com"
	git config user.name "Test Robot"
	git add .
	git commit -m "initial commit"
}
