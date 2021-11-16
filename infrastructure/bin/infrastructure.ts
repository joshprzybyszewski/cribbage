#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { DynamoDBStack } from '../lib/dynamodb-stack';
import { LambdaStack } from '../lib/lambda-stack';
import { DNSStack } from '../lib/dns-stack';
import { Environment } from '@aws-cdk/core';

const env: Environment = {
    account: process.env['AWS_ACCOUNT'] || process.env['CDK_DEFAULT_ACCOUNT'],
    region: process.env['AWS_REGION'] || process.env['CDK_DEFAULT_REGION'],
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

console.log('Creating a DNS stack...');

const dnsStack = new DNSStack(app, 'cribbage-dns', {
    lambda: lambdaStack.lambda,
    env,
});
console.log('Creating a DNS stack...Done!');
