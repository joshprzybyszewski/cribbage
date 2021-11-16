#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { DynamoDBStack } from '../lib/dynamodb-stack';
import { LambdaStack } from '../lib/lambda-stack';
import { DNSStack } from '../lib/dns-stack';
import { Environment } from '@aws-cdk/core';

const env: Environment = {
    account: process.env['AWS_ACCOUNT'],
    region: process.env['AWS_REGION'],
};

const app = new cdk.App();

console.log('building dynamodb table...');
const dynamoStack = new DynamoDBStack(app, 'cribbage-dynamodb', {
    env,
});
console.log('building dynamodb table...Done!');

console.log('building lambdas...');
const lambdaStack = new LambdaStack(app, 'cribbage-lambda', {
    table: dynamoStack.table,
    env,
});
console.log('building lambdas...Done!');

// TODO createa API gateway that talks through to the lambda

// TODO create a route 53 entry to pass hobbycribbage.com to the api gateway?

// TODO re-architect the frontend so that the "backend" calls go to a different sub-domain to hit
// the lambda, and all of hobbycribbage.com/ just goes to a CDN to get the page. It's a SPA, so
// I don't want to overthink this.
