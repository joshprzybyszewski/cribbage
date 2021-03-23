import {App, Duration, Stack, StackProps} from "@aws-cdk/core";
import {DatabaseInstance, DatabaseInstanceEngine, StorageType} from '@aws-cdk/aws-rds';
import {ISecret, Secret} from '@aws-cdk/aws-secretsmanager';
import {InstanceClass, InstanceSize, InstanceType, Peer, SubnetType, Vpc} from "@aws-cdk/aws-ec2";
import { Ec2Service } from "@aws-cdk/aws-ecs";

export interface RDSStackProps extends StackProps {
    vpc: Vpc;
}

export class RDSStack extends Stack {

    readonly secret: ISecret;
    readonly mySQLRDSInstance: DatabaseInstance;
    readonly dsnUser: string;
    readonly dsnPw: string;

    constructor(scope: App, id: string, props: RDSStackProps) {
        super(scope, id, props);

    // Place your resource definitions here
    this.secret = Secret.fromSecretAttributes(this, 'SamplePassword', {
        secretArn: 'arn:aws:secretsmanager:{region}:{organisation-id}:secret:ImportedSecret-sample',
    });


    // TODO
    // what's more, we're going to need to figure out how to create users
    // in it, and then how to get it initialized to have all of the tables
    // created for us.
    this.dsnUser = process.env.SECRET_DSN_USER || "TODO_setDsnUser";
    this.dsnPw = process.env.SECRET_DSN_PASSWORD || "TODO_setDsnPassword";
    
    this.mySQLRDSInstance = new DatabaseInstance(this, 'mysql-rds-instance', {
        engine: DatabaseInstanceEngine.MYSQL,
        instanceType: new InstanceType('t2.micro'),
        // instanceClass: InstanceType.of(InstanceClass.T2, InstanceSize.SMALL),
        credentials: {
            username: this.dsnUser.toString(),
            // TODO figure out the best way to get a secret here
            // password: dsnPw.toString(),
            // password: '',
        },
        vpc: props.vpc,
        vpcSubnets: {
            subnetType: SubnetType.ISOLATED,
        },
        storageEncrypted: true,
        multiAz: false,
        autoMinorVersionUpgrade: false,
        allocatedStorage: 25, // TODO what is this?
        storageType: StorageType.GP2, // TODO what is this?
        backupRetention: Duration.days(3), // TODO what is the cheapest?
        deletionProtection: false,
        databaseName: 'cribbage', // TODO we could/should relate this to CRIBBAGE_MYSQL_DB=cribbage
        port: 3306
    });


    
    }
}