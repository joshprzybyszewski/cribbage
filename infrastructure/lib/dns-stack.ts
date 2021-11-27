import * as cdk from '@aws-cdk/core';
import * as apigw from '@aws-cdk/aws-apigateway';
import { Function } from '@aws-cdk/aws-lambda';
import { ARecord, HostedZone, RecordTarget } from '@aws-cdk/aws-route53';
import * as route53Targets from '@aws-cdk/aws-route53-targets';
import * as acm from '@aws-cdk/aws-certificatemanager';
import { SPADeploy } from 'cdk-spa-deploy';

export interface DNSStackProps extends cdk.StackProps {
    lambda: Function;
}


export class DNSStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props: DNSStackProps) {
        super(scope, id, props);

        const zone = HostedZone.fromLookup(this, `${id}-hz`, {
            domainName: 'hobbycribbage.com',
        });
        // TODO it's safer to use `addResource` or `addMethod` or `addProxy` to point to this lambda
        // than to have `proxy: true,`, because that's gonna just route every request to this lambda.
        const restApi = new apigw.LambdaRestApi(this, 'cribbage-lambda-rest-api', {
            restApiName: 'lambda-rest-api',
            description: 'the bundled up golang server acting as a lambda',
            handler: props.lambda,
            proxy: true,
            deployOptions: {
                stageName: 'urlprefix',
                // tracingEnabled: true,
            },
            endpointTypes: [apigw.EndpointType.EDGE],
        });

        const cert = new acm.Certificate(this, 'cribbage-lamdba-cert', {
            domainName: '*.hobbycribbage.com',
            validation: acm.CertificateValidation.fromDns(zone),
        });
        const subdomain = 'lambda.hobbycribbage.com';
        restApi.addDomainName("domain_name", {
            domainName: subdomain,
            // securityPolicy: apigw.SecurityPolicy.TLS_1_2,
            certificate: cert,
            endpointType: apigw.EndpointType.EDGE,
        })

        new ARecord(this, `${id}-r53-a-record`, {
            zone: zone,
            recordName: subdomain,
            comment: 'A comment on the ARecord',
            target: RecordTarget.fromAlias(
                new route53Targets.ApiGateway(
                    restApi,
                ),
            ),
        });

        // TODO perhaps use https://github.com/aws-samples/aws-cdk-examples/blob/901ae3e11704fc378ade673f76f0eeae860a5daf/typescript/static-site/static-site.ts#L113-L127
        // as an example for deploying CDN-like
        // TODO re-architect the frontend so that the "backend" calls go to a different sub-domain to hit
        // the lambda, and all of hobbycribbage.com/ just goes to a CDN to get the page. It's a SPA, so
        // I don't want to overthink this.
        // https://github.com/cszczepaniak/oh-hell-scorecard/blob/7ea3b0ff4ce229d035bd5b100560399ac56be48c/infrastructure/lib/infrastructure-stack.ts#L8-L12
        const zipFileName = 'spa-bundle';
        const spaDeploy = new SPADeploy(this, `spa-deploy-${id}`).createSiteFromHostedZone({
            zoneName: 'hobbycribbage.com',
            indexDoc: 'index.html',
            websiteFolder:  './' + zipFileName + '.zip',
        });
    }
}
