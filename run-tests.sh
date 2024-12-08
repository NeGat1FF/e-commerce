#!/bin/bash
for dir in */; do
  if [ -f "$dir/go.mod" ]; then
    echo "Running tests in $dir"
    (cd "$dir" && go test ./... -v)
  fi
done
