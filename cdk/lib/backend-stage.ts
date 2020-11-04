import { Construct, StackProps, Stage } from "@aws-cdk/core";
import { BackendStack } from "./backend-stack";
import { DnsProps, EcrImageProps } from "./shared";

export interface BackendStageProps extends StackProps {
  dns: DnsProps;
}

export class BackendStage extends Stage {
  constructor(scope: Construct, id: string, props: BackendStageProps) {
    super(scope, id, props);

    // Backend application stack
    new BackendStack(this, "LolcatzBackend", props);
  }
}
