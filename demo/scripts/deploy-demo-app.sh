kubectl -n demo apply -f demo-deploy.yaml
while ( true ); do
  replicas=$(kubectl -n demo get deploy demo-app -ojsonpath='{.status.availableReplicas}')
  if [ "$replicas" == "1" ]; then
    kubectl -n demo get po
    break
  fi
  kubectl -n demo get po
done

url=$(kubectl -n demo get svc demo-app -ojsonpath='{"http://"}{.spec.clusterIP}{":"}{.spec.ports[?(@.name=="http")].port}')
echo application is started, access it in browser: $url
