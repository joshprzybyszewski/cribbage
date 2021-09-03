#!/bin/bash
set -e
set -x

aws dynamodb create-table \
    --endpoint-url http://dynamodb-local:8000 \
    --region us-west-2 \
    --billing-mode PAY_PER_REQUEST \
    --table-name cribbage \
    --attribute-definitions \
        AttributeName=DDBid,AttributeType=S \
        AttributeName=spec,AttributeType=S \
    --key-schema \
        AttributeName=DDBid,KeyType=HASH \
        AttributeName=spec,KeyType=Range