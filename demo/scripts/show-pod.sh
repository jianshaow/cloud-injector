. sourceme.sh

h_blue '================================ initContainers ================================'
set -x
kubectl -n demo get po -l app=demo-app -ojsonpath='{.items[0].spec.initContainers}' | jq
set +x
h_blue '================================================================================'
echo
read -p 'press enter...'
h_blue '================================== containers =================================='
set -x
kubectl -n demo get po -l app=demo-app -ojsonpath='{.items[0].spec.containers}' | jq
set +x
h_blue '================================================================================'
echo
read -p 'press enter...'
h_blue '==================================== volumes ==================================='
set -x
kubectl -n demo get po -l app=demo-app -ojsonpath='{.items[0].spec.volumes}' | jq
set +x
h_blue '================================================================================'
echo
