#!/usr/bin/env node

import "source-map-support/register"; 
import * as cdk from "@aws-cdk/core";
import { SecretValue } from "@aws-cdk/core";
import { BackendCdkPipeline } from "../lib/backend-pipeline";
import { GitHubSourceProps } from "../lib/shared";

const app = new cdk.App();

const env = {
  account: process.env.CDK_DEFAULT_ACCOUNT,
  region: process.env.CDK_DEFAULT_REGION,
};

const dns = {
  zoneName: "lolcatz.tv",
  recordName: "api-dev",
};

const image = {
  secretName: "docker/credentials/ghcr",
  repo: "ghcr.io/joerx/lolcatz-backend",
  tag: "latest"
}

const source: GitHubSourceProps = {
  repo: "lolcatz-backend",
  owner: "joerx",
  oauthToken: SecretValue.secretsManager("github/oauth-token"),
  branch: "master",
  subdirectory: "cdk",
};

new BackendCdkPipeline(app, "LolcatzBackendCdkPipeline", { source, dns, env, image });
