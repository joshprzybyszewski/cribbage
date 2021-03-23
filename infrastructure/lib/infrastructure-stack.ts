import * as ec2 from "@aws-cdk/aws-ec2";
import * as ecs_patterns from "@aws-cdk/aws-ecs-patterns";
import * as cdk from "@aws-cdk/core";
import * as ecs from "@aws-cdk/aws-ecs";

export class InfrastructureStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

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
          CRIBBAGE_DB: "memory", // TODO swap this to "mysql"
          CRIBBAGE_RESTPORT: serverPort.toString(),
          CRIBBAGE_DSN_HOST: "TODO_get_the_URL_to_the_rds",
          CRIBBAGE_DSN_PARAMS: "parseTime=true&timeout=90s&writeTimeout=90s&readTimeout=90s&tls=skip-verify&maxAllowedPacket=1000000000&rejectReadOnly=true",
          CRIBBAGE_DSN_USER: "TODO_figure_out_my_user",
          CRIBBAGE_DSN_PASSWORD: "TODO_give_that_user_a_pw",
        },
      },
    });
  }
}
