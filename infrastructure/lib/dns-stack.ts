import * as cdk from '@aws-cdk/core';
import * as apigw from '@aws-cdk/aws-apigateway';
import { Function } from '@aws-cdk/aws-lambda';
import { ARecord, HostedZone, RecordTarget } from '@aws-cdk/aws-route53';
import { ApiGateway, ApiGatewayDomain } from '@aws-cdk/aws-route53-targets';
import * as acm from '@aws-cdk/aws-certificatemanager';
import { SPADeploy } from 'cdk-spa-deploy';

export interface DNSStackProps extends cdk.StackProps {
    lambda: Function;
}


export class DNSStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props: DNSStackProps) {
        super(scope, id, props);

        // TODO somewhere between the r53 and the lambda we're losing the request
        const zone = HostedZone.fromLookup(this, `${id}-hz`, {
            domainName: 'hobbycribbage.com',
        });
        const restApi = new apigw.LambdaRestApi(this, 'cribbage-lambda-rest-api', {
            restApiName: 'lambda-rest-api',
            handler: props.lambda,
        });
        const cert = new acm.Certificate(this, 'cribbage-lamdba-cert', {
            domainName: '*.hobbycribbage.com',
            validation: acm.CertificateValidation.fromDns(zone),
        });
        const subdomain = 'lambda.hobbycribbage.com';
        const domain = restApi.addDomainName('DomainID', {
            certificate: cert,
            domainName: subdomain,
        })
        new ARecord(this, `${id}-r53-a-record`, {
            zone,
            recordName: subdomain,
            target: RecordTarget.fromAlias(
                new ApiGatewayDomain(domain),
            ),
        });

        // TODO perhaps use https://github.com/aws-samples/aws-cdk-examples/blob/901ae3e11704fc378ade673f76f0eeae860a5daf/typescript/static-site/static-site.ts#L113-L127
        // as an example for deploying CDN-like
        // TODO re-architect the frontend so that the "backend" calls go to a different sub-domain to hit
        // the lambda, and all of hobbycribbage.com/ just goes to a CDN to get the page. It's a SPA, so
        // I don't want to overthink this.
        // https://github.com/cszczepaniak/oh-hell-scorecard/blob/7ea3b0ff4ce229d035bd5b100560399ac56be48c/infrastructure/lib/infrastructure-stack.ts#L8-L12
        // new SPADeploy(this, id).createSiteFromHostedZone({
        //     indexDoc: 'index.html',
        //     websiteFolder: 'artifact',
        //     zoneName: 'hobbycribbage.com',
        // });
    }
}
