#!/usr/bin/env bash

set -eu

CURRENT_TAG=$(git describe --abbrev=0 --tags || echo 'v0.0.0')
CURRENT_TAG_NO_V=${CURRENT_TAG#"v"}

PATCH=$(echo $CURRENT_TAG_NO_V | cut -d'.' -f3)
MINOR=$(echo $CURRENT_TAG_NO_V | cut -d'.' -f2)
MAJOR=$(echo $CURRENT_TAG_NO_V | cut -d'.' -f1)

BUMP=0
while IFS= read -r line; do
  if [[ $line == *BREAKING\ CHANGE* ]]; then
    if [[ $BUMP -lt 3 ]]; then
      BUMP=3

      PATCH=0
      MINOR=0
      MAJOR=$((MAJOR+1))
    fi
  elif [[ $line == *feat:* || $line == *feat\(*\):* ]]; then
    if [[ $BUMP -lt 2 ]]; then
      BUMP=2

      PATCH=0
      MINOR=$((MINOR+1))
    fi
  elif [[ $line != "" ]]; then
    if [[ $BUMP -lt 1 ]]; then
      BUMP=1

      PATCH=$((PATCH+1))
    fi
  fi
done <<< "$(git log --oneline $(git describe --tags --abbrev=0 @^)..@)"

if [[ $BUMP == 0 ]]; then
  echo "There is nothing new to release"
  exit 1
fi

echo "v$MAJOR.$MINOR.$PATCH"
