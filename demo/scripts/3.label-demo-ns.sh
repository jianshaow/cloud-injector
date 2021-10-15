#set -x
kubectl label ns demo injection=enabled
kubectl get ns demo --show-labels=true
#set +x
