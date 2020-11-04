import * as codepipeline from "@aws-cdk/aws-codepipeline";
import * as actions from "@aws-cdk/aws-codepipeline-actions";
import { Construct, Stack, StackProps } from "@aws-cdk/core";
import { CdkPipeline, SimpleSynthAction } from "@aws-cdk/pipelines";
import { BackendStage } from "./backend-stage";
import { DnsProps, EcrImageProps, GitHubSourceProps } from "./shared";

export interface BackendCdkPipelineProps extends StackProps {
  source: GitHubSourceProps;
  dns: DnsProps;
}

export class BackendCdkPipeline extends Stack {
  constructor(scope: Construct, id: string, props: BackendCdkPipelineProps) {
    super(scope, id, props);

    const { account, region } = Stack.of(this);

    const sourceArtifact = new codepipeline.Artifact();
    const cloudAssemblyArtifact = new codepipeline.Artifact();

    const pipeline = new CdkPipeline(this, "BackendCdkPipeline", {
      cloudAssemblyArtifact,

      sourceAction: new actions.GitHubSourceAction({
        actionName: "GitHubSource",
        output: sourceArtifact,
        ...props.source,
      }),

      synthAction: SimpleSynthAction.standardNpmSynth({
        sourceArtifact: sourceArtifact,
        cloudAssemblyArtifact,
        subdirectory: props.source.subdirectory,
        environment: {
          privileged: true,
        },
        environmentVariables: {
          CDK_DEFAULT_ACCOUNT: { value: account },
          CDK_DEFAULT_REGION: { value: region },
        },
      }),
    });

    pipeline.addApplicationStage(
      new BackendStage(this, "PreProd", {
        dns: props.dns,
        env: props.env,
      })
    );
  }
}
