# How To Setup an Application in AWS EC2

In this tutorial, we're going to use AWS free tier to deploy your barebones app to the cloud. By the end of this doc, your app will be deployed, talking to RDS, and open to the internet.

Prereqs:
1. An application that works. That is, your app should be running locally/in CI and talking to mySQL already.
1. Docker image of your app published to docker Hub. This isn't terrible difficult. [do it in CI](https://docs.github.com/en/actions/language-and-framework-guides/publishing-docker-images), [some tutorial](https://ropenscilabs.github.io/r-docker-tutorial/04-Dockerhub.html), [another article](https://hackernoon.com/publish-your-docker-image-to-docker-hub-10b826793faf).
1. An AWS account.

Outcomes:
1. An RDS instance for your MySQL database.
1. An ALB (Application Load Balancer) that is open for the internet to talk to.
1. A Cluster, Service, and Task with one EC2 instance (t3.micro) running your app via Docker image.
1. A VPC (Virtual Private Cloud) and a Security Group that allows all of your AWS resources to talk to each other.
1. A HealthCheck for your Service that makes sense.

A few caveats:
1. The following guide is based off [this AWS article](https://aws.amazon.com/getting-started/hands-on/deploy-docker-containers/), but I found the detail to be insufficient. Also, it appears that it no longer falls in the "free tier" (since I'm not sure how AWS FarGate falls into the Free Tier).
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

1. Create an ECS Cluster. On [this page](https://console.aws.amazon.com/ecs/home#/clusters), click `Create Cluster`. [!Create Cluster Image](./images/create_cluster.png)
1. Select `EC2 Linux + Networking` and `Next Step`. [!Create EC2 Linux](./images/create_linux_cluster.png)
1. Choose the starting values for your Cluster.
   1. Set `Cluster name` to `my-app-cluster`.
   1. Set `EC2 instance type` to `t3.micro` for free-tier.
   1. For `VPC`, you should choose the one listed in your RDS under `Connectivity & security`'s `Networking` portion.
   1. Click `Create`. Once the Cluster is created, you should be able to see it as a card on [this page](https://console.aws.amazon.com/ecs/home#/clusters).
1. Setup the Task for the Service.
   1. Navigate to the Task Definition page (probably [here](https://console.aws.amazon.com/ecs/home#/taskDefinitions)).
   1. Click `Create new Task Definition`. [!Task Create page](./images/task_definition_page.png)
   1. Select `EC2`.
   1. Give your task a name (`my-app-task` will do).
   1. Click `Add Container`.
   1. Here is where you get to define your Docker container stats. This repo looks similar to [!this](./images/docker_container_defintion.png).
      - Scroll partway down to `Environment variables` and add one like this: [!env var in task definition](./images/envvar_task_definition.png)
         - PROTIP: You can edit these in future revisions by scrolling to bottom, selecting `Configure via JSON`, and then updating the JSON field `environment` under `containerDefinitions`. [!edit JSON](./images/edit_task_json.png) [!see JSON fields](./images/edit_task_json_fields.png)
      - The rest of these options can be configured if you like, but I left most of them alone.
      - Click `Add` when you're finished
   1. Click `Create`. You'll use this task when creating the Service.
1. Create a Load Balancer for the Service.
   1. Navigate to [the Load Balancer page](https://console.aws.amazon.com/ec2/v2/home#LoadBalancers:).
   1. Push `Create Load Balancer` and then `Create` under `Application Load Balancer`.
   1. Give it a name (to be used for the Service) and select your VPC and zones.
   1. Click through the rest of the pages and create the load balancer.
1. Setup the Service inside of the Cluster.
   1. Navigate to the cluster page (which looks like this: `https://console.aws.amazon.com/ecs/home#/clusters/<your-cluster-name>/services`).
   1. Click `Create` under the `Service` tab. [!New Cluster page](./images/new_cluster.png)
   1. Choose `EC2` for `Launch type`.
   1. Choose your task definition for `Task Definition`.
   1. Give the service a name (perhaps `my-app-service` :shrug:).
   1. Set `1` for `Number of tasks`.
   1. Click `Next`.
   1. Choose your VPC and subnets.
   1. Choose `Application Load Balancer` and then select the one you created.
   1. Click through and create.

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

## Recap

So now you should have an app deployed out into the cloud! Congrats! There's a lot of hullabaloo involved in this whole process, and software engineering is always changing. Best wishes, and good luck!