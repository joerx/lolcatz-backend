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
bin/server -bucket=<bucket> -region=<region> -cors-allow-origin=http://localhost:3000
```
