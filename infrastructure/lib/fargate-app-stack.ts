import * as ec2 from "@aws-cdk/aws-ec2";
import * as ecs_patterns from "@aws-cdk/aws-ecs-patterns";
import * as cdk from "@aws-cdk/core";
import * as ecs from "@aws-cdk/aws-ecs";
import {ISecret, Secret} from '@aws-cdk/aws-secretsmanager';

export interface FargateAppStackProps extends cdk.StackProps {
  vpc: ec2.Vpc;
  rdsEndpoint: string
  rdsDbUser: string
  rdsDbPassword?: string
}


export class FargateAppStack extends cdk.Stack {


  constructor(scope: cdk.Construct, id: string, props: FargateAppStackProps) {
    super(scope, id, props);


    const secret = Secret.fromSecretAttributes(this, props.rdsDbPassword || 'umwaat', {
      secretArn: 'arn:aws:secretsmanager:<region>:<organisationId>:secret:ImportedSecret-fPNc0O',
    });

    // TODO we need to pass these vars in as props that were created in the RDS-stack file
    // we may be able to use this dynamoDB creation example as a starter: https://github.com/aws-samples/aws-cdk-examples/blob/18f3429414cc80223a04a04d0249b12a7ae13cbb/typescript/api-cors-lambda-crud-dynamodb/index.ts#L10-L21
    var dsnPw = secret.secretValue.toString() || props.rdsDbPassword || 'TODOsetThePassword';
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
          CRIBBAGE_DSN_PASSWORD: dsnPw.toString(),
        },
      },
    });
  }
}
