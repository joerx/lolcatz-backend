# LolCatz Backend

Manage funny cat pictures. Or any other kind of pictures, really.

## Build

```sh
go build -o bin/server .
```

## Usage

- Build the server binary
- Create an S3 bucket and set up credentials via the [default AWS credentials chain](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)
- Start the [frontend](https://github.com/joerx/lolcatz-frontend) (take note of the browser URL for the CORS origin header)
- Note: you may omit `-cors-allow-origin` but this is recommended _only for development purposes_
- Start the server:

```sh
bin/server -bucket <bucket> -region <region> -cors-allow-origin http://localhost:3000
```

## Database Config

Default database config assumes postgres on `localhost:5432`, database `lolcatz`, password `default`. To customize:

```sh
bin/server -bucket <bucket> -region <region>
  -db-host <my-db-hostname>
  -db-user <my-db-user>
  -db-password <my-db-password>
  -db-name <my-db-name>
```

## Change HTTP Port

By default the http server will bind to `localhost:8000` which is suitable for development. To customize, use the `-bind` option.

```sh
bin/server -bind=localhost:9000 ...
```

## Docker

Using GitHub Container Registry. See https://github.blog/2020-09-01-introducing-github-container-registry/.

Example GHCR Login with token stored in SSM:

```
export CR_PAT=$(aws ssm get-parameter --name /github/ghcr-push --query Parameter.Value --output text --with-decryption)
echo $CR_PAT | docker login ghcr.io -u USERNAME --password-stdin
```

Build and publish:

```bash
make docker-build
make docker-push
```

### ECR Push

Can be useful when used with AWS CodePipeline. Assuming you want to use your current default AWS profile:

```bash
AWS_ACCOUNT=$(aws sts get-caller-identity --query Account --output text)
AWS_REGION=$(aws configure get region)

docker build -t lolcatz-backend .

aws ecr get-login-password | docker login --username AWS --password-stdin ${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com
docker tag lolcatz-backend:latest ${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com/lolcatz-backend:latest
docker push ${AWS_ACCOUNT}.dkr.ecr.${AWS_REGION}.amazonaws.com/lolcatz-backend:latest
```
