#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "@aws-cdk/core";
// import { SecretValue } from "@aws-cdk/core";
// import { BackendCdkPipeline } from "../lib/backend-pipeline";
import { BackendStack } from "../lib/backend-stack";
import { SecretValue } from "@aws-cdk/core";
import { BackendCdkPipeline } from "../lib/backend-pipeline";

const app = new cdk.App();

const env = {
  account: process.env.CDK_DEFAULT_ACCOUNT,
  region: process.env.CDK_DEFAULT_REGION,
};

const dns = {
  zoneName: "lolcatz.tv",
  recordName: "api-dev",
};

const source = {
  repo: "lolcatz-backend",
  owner: "joerx",
  oauthToken: SecretValue.secretsManager("github/oauth-token"),
  branch: "cdk-pipeline",
};

// new BackendStack(app, "LolcatzBackend", {
//   dns,
//   env,
// });

new BackendCdkPipeline(app, "LolcatzBackendCdkPipeline", { source, dns, env });
