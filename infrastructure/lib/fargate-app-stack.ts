import * as ec2 from "@aws-cdk/aws-ec2";
import * as ecs_patterns from "@aws-cdk/aws-ecs-patterns";
import * as cdk from "@aws-cdk/core";
import * as ecs from "@aws-cdk/aws-ecs";
import { Secret } from '@aws-cdk/aws-secretsmanager';

export interface FargateAppStackProps extends cdk.StackProps {
  vpc: ec2.Vpc
  dbSecretArn: string | undefined
}


export class FargateAppStack extends cdk.Stack {

  readonly albFargateService: ecs_patterns.ApplicationLoadBalancedFargateService;

  constructor(scope: cdk.Construct, id: string, props: FargateAppStackProps) {
    super(scope, id, props);

    const dbUserSecret = Secret.fromSecretCompleteArn(this,
      `cribbageRDSSecret`,
      props.dbSecretArn ?? 'badarn',
    )

    var dsnHost = dbUserSecret.secretValueFromJson('dbname').toString();
    var dbType = dsnHost.length > 0 ? "mysql" : "memory";

    const serverPort = 80;

    const vpc = props.vpc;

    const cluster = new ecs.Cluster(this,
      `${id}-cluster`,
      {
        vpc: vpc
      },
    );

    this.albFargateService = new ecs_patterns.ApplicationLoadBalancedFargateService(this,
      `${id}-fargate`,
      {
        cluster: cluster, // Required
        memoryLimitMiB: 1024, // Default is 512
        cpu: 512, // Default is 256
        desiredCount: 1, // Default is 1
        publicLoadBalancer: true, // Default is false
        taskImageOptions: {
          // here's how we'd grab the image from dockerhub:
          image: ecs.ContainerImage.fromRegistry(
            "joshprzybyszewski/cribbage:pr-94-merge" // TODO change this to use `master` or `latest-tag`, not `latest`.
          ),
          containerPort: serverPort,
          environment: {
            CRIBBAGE_RESTPORT: serverPort.toString(),
            CRIBBAGE_DB: dbType,
            CRIBBAGE_DSN_HOST: dbUserSecret.secretValueFromJson('host').toString(),
            CRIBBAGE_DSN_USER: dbUserSecret.secretValueFromJson('username').toString(),
            CRIBBAGE_DSN_PASSWORD: dbUserSecret.secretValueFromJson('password').toString(),
            CRIBBAGE_MYSQL_DB: dbUserSecret.secretValueFromJson('dbname').toString(),
            CRIBBAGE_MYSQL_CREATE_TABLES: 'true', // mysql_create_tables may or may not be a good thing...
            CRIBBAGE_DSN_PARAMS: "parseTime=true&timeout=90s&writeTimeout=90s&readTimeout=90s&tls=skip-verify&maxAllowedPacket=1000000000&rejectReadOnly=true",
            deploy: 'prod',
          },
        },
      },
    );
  }
}
