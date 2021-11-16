import * as ecs_patterns from '@aws-cdk/aws-ecs-patterns';
import * as cdk from '@aws-cdk/core';
import {Code, Function, Runtime}  from '@aws-cdk/aws-lambda';
import {  Table } from '@aws-cdk/aws-dynamodb';
import * as assets from '@aws-cdk/aws-s3-assets';

export interface LambdaStackProps extends cdk.StackProps {
    readonly table: Table;
}

export class LambdaStack extends cdk.Stack {
    readonly albFargateService: ecs_patterns.ApplicationLoadBalancedFargateService;
    readonly function: Function;

    constructor(scope: cdk.Construct, id: string, props: LambdaStackProps) {
        super(scope, id, props);

        const asset = new assets.Asset(this, 'CribbageZip', {
            path: '../cribbage-lambda.zip',
          });

        this.function = new Function(this, 'cribbage-lambda-id', {
            runtime: Runtime.GO_1_X,
            code: Code.fromBucket(asset.bucket, asset.s3ObjectKey),
            handler: 'cribbage-lambda',
            timeout: cdk.Duration.seconds(15),
            environment: {
                'CRIBBAGE_DB': 'dynamodb',
                'CRIBBAGE_LAMBDA': 'true',
                'TABLE_NAME': props.table.tableName,
            },
        })

    }
}