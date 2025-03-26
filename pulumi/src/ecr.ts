import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";

const repo = new awsx.ecr.Repository("repo-mymail", {
  forceDelete: true,
});

const image = new awsx.ecr.Image("ecr-image-mymail", {
  repositoryUrl: repo.url,
  context: "../..",
  platform: "linux/amd64",
});

export const imageUri = image.imageUri
