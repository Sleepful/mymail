import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";
import { publicSubnetIds, vpcId } from "./vpc";
import { imageUri } from "./ecr";
import { sgAllowHttp, sgAllowTls } from "./securityGroup";


// const example = new aws.acm.Certificate("example", {});

// Create a load balancer in the default VPC listening on port 80.
const lb = new awsx.lb.ApplicationLoadBalancer("lb-mymail", {
  listeners: [
    {
      // Redirect HTTP traffic on port 80 to port 443.
      port: 80,
      protocol: "HTTP",
      defaultActions: [
        {
          type: "redirect",
          redirect: {
            port: "443",
            protocol: "HTTPS",
            statusCode: "HTTP_301",
          },
        },
      ],
    },
    {
      // Accept HTTPS traffic on port 443.
      port: 443,
      protocol: "HTTPS",
    },
  ],
  defaultTargetGroup: {
    port: 8090,
  },
  securityGroups: [sgAllowTls.id, sgAllowHttp.id]
});

const cluster = new aws.ecs.Cluster("ecs-cluster-mymail", {});

const service = new awsx.ecs.FargateService("ecs-fargate-mymail", {
  cluster: cluster.arn,
  networkConfiguration: {
    subnets: publicSubnetIds,
    securityGroups: [sgAllowHttp.id],
  },
  desiredCount: 1,
  deploymentMaximumPercent: 0,
  taskDefinitionArgs: {
    container: {
      name: "myPocketbaseService",
      image: imageUri,
      cpu: 128,
      memory: 512,
      essential: true,
      portMappings: [
        {
          // containerPort: 8090,
          // containerPort: 80,
          containerPort: 8090,
          targetGroup: lb.defaultTargetGroup,
        },
      ],
    },
  },
});

export const url = pulumi.interpolate`http://${lb.loadBalancer.dnsName}`;
