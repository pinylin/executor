set -ex

# git pull

version=$1
echo "version: $version"

# export GOPATH=/Users/pinglin/workspace/unicom:$GOPATH
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wechatService wechatSvc/main.go

docker build -t executor/wechatSvc:$version -f Dockerfile .

docker tag executor/wechatSvc:$version 47.98.200.101:5000/executor/wechatSvc:$version

docker push 47.98.200.101:5000/executor/wechatSvc:$version

docker rmi $(docker images -f "dangling=true" -q)
docker images