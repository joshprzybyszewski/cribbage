import * as cdk from "@aws-cdk/core";
import * as ecs from "@aws-cdk/aws-ecs";
import * as ecsPatterns from "@aws-cdk/aws-ecs-patterns";

export class InfrastructureStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const serverPort = 8080;

    new ecsPatterns.ApplicationLoadBalancedFargateService(this, `${id}-ecs`, {
      memoryLimitMiB: 1024,
      cpu: 512,
      taskImageOptions: {
        image: ecs.ContainerImage.fromAsset("../", {
          exclude: ["infrastructure/*"],
        }),
        // here's how we'd grab the image from docker hub:
        // image: ecs.ContainerImage.fromRegistry(
        //   "joshprzybyszewski/cribbage:latest"
        // ),
        containerPort: serverPort,
        environment: {
          CRIBBAGE_DB: "memory",
          CRIBBAGE_RESTPORT: serverPort.toString(),
        },
      },
    });
  }
}
