# cloud-injector

## Build go
~~~ shell
CGO_ENABLED=0 GOOS=linux go build -a -v -o build/pod-injector cmd/main.go
~~~

## Build docker
~~~ shell
docker build -t jianshao/pod-injector:0.0.1 docker/
docker push jianshao/pod-injector:0.0.1
~~~

## Create webhook
~~~ shell
kubectl create ns injector
kubectl -n injector apply -f manifests/injector-certs.yaml
kubectl -n injector apply -f manifests/injector-deploy.yaml
kubectl -n injector apply -f manifests/injector-deploy.yaml
~~~

## Demo pod injection
~~~ shell
kubectl create ns bar
kubectl label ns bar injection=enabled
kubectl -n bar apply -f test/test-cm.yaml
kubectl -n bar apply -f test/test-pod.yaml
~~~
