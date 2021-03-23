import * as ec2 from "@aws-cdk/aws-ec2";
import * as ecs_patterns from "@aws-cdk/aws-ecs-patterns";
import * as cdk from "@aws-cdk/core";
import * as ecs from "@aws-cdk/aws-ecs";

export class InfrastructureStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // TODO we need to figure out best practices for standing up the DB
    // either in the same CF stack, or in a different one.
    // what's more, we're going to need to figure out how to create users
    // in it, and then how to get it initialized to have all of the tables
    // created for us.
    // we may be able to use this dynamoDB creation example as a starter: https://github.com/aws-samples/aws-cdk-examples/blob/18f3429414cc80223a04a04d0249b12a7ae13cbb/typescript/api-cors-lambda-crud-dynamodb/index.ts#L10-L21
    var dsnUser = process.env.SECRET_DSN_USER || "TODO_setDsnUser";
    var dsnPw = process.env.SECRET_DSN_PASSWORD || "TODO_setDsnPassword";
    var dsnHost = process.env.SECRET_DSN_HOST || ""; // TODO figure out how to ref the RDS
    var dbType = dsnHost.length > 0 ? "mysql" : "memory";

    const serverPort = 80;

    const vpc = new ec2.Vpc(this, `${id}-vpc`, {
      maxAzs: 3 // Default is all AZs in region
    });

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
          CRIBBAGE_DSN_USER: dsnUser.toString(),
          CRIBBAGE_DSN_PASSWORD: dsnPw.toString(),
        },
      },
    });
  }
}
