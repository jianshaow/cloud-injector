set -x
kubectl -n demo rollout restart deploy demo-app
set +x
