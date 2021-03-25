import { App, Stack, StackProps } from "@aws-cdk/core";
import { DatabaseInstance, DatabaseInstanceEngine, StorageType, MysqlEngineVersion, Credentials } from '@aws-cdk/aws-rds';
import { InstanceType, Vpc } from "@aws-cdk/aws-ec2";

export interface RDSStackProps extends StackProps {
    vpc: Vpc;
    rdsIngressPort: number;
}

export class RDSStack extends Stack {

    readonly mySQLRDSInstance: DatabaseInstance;
    readonly creds: Credentials;
    readonly dbName: string;

    constructor(scope: App, id: string, props: RDSStackProps) {
        super(scope, id, props);

        this.dbName = 'cribbage';
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
        });



    }
}