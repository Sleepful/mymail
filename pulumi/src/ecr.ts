import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";
import { vpcId } from "./vpc";
import { sgAllowHttp, sgAllowTls } from "./securityGroup";

const repo = new awsx.ecr.Repository("repo-mymail", {
  forceDelete: true,
});

const image = new awsx.ecr.Image("ecr-image-mymail", {
  repositoryUrl: repo.url,
  context: "..",
  platform: "linux/amd64",
});

export const imageUri = image.imageUri

// const ecrApi = new aws.ec2.VpcEndpoint("vpcendpoint-ecrApi", {
//   vpcId: vpcId,
//   serviceName: "com.amazonaws.us-east-1.ecr.api",
//   vpcEndpointType: "Interface",
//   securityGroupIds: [sgAllowHttp.id, sgAllowTls.id]
// });

