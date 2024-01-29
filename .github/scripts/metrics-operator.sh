#!/bin/bash

# Keptn Helm Testing
#
# This script supports the comparison of standard values and expected templated results to helm chart
# it is used to make sure changes to the chart are intentional and produce expected outcomes

echo "copying manifests"

dir="metrics-operator/config/crd/bases"
helm_dir="metrics-operator/chart/templates"

configElements=(analysisdefinition analysisvaluetemplate keptnmetricsprovider)

truncate -s 0 $helm_dir/keptnmetric-crd.yaml
cat $dir/metrics.keptn.sh_keptnmetrics.yaml >> $helm_dir/keptnmetric-crd.yaml

truncate -s 0 $helm_dir/analysis-crd.yaml
cat $dir/metrics.keptn.sh_analyses.yaml >> $helm_dir/analysis-crd.yaml

n=3
# Loop through each file in the directory
for file in "$dir"/* ; do
    # Extract the basename of the file
    configCrds=$(basename "$file" .yaml)
    echo "Processing file: $filename"

    # Loop through each element in the pickElement array
    for element in "${configElements[@]}"; do
        echo "Checking element: $element"

        # Check if the element is present in the filename
        if echo "$configCrds" | grep -q "$element" && [[ $n > -1 ]] ; then
            echo "Match found: $element in $filename"

            # Loop through files in the helm directory
            for crds in "$helm_dir"/* ; do
                # Extract the basename of the helm file
                helm_filename=$(basename "$crds" .yaml)

                echo "Checking helm file: $helm_filename"

                # Check if the element is present in the helm filename
                if echo "$helm_filename" | grep -q "$element" && [[ $n > -1 ]] ; then

                    echo "Match found: $element in $helm_filename"
                      truncate -s 0 "$crds"
                      ((n=n-1))
                    # Concatenate the content of the file into the helm file
                    cat "$file" >> "$crds"
                    break
                fi
            done
        fi
    done
done

for file in "$helm_dir"/*; do
    filename=$(basename "$file" .yaml)
    if [[ $filename == k* || $filename == a* ]]; then
                sed -i '/controller-gen.kubebuilder.io\/version: v0.14.0/a\
    {{- with .Values.global.caInjectionAnnotations }}\
    {{- toYaml . | nindent 4 }}\
    {{- end }}\
    {{- include "common.annotations" ( dict "context" . ) }}\
  labels:\
    app.kubernetes.io/part-of: keptn\
    crdGroup: metrics.keptn.sh\
    keptn.sh/inject-cert: "true"\
  {{- include "common.labels.standard" ( dict "context" . ) | nindent 4 }}' "$helm_dir/$filename.yaml"
        fi
        if [[ $filename == "analysis-crd" ]] ; then
        echo "found"
                sed -i "/{{- end }}/a\\
    cert-manager.io/inject-ca-from: '{{ .Release.Namespace }}\/keptn-certs'" "$helm_dir/$filename.yaml"
        fi
done