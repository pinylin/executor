set -ex

# git pull

version=$1
echo "version: $version"

# export GOPATH=/Users/pinglin/workspace/unicom:$GOPATH
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tasksService tasksSvc/main.go

docker build -t executor/tasksSvc:$version -f Dockerfile .

docker tag executor/tasksSvc:$version 47.98.200.101:5000/executor/tasksSvc:$version

docker push 47.98.200.101:5000/executor/tasksSvc:$version

docker rmi $(docker images -f "dangling=true" -q)
docker images