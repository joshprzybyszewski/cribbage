import * as cdk from '@aws-cdk/core';
import { ARecord, HostedZone, RecordSet, RecordTarget, RecordType } from '@aws-cdk/aws-route53';
import { LoadBalancerTarget } from '@aws-cdk/aws-route53-targets';
import { ILoadBalancerV2 } from '@aws-cdk/aws-elasticloadbalancingv2';

export interface DNSStackProps extends cdk.StackProps {
    loadBalancer: ILoadBalancerV2;
}

export class DNSStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props: DNSStackProps) {
        super(scope, id, props);

        const zone = HostedZone.fromLookup(this, `${id}-hz`, {
            domainName: 'hobbycribbage.com',
        });
        new ARecord(this, `${id}-r53-a-record`, {
            zone,
            target: RecordTarget.fromAlias(
                new LoadBalancerTarget(props.loadBalancer),
            ),
        });
    }
}
