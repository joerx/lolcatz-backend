import { SecretValue } from "@aws-cdk/core";

export interface DnsProps {
  recordName?: string;
  zoneName: string;
}

export interface GitHubSourceProps {
  oauthToken: SecretValue;
  repo: string;
  owner: string;
  branch?: string;
  subdirectory?: string;
}

export interface EcrImageProps {
  repository: string;
  tag?: string;
}
