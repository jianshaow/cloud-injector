. sourceme.sh

while ( true ); do
  unavailable_replicas=$(kubectl -n demo get deploy demo-app -ojsonpath='{.status.unavailableReplicas}')
  if [ "$unavailable_replicas" != "1" ]; then
    kubectl -n demo get po
    break
  fi
  kubectl -n demo get po
done
green 'application is ready!'
