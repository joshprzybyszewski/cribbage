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
console.log('you might need to bootstrap');
// cdk bootstrap aws://ACCOUNT-NUMBER/REGION
const lambdaStack = new LambdaStack(app, 'cribbage-lambda', {
    table: dynamoStack.table,
    env,
});
console.log('building lambdas...Done!');

console.log('Granting RW access to dynamo from lambdas...');
console.log('this does not seem to be enough');
/*
2021/11/15 02:39:33 DescribeTable ERROR: operation error DynamoDB: 
DescribeTable, https response error StatusCode: 400, 
RequestID: 5Q1B31KL6H7KIMUTDPEHCJM1EVVV4KQNSO5AEMVJF66Q9ASUAAJG, 
api error AccessDeniedException: 
User: arn:aws:sts::971042860856:assumed-role/cribbage-lambda-cribbagelambdaidServiceRole97F1334-1JABCY835OEX/cribbage-lambda-cribbagelambdaid0ADEAD21-uOm3ajM1vXir 
 is not authorized to perform: dynamodb:DescribeTable on 
 resource: arn:aws:dynamodb:us-east-2:971042860856:table/cribbage

*/
dynamoStack.table.grantReadWriteData(lambdaStack.function);
console.log('Granting RW access to dynamo from lambdas...Done!');

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
