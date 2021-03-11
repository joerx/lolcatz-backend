import { Construct, StackProps, Stage } from "@aws-cdk/core";
import { BackendStack } from "./backend-stack";
import { DnsProps, RegistryImageProps } from "./shared";

export interface BackendStageProps extends StackProps {
  dns: DnsProps;
  image: RegistryImageProps
}

export class BackendStage extends Stage {
  constructor(scope: Construct, id: string, props: BackendStageProps) {
    super(scope, id, props);

    // Backend application stack
    new BackendStack(this, "LolcatzBackend", props);
  }
}
