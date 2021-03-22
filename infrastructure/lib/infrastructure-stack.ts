import * as cdk from "@aws-cdk/core";
import * as ecs from "@aws-cdk/aws-ecs";
import * as ecsPatterns from "@aws-cdk/aws-ecs-patterns";

export class InfrastructureStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    new ecsPatterns.ApplicationLoadBalancedFargateService(this, `${id}-ecs`, {
      ...props,
      memoryLimitMiB: 1024,
      cpu: 512,
      taskImageOptions: {
        image: ecs.ContainerImage.fromAsset("../", {
          exclude: ["infrastructure/*"],
        }),
        containerPort: 8080,
      },
    });
  }
}
