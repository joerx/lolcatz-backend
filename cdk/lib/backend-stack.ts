import * as cdk from "@aws-cdk/core";
import * as ec2 from "@aws-cdk/aws-ec2";
import * as ecs from "@aws-cdk/aws-ecs";
import * as rds from "@aws-cdk/aws-rds";
import * as logs from "@aws-cdk/aws-logs";
import * as s3 from "@aws-cdk/aws-s3";
import * as sm from "@aws-cdk/aws-secretsmanager";
import * as lbv2 from "@aws-cdk/aws-elasticloadbalancingv2";
import * as r53 from "@aws-cdk/aws-route53";
import * as acm from "@aws-cdk/aws-certificatemanager";
import * as ecr from "@aws-cdk/aws-ecr";
import * as sns from "@aws-cdk/aws-sns";
import * as s3n from "@aws-cdk/aws-s3-notifications";
import { LoadBalancerTarget } from "@aws-cdk/aws-route53-targets";
import { CfnOutput } from "@aws-cdk/core";
import { DnsProps, RegistryImageProps } from "./shared";
import * as secretsmanager from "@aws-cdk/aws-secretsmanager";

export interface BackendStackProps extends cdk.StackProps {
  dns: DnsProps;
  image: RegistryImageProps;
}

export class BackendStack extends cdk.Stack {
  public readonly apiUrlOutput: CfnOutput;

  public readonly ecsService: ecs.IBaseService;

  public readonly ecrRepository: ecr.IRepository;

  constructor(scope: cdk.Construct, id: string, props: BackendStackProps) {
    super(scope, id, props);

    const region = cdk.Stack.of(this).region;
    const applicationPort = 3000;
    const apiUrl = `${props.dns.recordName}.${props.dns.zoneName}`;

    // VPC
    const vpc = new ec2.Vpc(this, "Vpc", {
      cidr: "10.0.0.0/16",
      maxAzs: 3,
      natGateways: 1,
      subnetConfiguration: [
        {
          subnetType: ec2.SubnetType.PUBLIC,
          name: "public",
        },
        {
          subnetType: ec2.SubnetType.PRIVATE,
          name: "private",
        },
      ],
      enableDnsHostnames: true,
      enableDnsSupport: true,
    });

    // Hosted zone, DNS cert

    const hostedZone = r53.HostedZone.fromLookup(this, "DnsZone", {
      domainName: props.dns.zoneName,
    });

    const cert = new acm.DnsValidatedCertificate(this, "AcmCert", {
      domainName: apiUrl,
      hostedZone,
    });

    // S3 bucket for storing images

    const bucket = new s3.Bucket(this, "ImageStore", {
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      cors: [
        {
          allowedOrigins: ["*"],
          allowedMethods: [s3.HttpMethods.GET],
        },
      ],
    });

    const topic = new sns.Topic(this, "ImageUploadNotifications");
    bucket.addEventNotification(s3.EventType.OBJECT_CREATED, new s3n.SnsDestination(topic));

    // Security groups

    const serviceSg = new ec2.SecurityGroup(this, "ServiceSg", {
      vpc,
    });

    const lbSg = new ec2.SecurityGroup(this, "LbSg", {
      vpc,
    });

    serviceSg.addIngressRule(lbSg, ec2.Port.tcp(applicationPort), "LB ingress", false);

    lbSg.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(80), "HTTP ingress");
    lbSg.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(443), "HTTPS ingress");

    // Database stuff

    const dbPasswordSecret = new sm.Secret(this, "DbPassword", {
      generateSecretString: {
        includeSpace: false,
        excludePunctuation: true,
      },
    });
    const dbUsername = "postgres";
    const dbName = "lolcatz";

    const dbSg = new ec2.SecurityGroup(this, "DbSg", {
      vpc,
      allowAllOutbound: false,
    });

    dbSg.addIngressRule(serviceSg, ec2.Port.tcp(5432), "Service ingress to DB", false);

    const db = new rds.DatabaseInstance(this, "Database", {
      engine: rds.DatabaseInstanceEngine.postgres({
        version: rds.PostgresEngineVersion.VER_12,
      }),
      vpc,
      vpcPlacement: {
        subnetType: ec2.SubnetType.PRIVATE,
      },
      credentials: {
        username: dbUsername,
        password: dbPasswordSecret.secretValue,
      },
      databaseName: dbName,
      multiAz: false,
      instanceType: ec2.InstanceType.of(ec2.InstanceClass.T3, ec2.InstanceSize.SMALL),
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      backupRetention: cdk.Duration.days(0),
      deleteAutomatedBackups: true,
      deletionProtection: false,
      cloudwatchLogsRetention: logs.RetentionDays.THREE_DAYS,
      securityGroups: [dbSg],
    });

    // ECS service + container

    const ecsCluster = new ecs.Cluster(this, "EcsCluster", {
      vpc,
    });

    const lg = new logs.LogGroup(this, "TaskLogs");

    const taskDef = new ecs.FargateTaskDefinition(this, "Backend", {});

    if (taskDef.taskRole) {
      bucket.grantReadWrite(taskDef.taskRole);
    }

    const containerCmd = [
      `-bind=":${applicationPort}"`,
      `-cors-allow-origin="*"`,
      `-bucket="${bucket.bucketName}"`,
      `-region="${region}"`,
      `-db-host="${db.dbInstanceEndpointAddress}"`,
      `-db-port="${db.dbInstanceEndpointPort}"`,
      `-db-name="${dbName}"`,
      `-db-user="${dbUsername}"`,
      '-db-password="${DB_PASSWORD}"',
    ].join(" ");

    // Docker image
    const imageName = props.image.repo + ":" + props.image.tag;
    const credentials = secretsmanager.Secret.fromSecretName(this, "ContainerImageCreds", props.image.secretName);
    const image = ecs.ContainerImage.fromRegistry(imageName, {credentials});

    const container = taskDef.addContainer("backend", {
      image,
      logging: ecs.LogDriver.awsLogs({
        streamPrefix: "backend",
        logGroup: lg,
      }),
      command: [containerCmd],
      secrets: {
        DB_PASSWORD: ecs.Secret.fromSecretsManager(dbPasswordSecret),
      },
    });

    container.addPortMappings({
      containerPort: applicationPort,
    });

    // ECS Service
    const ecsService = new ecs.FargateService(this, "Service", {
      cluster: ecsCluster,
      taskDefinition: taskDef,
      desiredCount: 1,
      minHealthyPercent: 50,
      securityGroups: [serviceSg],
      deploymentController: {
        type: ecs.DeploymentControllerType.ECS,
      },
      vpcSubnets: {
        subnetType: ec2.SubnetType.PRIVATE,
      },
    });

    this.ecsService = ecsService;

    // Load balancer, target group, etc.
    const lb = new lbv2.ApplicationLoadBalancer(this, "ServiceLb", {
      vpc: vpc,
      internetFacing: true,
      securityGroup: lbSg,
    });

    lb.addRedirect({
      sourcePort: 80,
      sourceProtocol: lbv2.ApplicationProtocol.HTTP,
      targetPort: 443,
      targetProtocol: lbv2.ApplicationProtocol.HTTPS,
    });

    // SSL Listener
    const secureListener = lb.addListener("secure", {
      port: 443,
      protocol: lbv2.ApplicationProtocol.HTTPS,
      sslPolicy: lbv2.SslPolicy.RECOMMENDED,
      certificates: [cert],
    });

    secureListener.addTargets("ServiceFleet", {
      targets: [ecsService],
      port: applicationPort,
      protocol: lbv2.ApplicationProtocol.HTTP,
      healthCheck: {
        enabled: true,
        port: `${applicationPort}`,
        protocol: lbv2.Protocol.HTTP,
        path: "/health",
      },
    });

    new r53.ARecord(this, "BackendDnsRecord", {
      target: r53.RecordTarget.fromAlias(new LoadBalancerTarget(lb)),
      zone: hostedZone,
      recordName: props.dns.recordName,
    });

    this.apiUrlOutput = new CfnOutput(this, "apiUrlOutput", {
      value: apiUrl,
    });
  }
}
