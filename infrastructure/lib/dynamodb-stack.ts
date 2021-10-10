import { App, Stack, StackProps } from '@aws-cdk/core';
import {
    Credentials,
    DatabaseInstance,
    DatabaseInstanceEngine,
    MysqlEngineVersion,
    OptionGroup,
    StorageType,
} from '@aws-cdk/aws-rds';
import { InstanceType, ISecurityGroup, Peer, Port, SecurityGroup, Vpc } from '@aws-cdk/aws-ec2';

export interface DynamoDBStackProps extends StackProps {
    vpc: Vpc;
    rdsIngressPort: number;
}

export class DynamoDBStack extends Stack {
    readonly mySQLRDSInstance: DatabaseInstance;
    readonly creds: Credentials;
    readonly dbName: string;

    constructor(scope: App, id: string, props: DynamoDBStackProps) {
        super(scope, id, props);

        this.dbName = 'cribbage';
        this.creds = Credentials.fromGeneratedSecret('cribbageApp', {});

        const mysqlFullVersion = '8.0.19'; // matches the version in CI and install.sh
        const mysqlMajorVersion = '8.0';

        // docs for security group loosely based on: https://github.com/aws/aws-cdk/blob/77a6268d696dc0f33fbce4c973f45df29da7aef5/packages/%40aws-cdk/aws-ec2/README.md#allowing-connections
        // We create this security group because otherwise our RDS has egress capability to the whole world.
        const securityGroup: ISecurityGroup = new SecurityGroup(this, 'SecurityGroup', {
            vpc: props.vpc,
            description: 'Allow mysql access within the VPC',
            allowAllOutbound: false, // don't let the DB talk to the world.
        });
        securityGroup.addIngressRule(
            Peer.ipv4('10.0.0.0/24'), // defined in vpc-stack as the cidrMask for our "public" ingress subnet
            Port.tcp(props.rdsIngressPort),
            'Allow containers in the VPC to talk to the DB',
        );
        securityGroup.addEgressRule(
            Peer.ipv4('10.0.0.0/24'), // defined in vpc-stack as the cidrMask for our "public" ingress subnet
            Port.tcp(props.rdsIngressPort),
            'Allow the DB to talk to containers in the VPC',
        );

        this.mySQLRDSInstance = new DatabaseInstance(this, 'mysql-rds-instance', {
            engine: DatabaseInstanceEngine.mysql({
                version: MysqlEngineVersion.of(mysqlFullVersion, mysqlMajorVersion),
            }),
            instanceType: new InstanceType('t2.micro'),
            credentials: this.creds,
            vpc: props.vpc,
            storageEncrypted: false, // not supported by db.t2.micro, otherwise it'd be true
            multiAz: false,
            autoMinorVersionUpgrade: false,
            allocatedStorage: 20, // number of GB. Default is 100
            deletionProtection: true,
            databaseName: this.dbName,
            instanceIdentifier: 'cribbage-storage-mysql',
            port: props.rdsIngressPort,
            securityGroups: [securityGroup],
        });
    }
}
