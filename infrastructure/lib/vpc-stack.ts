import { App, Stack, StackProps } from '@aws-cdk/core';
import { Peer, Port, SecurityGroup, Subnet, SubnetType, Vpc } from '@aws-cdk/aws-ec2'

export class VpcStack extends Stack {
    readonly vpc: Vpc;
    // readonly ingressSecurityGroup: SecurityGroup;
    // readonly egressSecurityGroup: SecurityGroup;

    constructor(scope: App, id: string, props?: StackProps) {
        super(scope, id, props);

        this.vpc = new Vpc(this, `${id}-vpc`, {
            maxAzs: 3, // Default is all AZs in region
            subnetConfiguration: [
                {
                    cidrMask: 24,
                    name: 'ingress',
                    subnetType: SubnetType.PUBLIC,
                },
                {
                    cidrMask: 28,
                    name: 'rds',
                    subnetType: SubnetType.PRIVATE,
                },
            ]
        });
    }
}