# How To Setup an Application in AWS EC2

In this tutorial, we're going to use AWS free tier to deploy your barebones app to the cloud. By the end of this doc, your app will be deployed, talking to RDS, and open to the internet.

Prereqs:
1. An application that works. That is, your app should be running locally/in CI and talking to mySQL already.
1. Docker image of your app published to docker Hub. This isn't terrible difficult. (will link an article).
1. An AWS account.

Outcomes:
1. An ALB (Application Load Balancer) that is open for the internet to talk to.
1. A Cluster, Service, and Task with one EC2 instance (t3.micro) running your app via Docker image.
1. An RDS instance for your MySQL database.
1. A VPC (Virtual Private Cloud) and a Security Group that allows all of your AWS resources to talk to each other.
1. A HealthCheck for your Service that makes sense.

A few caveats:
1. The following guide is based off [this AWS article](https://aws.amazon.com/getting-started/hands-on/deploy-docker-containers/), but I found the detail to be insufficient. Also, it appears that it no longer falls in the "free tier" (I'm not sure how AWS charges for FarGate).
1. I'm not yet an AWS security-minded expert, so a lot of what I'll suggest in this article is intended to get your app on its feet rather than "the right choice". I do not claim to be experienced in this at all yet; I'm just sharing what I've found so far.

Let's dive in.

## Deploying your app

Let's do this in three stages: Persistence, Application, and Everything Else. So first we'll get your RDS stood up, then we'll deploy your docker image, and last of all we'll tweak all of the little things.

## Persistence (RDS)

I followed the [easy-enough instructions here](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Tutorials.WebServerDB.CreateDBInstance.html) to get my RDS set up. I'd suggest the same.

Other Reading:
- [RDS Getting Started Resources](https://aws.amazon.com/rds/resources/)
- [RDS Free Tier Details](https://aws.amazon.com/rds/free/)
- [RDS Tutorials](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Tutorials.html)
- [RDS User Guide](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Welcome.html)
- [Setup RDS with Web Server Example](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Tutorials.WebServerDB.CreateWebServer.html) (I found this tutorial lacking because it doesn't explain how to run the App from a Docker image)

## Application

Alright, so here's where things get fun. So if we start with [this article](https://aws.amazon.com/getting-started/hands-on/deploy-docker-containers/), and then tweak it a little, it'll be all-good.

1. Create an ECS Cluster. On [this page](https://console.aws.amazon.com/ecs/home#/clusters), click `Create Cluster`. [!Create Cluster Image](images/create_cluster.png)
1. Select `EC2 Linux + Networking` and `Next Step`. [!Create EC2 Linux](images/create_linux_cluster.png)
1. Choose the starting values for your Cluster.
   1. Set `Cluster name` to `my-app-cluster`.
   1. Set `EC2 instance type` to `t3.micro` for free-tier.
   1. For `VPC`, you should choose the one listed in your RDS under `Connectivity & security`'s `Networking` portion.
   1. Click `Create`. Once the Cluster is created, you should be able to see it as a card on [this page](https://console.aws.amazon.com/ecs/home#/clusters).
1. Setup the Service inside of the Cluster.
   - TODO
1. Setup the Task for the Service.
   - TODO
1. Create a Load Balancer for the Service.
   - TODO

Badda-bing, badda-boom. Now we'll just have to smooth out the cracks.

## Everything Else

Alright, so here's some of the mistakes I made, what they looked like, and how I resolved them.

- A `connection refused` or `i/o timeout` message coming out of that Task's logs when trying to connect to the RDS.
    - Likely, these both have the same root cause. Make sure your RDS and your Cluster are in the same VPC and have the same Security group. Then, make sure the Security Group on RDS has an "inbound rule" to accept traffic from your EC2 Cluster on port `3306` (or whatever port you're using).
- A `ERROR 1045 (28000): Access denied for user` message from the Task's logs when trying to connect to the RDS.
    - Hooray! This just means your username/password for accessing the DB is wrong and you need to update it in the DSN.
- Your app successfully connects to RDS (via logs), but you can't access your app via the ALB endpoint (it keeps saying 503).
    - For me, this happened because the ALB needs to health check anything it routes traffic to. So this means a few things:
        - Make sure your app is hosting HTTP traffic on the same port the ALB is health checking. For me, my app defaults to serving at `:8080` instead of `:80`, and this needed to be updated via environment variable.
        - Make sure your ALB is doing something sane for a health check. By default, it should assume "healthy" if it gets a 200 back from the path `/`, but my target group health check was messing that up. Find your HealthChecks by accessing TargetGroups [here](https://console.aws.amazon.com/ec2/home#TargetGroups;sort=targetGroupName), and select the one for your load-balancer. Then, at the bottom, select the `Health checks` tab and make sure the path looks like `/` (not `/my-service-name`). [!Good ALB Health Check](images/good_alb_healthcheck.png)