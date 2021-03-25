#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { VpcStack } from '../lib/vpc-stack';
import { FargateAppStack } from '../lib/fargate-app-stack';
import { RDSStack } from '../lib/rds-stack';
import { Port } from '@aws-cdk/aws-ec2';

const app = new cdk.App();

console.log('building vpc...');
const vpcStackEntity = new VpcStack(app, 'cribbage-vpc');
console.log('building vpc...Done!');

console.log('building rds...');
const rdsStack = new RDSStack(app, 'cribbage-rds', {
    vpc: vpcStackEntity.vpc,
    rdsIngressPort: vpcStackEntity.rdsPort,
});
console.log('building rds...Done!');

console.log('building fargate app...');
const fargateStack = new FargateAppStack(app, 'cribbage-app', {
    vpc: vpcStackEntity.vpc,
    dbSecretArn: rdsStack.mySQLRDSInstance.secret?.secretArn,
});
console.log('building fargate app...Done!');

// I read the following doc which recommended peering two constructs using connections in python: https://docs.aws.amazon.com/cdk/api/latest/python/aws_cdk.aws_ec2/SecurityGroup.html
console.log('Allowing connections between constructs...');
fargateStack.albFargateService.service.connections.allowTo(rdsStack.mySQLRDSInstance, Port.tcp(vpcStackEntity.rdsPort));
console.log('Allowing connections between constructs...Done!');
