. sourceme.sh

h_blue '============================== injector namespace =============================='
#set -x
kubectl -n injector get all
#set +x
h_blue '================================================================================'
echo
read -p 'press enter...'
h_green '================================ demo namespace ================================'
#set -x
kubectl get ns demo --show-labels=true
kubectl -n demo get all
#set +x
h_green '================================================================================'
url=$(kubectl -n demo get svc demo-app -ojsonpath='{"http://"}{.spec.clusterIP}{":"}{.spec.ports[?(@.name=="http")].port}')
echo demo application url: $url
