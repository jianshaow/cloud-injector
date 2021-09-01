# cloud-injector

~~~ shell
kubectl create ns bar
kubectl label ns bar injection=enabled
kubectl -n bar apply -f test/test-cm.yaml
kubectl -n bar apply -f test/test-pod.yaml
~~~