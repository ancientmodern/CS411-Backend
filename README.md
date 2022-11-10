# Illini Foodie Backend
## Requirement
- Go >= 1.19
- MySQL

## Deployment
1. Install the gcloud CLI locally ([doc](https://cloud.google.com/sdk/docs/install))
```shell
./google-cloud-sdk/install.sh
./google-cloud-sdk/bin/gcloud init
```
2. Enable gcloud auth application-default login ([doc](https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login))
```shell
source ~/.zshrc
gcloud auth application-default login
```
3. Install go pkg dependencies (either `go mod` or `go get`)
```shell
go mod tidy
# go get .
```
4. Build go executable file
```shell
go build main.go
```
5. Install and run Google Cloud SQL Auth proxy in the background ([doc](https://cloud.google.com/sql/docs/mysql/sql-proxy))
```shell
# the link is for M1 Mac, see doc for other platforms
curl -o cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.darwin.arm64
chmod +x cloud_sql_proxy
./cloud_sql_proxy -instances=cs411-team067:us-central1:sometimesnaive=tcp:3306 &
# to kill a bind
# ps | grep cloud_sql_proxy
# kill ...
```
6. Run the program with environment variables
```shell
DB_USER=${user} DB_PASS=${password} DB_NAME=test411 DB_PORT=3306 INSTANCE_HOST=127.0.0.1 ./main
```
## REST API format
https://docs.google.com/document/d/1A0gR0ikaMZ6yH2HnQQ4vFLo0zqM2efcOIIjumdQ0KfI/edit