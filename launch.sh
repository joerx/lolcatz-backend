#!/bin/sh

# Launcher script for development machines. Production machines should use systemd or docker
# to configure the application before launch.

make clean build

OS=$(go env GOOS)
ARCH=$(go env GOARCH)

if [[ -f .env ]]; then
    source .env
fi

fatal() {
    echo $1
    exit 1
}

[[ -z ${S3_BUCKET} ]] && fatal "S3_BUCKET must be set"
[[ -z ${S3_REGION} ]] && fatal "S3_REGION must be set"


echo "Bucket: ${S3_BUCKET}"
echo "Region: ${S3_REGION}"
echo "Database: postgres://${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
echo "---"

./bin/lolcatz-backend-${OS}-${ARCH} \
    -cors-allow-origin='*' \
    -bucket=${S3_BUCKET} \
    -region=${S3_REGION} \
    -db-host=${DB_HOST} \
    -db-name=${DB_NAME} \
    -db-port=${DB_PORT} \
    -db-user=${DB_USER} \
    -db-password=${DB_PASSWORD}
