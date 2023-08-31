#!/usr/bin/bash

# checkout https://github.com/programme-lv/example-tasks into tmp dir
# copy all directories from the tmp dir into the upload dir
# use absolute location ../upload

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
UPLOAD_DIR=$SCRIPT_DIR/../upload
REPO_URL="https://github.com/programme-lv/example-tasks"

# make sure upload dir exists
if [ ! -d "$UPLOAD_DIR" ]; then
	mkdir "$UPLOAD_DIR"
fi

TMP_DIR=$(mktemp -d -t git_clone_XXXXXX)
git clone "$REPO_URL" "$TMP_DIR"
cd "$TMP_DIR"

for dir in $TMP_DIR/*; do
	if [ -d "$dir" ]; then
		if [ -d "$UPLOAD_DIR/$(basename "$dir")" ]; then
			rm -rf "$UPLOAD_DIR/$(basename "$dir")"
		fi
		cp -r "$dir" "$UPLOAD_DIR"
	fi
done
