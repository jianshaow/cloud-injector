# cloud-injector

## Build go manually (optional)
~~~ shell
CGO_ENABLED=0 GOOS=linux go build -a -v -o build/pod-injector cmd/main.go
~~~

## Run locally (for test)
~~~ shell
go run cmd/main.go -v=2 --cert-file=test/server.cer --key-file=test/server.key --patch-file=test/patch.json
~~~

## Build docker
~~~ shell
docker build -t jianshao/pod-injector:0.0.1 .
docker push jianshao/pod-injector:0.0.1
~~~

## Run with docker
~~~ shell
docker run -d --name pod-injector --rm -v $PWD/test:/certs -p 8443:8443 jianshao/pod-injector:0.0.1 pod-injector -v=2
~~~

## Verify locally
~~~ shell
curl -v --cacert test/ca.cer -H"Content-Type: application/json" https://localhost:8443/inject -d@test/request.json
curl -s --cacert test/ca.cer -H"Content-Type: applicationjson" https://localhost:8443/inject -d@test/request.json|jq -r '.response.patch'|base64 -d
~~~

## Create webhook
~~~ shell
kubectl create ns injector
kubectl apply -f manifests/injector-certs.yaml
kubectl apply -f manifests/injector-patchs-example.yaml
kubectl apply -f manifests/injector-deploy.yaml
kubectl apply -f manifests/injector-webhook.yaml
~~~

## Demo pod injection
~~~ shell
kubectl create ns bar
kubectl label ns bar injection=enabled
kubectl -n bar apply -f test/test-cm.yaml
kubectl -n bar apply -f test/test-pod.yaml
~~~
