import { App, Stack, StackProps } from '@aws-cdk/core';
import { AttributeType, BillingMode, Table } from '@aws-cdk/aws-dynamodb';

export interface DynamoDBStackProps extends StackProps {
}

export class DynamoDBStack extends Stack {
    readonly table: Table;

    constructor(scope: App, id: string, props: DynamoDBStackProps) {
        super(scope, id, props);

        this.table = new Table(this, 'cribbage', {
          partitionKey: { name: 'DDBid', type: AttributeType.STRING },
          sortKey: { name: 'spec', type: AttributeType.STRING },
          billingMode: BillingMode.PAY_PER_REQUEST,
        // --region us-west-2 
        });
    }
}
