import * as ecs_patterns from '@aws-cdk/aws-ecs-patterns';
import * as cdk from '@aws-cdk/core';
import { Code, Function, Runtime } from '@aws-cdk/aws-lambda';
import { Table } from '@aws-cdk/aws-dynamodb';
import * as assets from '@aws-cdk/aws-s3-assets';

export interface LambdaStackProps extends cdk.StackProps {
    readonly table: Table;
}

export class LambdaStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props: LambdaStackProps) {
        super(scope, id, props);

        const zipFileName = 'cribbage-lambda';

        const asset = new assets.Asset(this, 'CribbageZip', {
            path: '../' + zipFileName + '.zip',
        });

        const lambda = new Function(this, 'serverlessLambdaHandler', {
            runtime: Runtime.GO_1_X,
            code: Code.fromBucket(asset.bucket, asset.s3ObjectKey),
            handler: zipFileName,
            timeout: cdk.Duration.seconds(15),
            environment: {
                'CRIBBAGE_DB': 'dynamodb',
                'CRIBBAGE_LAMBDA': 'true',
                'TABLE_NAME': props.table.tableName,
            },
        })


        props.table.grantReadWriteData(lambda);
        // DescribeTable is not granted by default, but this should resolve it according to: https://github.com/aws/aws-cdk/issues/7633#issuecomment-621672844
        props.table.grant(lambda, 'dynamodb:DescribeTable');
    }
}
