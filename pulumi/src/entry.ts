import * as pulumi from "@pulumi/pulumi";
import * as awsx from "@pulumi/awsx";
import { vpcId } from "./vpc"
import { url } from "./ecs"

export const out = { vpcId, urlAlb: url }
