#!/usr/bin/env bash

. `git rev-parse --show-toplevel`/scripts/common.sh

OLD_TAG=`git describe --tags --abbrev=0`
CURRENT_TAG=`git describe --tags`

echo -n "Last explicit tag: "
yellow $OLD_TAG
echo -n "Current tag: "
yellow $CURRENT_TAG

echo -n "Enter the new tag: "
read NEW_TAG

echo -n "New tag: "
green $NEW_TAG

echo -n "Does this look OK? [Press 'y'] "
read -n 1 CONFIRM
echo

if [ "$CONFIRM" != "y" ]; then
  red "Stopping"
  exit 1
fi
echo

echo -n "Checking for dirty git tree... "
if test -n "$(git status --porcelain)"; then
  red "Uncommitted changes!"
  exit 1
fi
green "OK"

echo -n "Tagging release... "
git tag --annotate $NEW_TAG -m "Version $NEW_TAG ($NEW_VER)"
abort_if_error
green "OK"

echo -n "Creating release dir if it doesn't exist... "
mkdir -p release
abort_if_error
green "OK"

echo -n "Outputting release notes... "
git log --abbrev-commit --pretty=oneline --reverse \
  $OLD_TAG..HEAD > release/notes_${NEW_TAG}.txt
abort_if_error
green "OK"

echo -n "Comitting changes... "
git add release
git commit -m "New release - $NEW_VER"
abort_if_error
green "OK"

echo -n "Updating tag... "
git tag --annotate --force $NEW_TAG -m "Version $NEW_TAG ($NEW_VER)"
abort_if_error
green "OK"

echo -n "Pushing to origin... "
git push origin
abort_if_error
git push origin --tags
abort_if_error
green "OK"
