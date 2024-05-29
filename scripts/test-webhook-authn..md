## Steps

1. Build Horus.
```bash
# It also pushes the built image to K3s.
$ ./scripts/build-image.sh
```

2. Ensure that to pod uses latest Horus.
```bash
$ kubectl get pods
NAME           READY   STATUS    RESTARTS   AGE
horus-royale   1/1     Running   0          17m

$ kubectl delete pods/horus-royal
$ kubectl get pods
NAME           READY   STATUS    RESTARTS   AGE
horus-cheese   1/1     Running   0          2s
```

3. Run test.
```bash
$ ./test-webhook-authn.mjs
```
