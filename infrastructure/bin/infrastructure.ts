#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { VpcStack } from "../lib/vpc-stack";
import { FargateAppStack } from "../lib/fargate-app-stack";
import { RDSStack } from "../lib/rds-stack";
import { exit } from "process";

const app = new cdk.App();

console.log('building vpc...');
const vpcStackEntity = new VpcStack(app, 'cribbage-vpc');
console.log('building vpc...Done!');

console.log('building rds...');
const rdsStack = new RDSStack(app, 'cribbage-rds', {
    vpc: vpcStackEntity.vpc,
});
console.log('building rds...Done!');

console.log('building fargate app...');
new FargateAppStack(app, 'cribbage-app', {
    vpc: vpcStackEntity.vpc,
    rdsEndpoint: rdsStack.mySQLRDSInstance.dbInstanceEndpointAddress,
    rdsDbUser: rdsStack.creds.username,
    dsnPassword: rdsStack.creds.password,
    dbName: rdsStack.dbName,
    // subnetName: vpcStackEntity._subnetName
});
console.log('building fargate app...Done!');
