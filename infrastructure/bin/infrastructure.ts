#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { VpcStack } from '../lib/vpc-stack';
import { DynamoDBStack } from '../lib/dynamodb-stack';
import { LambdaStack } from '../lib/lambda-stack';
// import { FargateAppStack } from '../lib/fargate-app-stack';
// import { RDSStack } from '../lib/rds-stack';
import { Port } from '@aws-cdk/aws-ec2';
import { DNSStack } from '../lib/dns-stack';
import { Environment } from '@aws-cdk/core';

const env: Environment = {
    account: process.env['AWS_ACCOUNT'],
    region: process.env['AWS_REGION'],
};

const app = new cdk.App();

// console.log('building vpc...');
// const vpcStackEntity = new VpcStack(app, 'cribbage-vpc', { env });
// console.log('building vpc...Done!');

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

// console.log('building rds...');
// const rdsStack = new RDSStack(app, 'cribbage-rds', {
//     vpc: vpcStackEntity.vpc,
//     rdsIngressPort: vpcStackEntity.rdsPort,
//     env,
// });
// console.log('building rds...Done!');

// console.log('building fargate app...');
// const fargateStack = new FargateAppStack(app, 'cribbage-app', {
//     vpc: vpcStackEntity.vpc,
//     dbSecretArn: rdsStack.mySQLRDSInstance.secret?.secretArn,
//     env,
// });
// console.log('building fargate app...Done!');

// console.log('building dns...');
// new DNSStack(app, 'cribbage-dns', {
//     loadBalancer: fargateStack.albFargateService.loadBalancer,
//     env,
// });
// console.log('building dns...Done!');

// I read the following doc which recommended peering two constructs 
// using connections in python: https://docs.aws.amazon.com/cdk/api/latest/python/aws_cdk.aws_ec2/SecurityGroup.html
// console.log('Allowing connections between constructs...');
// fargateStack.albFargateService.service.connections.allowTo(rdsStack.mySQLRDSInstance, Port.tcp(vpcStackEntity.rdsPort));
// console.log('Allowing connections between constructs...Done!');
