#!/usr/bin/env bash

load "$BATS_ROOT"/lib/bats/assert.bash

# ------------------------------------------------------------------------------
# Setup
# ------------------------------------------------------------------------------

export PATH=$PWD:$PATH

setup() {
	[[ ! -f ${BATS_PARENT_TMPNAME}.skip ]] || skip "previous test failed! skipping remaining tests..."

	fixture_enter
	file_create t
	git_create
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
	local arg
	for arg; do
		echo -e "\\e[1;38;5;11m$arg\\e[0m" >&3
	done
}

die() {
	local arg
	for arg; do
		echo -e "\\e[1;38;5;198m$arg\\e[0m" >&3
	done

	exit 1
}

see() {
	"$@" >&3 2>&3
}

join() {
	local separator=$1
	shift
	echo "$(IFS="$separator"; echo "${*}")"
}

cr() {
	echo -en "$1\\r"
}

file_create() {
	local file=$1
	shift

	mkdir -p "$(dirname "$file")"

	[[ -t 0 ]] || cat - >"$file"
	echo "$*" >>"$file"
}

declare -ag cleanup=()

fixture_enter() {
	local tmpdir

	tmpdir=$(mktemp -d "${BATS_TMPDIR:-/tmp}/it.XXXXX")
	cleanup+=("$tmpdir")

	pushd "$tmpdir" &>/dev/null || exit 1
	mkdir -p fixture
	pushd fixture &>/dev/null || exit 1
}

git_create() {
	git init
	git config user.email "robot@example.com"
	git config user.name "Test Robot"
	git add .
	git commit -m "initial commit"
}
