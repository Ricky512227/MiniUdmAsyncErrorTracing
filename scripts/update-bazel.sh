#!/bin/bash
# Script to update Bazel BUILD files using Gazelle

set -e

echo "Updating Bazel BUILD files with Gazelle..."

# Run Gazelle to update BUILD files
bazel run //:gazelle

echo "Done! BUILD files have been updated."

