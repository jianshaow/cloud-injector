#set -x
kubectl label ns demo injection-
kubectl get ns demo --show-labels=true
#set +x
