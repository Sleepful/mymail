import * as aws from "@pulumi/aws";
import { vpc } from "./vpc";

const allowTls = new aws.ec2.SecurityGroup("allow_tls", {
  name: "allow_tls",
  description: "Allow TLS inbound traffic and all outbound traffic",
  vpcId: vpc.vpcId,
  tags: {
    Name: "allow_tls",
  },
});

const allowTlsIpv4 = new aws.vpc.SecurityGroupIngressRule("allow_tls_ipv4", {
  securityGroupId: allowTls.id,
  cidrIpv4: "0.0.0.0/0",
  fromPort: 443,
  ipProtocol: "tcp",
  toPort: 443,
})

const allowAllTrafficIpv4 = new aws.vpc.SecurityGroupEgressRule("allow_all_traffic_ipv4", {
  securityGroupId: allowTls.id,
  cidrIpv4: "0.0.0.0/0",
  ipProtocol: "-1",
})

const allowHttp = new aws.ec2.SecurityGroup("allow_http", {
  name: "allow_http",
  description: "Allow HTTP inbound traffic and all outbound traffic",
  vpcId: vpc.vpcId,
  tags: {
    Name: "allow_http",
  },
});

const allowHttpIpv4 = new aws.vpc.SecurityGroupIngressRule("allow_http_ipv4", {
  securityGroupId: allowHttp.id,
  cidrIpv4: "0.0.0.0/0",
  fromPort: 80,
  ipProtocol: "tcp",
  toPort: 80,
})

const allowAllTrafficIpv4Http = new aws.vpc.SecurityGroupEgressRule("allow_all_traffic_ipv4_http", {
  securityGroupId: allowHttp.id,
  cidrIpv4: "0.0.0.0/0",
  ipProtocol: "-1",
})

const allow8090Http = new aws.ec2.SecurityGroup("allow_8090_http", {
  name: "allow_8090_http",
  description: "Allow HTTP inbound traffic on 8090 and all outbound traffic",
  vpcId: vpc.vpcId,
  tags: {
    Name: "allow_8090_http",
  },
});

const allow8090HttpIpv4 = new aws.vpc.SecurityGroupIngressRule("allow_8090_http_ipv4", {
  securityGroupId: allow8090Http.id,
  cidrIpv4: "0.0.0.0/0",
  fromPort: 8090,
  ipProtocol: "tcp",
  toPort: 8090,
})

const allowAllEgressHttp8090 = new aws.vpc.SecurityGroupEgressRule("allow_all_traffic_ipv4_http_8090", {
  securityGroupId: allow8090Http.id,
  cidrIpv4: "0.0.0.0/0",
  ipProtocol: "-1",
})
export const sgAllowHttp = allowHttp
export const sgAllow8090Http = allow8090Http
export const sgAllowTls = allowTls
