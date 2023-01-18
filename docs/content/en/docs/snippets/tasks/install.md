
You can install the lifecycle toolkit using the current release manifest:
<!---x-release-please-start-version-->
```
kubectl apply -f https://github.com/keptn/lifecycle-toolkit/releases/download/v0.5.0/manifest.yaml
kubectl wait --for=condition=Available deployment/klc-controller-manager -n keptn-lifecycle-toolkit-system --timeout=120s
```
<!---x-release-please-end-->

Now, the Lifecycle Toolkit and its dependency is installed and ready to use.