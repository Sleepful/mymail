import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";
import { publicSubnetIds, vpcId, vpc } from "./vpc";
import { imageUri } from "./ecr";
import { sgAllow8090Http, sgAllowHttp, sgAllowTls } from "./securityGroup";


// const example = new aws.acm.Certificate("example", {});
const CERT_ARN = "arn:aws:acm:us-east-1:791951373389:certificate/26af2906-829d-4fa2-abd3-29f6446a52c8"

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
      certificateArn: CERT_ARN
    },
  ],
  defaultTargetGroup: {
    port: 8090,
  },
  securityGroups: [sgAllowTls.id, sgAllowHttp.id],
  subnetIds: publicSubnetIds,
});

const cluster = new aws.ecs.Cluster("ecs-cluster-mymail", {});

const service = new awsx.ecs.FargateService("ecs-fargate-mymail", {
  cluster: cluster.arn,
  networkConfiguration: {
    subnets: publicSubnetIds,
    securityGroups: [sgAllow8090Http.id],

    // NOTE: a NAT Gateway would be better, for now use cheaper EIP option,
    // required to initiate outbound connections: ecr, apis, etc.
    assignPublicIp: true,
  },
  desiredCount: 1,
  // TODO:
  // blue green deployments
  // defaults to "ECS" type:
  // - default: https://github.com/pulumi/pulumi-aws/blob/master/sdk/dotnet/Ecs/Outputs/ServiceDeploymentController.cs#L17
  // - green/blue updates: 
  //    - https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-type-bluegreen.html
  //    - https://docs.aws.amazon.com/AmazonECS/latest/developerguide/create-blue-green.html
  // - api: https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_DeploymentController.html
  // - rolling updates: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deployment-type-ecs.html
  // deploymentController: { type: "CODE_DEPLOY" },
  //
  // Alternative with rolling deployment:
  deploymentMaximumPercent: 200,
  deploymentMinimumHealthyPercent: 100,
  // deploymentMaximumPercent: 100,
  // deploymentMinimumHealthyPercent: 0,
  taskDefinitionArgs: {
    container: {
      name: "myPocketbaseService",
      image: imageUri,
      cpu: 128,
      memory: 512,
      essential: true,
      portMappings: [
        {
          containerPort: 8090,
          targetGroup: lb.defaultTargetGroup,
        },
      ],
    },
  },
});

export const url = pulumi.interpolate`http://${lb.loadBalancer.dnsName}`;
