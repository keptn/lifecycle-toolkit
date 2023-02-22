At the moment, the lifecycle controller needs *cert-manager* to be installed. Therefore, you can install cert-manager using:

<!-- 
[cert-manager](https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml)
-->
```
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
kubectl wait --for=condition=Available deployment/cert-manager-webhook -n cert-manager --timeout=60s
```

After that, you can install the lifecycle toolkit using the current release manifest:
<!---x-release-please-start-version-->
```
kubectl apply -f https://github.com/keptn/lifecycle-toolkit/releases/download/v0.5.0/manifest.yaml
kubectl wait --for=condition=Available deployment/klc-controller-manager -n keptn-lifecycle-toolkit-system --timeout=120s
```
<!---x-release-please-end-->

Now, the Lifecycle Toolkit and its dependency is installed and ready to use.