import * as ec2 from "@aws-cdk/aws-ec2";
import * as ecs_patterns from "@aws-cdk/aws-ecs-patterns";
import * as cdk from "@aws-cdk/core";
import * as ecs from "@aws-cdk/aws-ecs";

export interface FargateAppStackProps extends cdk.StackProps {
  vpc: ec2.Vpc
  rdsEndpoint: string
  rdsDbUser: string
  dsnPassword: cdk.SecretValue | undefined
  dbName: string
}


export class FargateAppStack extends cdk.Stack {

  constructor(scope: cdk.Construct, id: string, props: FargateAppStackProps) {
    super(scope, id, props);

    var dsnHost = props.rdsEndpoint;
    var dbType = dsnHost.length > 0 ? "mysql" : "memory";

    const serverPort = 80;

    const vpc = props.vpc;

    const cluster = new ecs.Cluster(this, `${id}-cluster`, {
      vpc: vpc
    });

    new ecs_patterns.ApplicationLoadBalancedFargateService(this, `${id}-fargate`, {
      cluster: cluster, // Required
      memoryLimitMiB: 1024, // Default is 512
      cpu: 512, // Default is 256
      desiredCount: 1, // Default is 1
      publicLoadBalancer: true, // Default is false
      taskImageOptions: {
        // here's how we'd grab the image from dockerhub:
        image: ecs.ContainerImage.fromRegistry(
          "joshprzybyszewski/cribbage:latest" // TODO change this to use `master`, not `latest`.
        ),
        containerPort: serverPort,
        environment: {
          CRIBBAGE_DB: dbType,
          CRIBBAGE_RESTPORT: serverPort.toString(),
          CRIBBAGE_DSN_PARAMS: "parseTime=true&timeout=90s&writeTimeout=90s&readTimeout=90s&tls=skip-verify&maxAllowedPacket=1000000000&rejectReadOnly=true",
          CRIBBAGE_DSN_HOST: dsnHost?.toString() ?? "",
          CRIBBAGE_DSN_USER: props.rdsDbUser.toString(),
          CRIBBAGE_DSN_PASSWORD: props.dsnPassword?.toString() || 'bigbigfail',
          CRIBBAGE_MYSQL_DB: props.dbName,
          CRIBBAGE_MYSQL_CREATE_TABLES: 'true', // mysql_create_tables may or may not be a good thing...
          deploy: 'prod',
        },
      },
    });
  }
}
