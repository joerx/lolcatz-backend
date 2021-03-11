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
  name: "ghcr.io/joerx/lolcatz-backend:latest"
}

// How is branching going to work? Pipeline is deployed only on change to master?
// What about multiple team members working on the same application? How do they 
// test their changes locally? 
// Presumably the deployed changes are staging or pre-prod, so they represent the 
// integration stage. For local development, engineers use local deployment (e.g.
// SAM, docker compose, etc.), unit tests, etc. Integration tests are either part
// of the CDK pipeline or a separate CI pipeline.
// So the workflow looks roughly like this:
//
// - Code, build and test locally
// - Push changes, create a PR
// - Automated unit and integration tests
// - PR review, appproval
// - Merge to master updates pipeline, deploys preprod
// - Additional stages, manual approval, etc.

const source: GitHubSourceProps = {
  repo: "lolcatz-backend",
  owner: "joerx",
  oauthToken: SecretValue.secretsManager("github/oauth-token"),
  branch: "master",
  subdirectory: "cdk",
};

new BackendCdkPipeline(app, "LolcatzBackendCdkPipeline", { source, dns, env, image });
