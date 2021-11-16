import { App, Stack, StackProps } from '@aws-cdk/core';
import { SubnetType, Vpc } from '@aws-cdk/aws-ec2';

export class VpcStack extends Stack {
    readonly vpc: Vpc;
    readonly rdsPort: number;

    constructor(scope: App, id: string, props?: StackProps) {
        super(scope, id, props);

        // Check out this documentation: https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Subnets.html#vpc-subnet-basics
        // But that's confusing, so just look at the docs on this struct.
        this.vpc = new Vpc(this, `${id}-vpc`, {
            maxAzs: 3, // Default is all AZs in region
            subnetConfiguration: [
                {
                    cidrMask: 24,
                    name: 'publicIngressSubnet',
                    subnetType: SubnetType.PUBLIC,
                },
                {
                    cidrMask: 28,
                    name: 'rds',
                    subnetType: SubnetType.PRIVATE,
                },
            ],
        });

        this.rdsPort = 3306;
    }
}
