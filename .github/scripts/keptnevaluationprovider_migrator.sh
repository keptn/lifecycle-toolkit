#!/bin/bash

getNamespaces() {
  out=$(kubectl get ns -o name)
  for n in $out
  do
    nn=${n#"namespace/"}
    namespaces+=("$nn")
  done
}

getProvidersOfNamespace() {
  out=$(kubectl get keptnevaluationproviders -o name -n "$1")
  for n in $out
  do
    nn=${n#"keptnevaluationprovider.lifecycle.keptn.sh/"}
    providers+=("$nn -n $1")
  done
}

echo -e "------------------------------\n"
echo -e "Migrating manifests.\n"
echo -e "------------------------------\n"

declare -a namespaces
declare -a providers

DATE=$(date +%s)
MANIFESTS_FILE="manifests-$DATE.yaml"

getNamespaces

for n in "${namespaces[@]}"
do
  getProvidersOfNamespace "$n"
done

for n in "${providers[@]}"
do
  echo "---" >> $MANIFESTS_FILE
  kubectl get keptnevaluationproviders $n -oyaml >> $MANIFESTS_FILE
done

sed -i 's#lifecycle\.keptn\.sh/v1alpha2#metrics\.keptn\.sh/v1alpha2#g' $MANIFESTS_FILE
sed -i 's/KeptnEvaluationProvider/KeptnMetricsProvider/g' $MANIFESTS_FILE

echo -e "------------------------------\n"
echo -e "Manifests migrated successfully.\n"
echo -e "To review the newly created manifests, check ./$MANIFESTS_FILE file.\n"
echo -e "------------------------------\n"

read -p "Do you want to apply the newly created KeptnMetricsProvider resources? [y/N]" -n 1 -r
echo -e "\n"
if [[ $REPLY =~ ^[Yy]$ ]]
then
  kubectl apply -f $MANIFESTS_FILE
  echo -e "\nManifests applied.\n"
else
  echo -e "Manifests not applied.\n"
fi

echo -e "------------------------------\n"

read -p "Do you want to delete the old KeptnEvaluationProvider resources? [y/N]" -n 1 -r
echo -e "\n"
if [[ $REPLY =~ ^[Yy]$ ]]
then
  for n in "${providers[@]}"
  do
    kubectl delete keptnevaluationproviders $n
  done
  echo -e "\nResources deleted.\n"
else
  echo -e "Resources not deleted.\n"
fi
