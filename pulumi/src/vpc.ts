import * as pulumi from "@pulumi/pulumi";
import * as awsx from "@pulumi/awsx";
import * as aws from "@pulumi/aws";
import { sgAllowHttp, sgAllowTls } from "./securityGroup";

// Allocate a new VPC with the default settings.
export const vpc = new awsx.ec2.Vpc("vpc-mymail", {
  cidrBlock: "10.0.4.0/24",
  numberOfAvailabilityZones: 3,
  subnetSpecs: [
    {
      type: awsx.ec2.SubnetType.Public,
      cidrMask: 28,
    },
  ], natGateways: {
    strategy: awsx.ec2.NatGatewayStrategy.None,
  },
  subnetStrategy: awsx.ec2.SubnetAllocationStrategy.Auto,
});


// Export a few properties to make them easy to use.
export const vpcId = vpc.vpcId;
// export const privateSubnetIds = vpc.privateSubnetIds;
export const publicSubnetIds = vpc.publicSubnetIds;
