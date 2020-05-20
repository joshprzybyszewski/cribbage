#!/bin/bash

# Check that you have a version of go we know works
GO_VERSION=$(go version)
if [[ $GO_VERSION != *"go1.14.2"* ]]; then
  echo "WARNING: We recommend using go version 1.14.2, but you have \"${GO_VERSION}\""
fi

# Check that you have mongodb and that it's a version known to work
if [[ ! $(which mongo) ]]; then
    echo "You don't have mongodb installed"
    echo "You can download it here: https://www.mongodb.com/download-center/community"
    echo "And follow instructions here: https://docs.mongodb.com/manual/installation/"
    echo "OR use brew: "
    echo "  brew tap mongodb/brew"
    echo "  brew install mongodb-community@4.2"
    echo "if you're on macOS"
else
    # Check that you have a version of mongo we know works
    MONGO_VERSION=$(mongo --version)
    if [[ $MONGO_VERSION != *"v4.2.1"* ]]; then
        echo "WARNING: We recommend using MONGO at version v4.2.1, but you have \"${MONGO_VERSION}\""
    fi
    # Check that mongodb is running.
    if [[ ! $(pgrep mongo) ]]; then
        echo "Did not find MongoDB running. Try 'mongo' to start it."
    fi
fi

echo ""
echo "Also for mongo, check out replica set install instructions here: http://thecodebarbarian.com/introducing-run-rs-zero-config-mongodb-runner"
echo "For mongo replicaset, it may simply be:"
echo "  npm install run-rs -g"
echo "  run-rs -v 4.2.1 --shell"
echo "Or it may be more difficult. Good luck."
echo ""
