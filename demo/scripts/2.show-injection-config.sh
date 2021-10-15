. sourceme.sh

h_blue '=============================== injection config ==============================='
#set -x
kubectl -n injector get cm injection-config -ojsonpath='{.data.injection\.yaml}'
#set +x
h_blue '================================================================================'
