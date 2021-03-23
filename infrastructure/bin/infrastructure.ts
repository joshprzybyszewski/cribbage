#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
import { VpcStack } from "../lib/vpc-stack";
import { FargateAppStack } from "../lib/fargate-app-stack";
import { RDSStack } from "../lib/rds-stack";

const app = new cdk.App();
const vpcStackEntity  = new VpcStack(app, 'cribbage-vpc');
const rdsStack = new RDSStack(app, 'cribbage-rds', {
    vpc: vpcStackEntity.vpc
});
new FargateAppStack(app, 'cribbage-app', {
    vpc: vpcStackEntity.vpc,
    // TODO
    rdsEndpoint: rdsStack.mySQLRDSInstance.dbInstanceEndpointAddress,
    rdsDbUser: rdsStack.dsnUser,
    rdsDbPassword: rdsStack.dsnPw,
    // subnetName: vpcStackEntity._subnetName
});
