import { App, Duration, Stack, StackProps } from "@aws-cdk/core";
import { DatabaseInstance, DatabaseInstanceEngine, StorageType, MysqlEngineVersion, Credentials } from '@aws-cdk/aws-rds';
import { ISecret, Secret } from '@aws-cdk/aws-secretsmanager';
import { InstanceClass, InstanceSize, InstanceType, Peer, SubnetType, Vpc } from "@aws-cdk/aws-ec2";
import { Ec2Service } from "@aws-cdk/aws-ecs";

export interface RDSStackProps extends StackProps {
    vpc: Vpc;
}

export class RDSStack extends Stack {

    readonly mySQLRDSInstance: DatabaseInstance;
    readonly creds: Credentials;
    readonly dbName: string;

    constructor(scope: App, id: string, props: RDSStackProps) {
        super(scope, id, props);

        this.dbName = 'cribbage';
        // this.dsnPassword = process.env.SECRET_CDK_DSN_PW || 'bigfailure';
        this.creds = Credentials.fromGeneratedSecret(
            'cribbageApp',
            {},
       );

        const mysqlFullVersion = '8.0.19'; // matches the version in CI and install.sh
        const mysqlMajorVersion = '8.0';
        this.mySQLRDSInstance = new DatabaseInstance(this, 'mysql-rds-instance', {
            engine: DatabaseInstanceEngine.mysql({
                version: MysqlEngineVersion.of(mysqlFullVersion, mysqlMajorVersion),
            }),
            instanceType: new InstanceType('db.t2.micro'),
            credentials: this.creds,
            vpc: props.vpc,
            vpcSubnets: {
                subnetType: SubnetType.ISOLATED,
            },
            storageEncrypted: true,
            multiAz: false,
            autoMinorVersionUpgrade: false,
            allocatedStorage: 25, // number of GB. Default is 100
            deletionProtection: true,
            databaseName: this.dbName,
            port: 3306
        });



    }
}