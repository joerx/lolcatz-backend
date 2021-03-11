import * as codepipeline from "@aws-cdk/aws-codepipeline";
import * as actions from "@aws-cdk/aws-codepipeline-actions";
import { Construct, Stack, StackProps } from "@aws-cdk/core";
import { CdkPipeline, SimpleSynthAction } from "@aws-cdk/pipelines";
import { BackendStage } from "./backend-stage";
import { DnsProps, RegistryImageProps, GitHubSourceProps } from "./shared";

export interface BackendCdkPipelineProps extends StackProps {
  source: GitHubSourceProps;
  image: RegistryImageProps;
  dns: DnsProps;
}

export class BackendCdkPipeline extends Stack {
  constructor(scope: Construct, id: string, props: BackendCdkPipelineProps) {
    super(scope, id, props);

    // console.log("ENVIRONMENT VARIABLES");
    // console.log(process.env);
    // console.log("---------------------")
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

    // If image hasn't been built before the pipeline executes, the container will fail to start. ECS should keep trying to 
    // launch the container until the image ultimately becomes available.
    // During initial deploy this may leave the ECS service in a pending state
    // During subsequent deploys, the service should simply keep running the old version
    const image = {
      secretName : props.image.secretName,
      repo : props.image.repo,
      tag : process.env.CODEBUILD_RESOLVED_SOURCE_VERSION ? process.env.CODEBUILD_RESOLVED_SOURCE_VERSION.substring(0, 8) : props.image.tag
    }

    console.log("Deploy image", image)

    pipeline.addApplicationStage(
      new BackendStage(this, "PreProd", {
        dns: props.dns,
        env: props.env,
        image
      })
    );
  }
}
