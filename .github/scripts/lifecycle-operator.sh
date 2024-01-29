#!/bin/bash

# Keptn script for endtasks
#
# This script copies the crds from lifecycle-operator(lifecycle-operator/config/crd/bases) to the directory's helm chart(lifecycle-operator/chart/templates)
# It also appends the respective annotations

echo "copying manifests"

dir="lifecycle-operator/config/crd/bases"
helm_dir="lifecycle-operator/chart/templates"

configElements=(keptnapps keptnappcontext keptnappcreationrequest keptnappversion keptnconfig keptnevaluations keptnevaluationdefinition keptntasks keptntaskdefinition keptnworkloads keptnworkloadversion)

n=10
for file in "$dir"/* ; do
    # Extract the basename of the file
    configCrds=$(basename "$file" .yaml)

    # Loop through each element in the pickElement array
    for element in "${configElements[@]}"; do

        # Check if the element is present in the filename
        if  echo "$configCrds" | grep -q "$element" && [[ $n > -1 ]] ; then
            echo "Match found: $element in $configCrds"

            # Loop through files in the helm directory
            for crds in "$helm_dir"/* ; do
                # Extract the basename of the helm file
                helm_filename=$(basename "$crds" .yaml)

                # Check if the element is present in the helm filename
                if [[ $helm_filename == k* ]] && echo "$helm_filename" | grep -q "$element" && [[ $n > -1 ]] ; then

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
        if [[ $filename == k* ]] ; then
                sed -i '/controller-gen.kubebuilder.io\/version: v0.14.0/a\
    {{- with .Values.global.caInjectionAnnotations  }}\
    {{- toYaml . | nindent 4 }}\
    {{- end }}\
  {{- include "common.annotations" ( dict "context" . ) }}\
  labels:\
    app.kubernetes.io/part-of: keptn\
    crdGroup: lifecycle.keptn.sh\
    keptn.sh/inject-cert: "true"\
{{- include "common.labels.standard" ( dict "context" . ) | nindent 4 }}' "$helm_dir/$filename.yaml"
        fi
            if [[ "$filename" == "keptnappcontext-crd" ]] ; then
            sed -i "s|{{- with .Values.global.caInjectionAnnotations  }}|cert-manager.io/inject-ca-from: '{{ .Release.Namespace }}/keptn-certs'|g" "$helm_dir/$filename.yaml"
            sed -i "/{{- end }}/d" "$helm_dir/$filename.yaml"
            sed -i "/{{- toYaml \. /d" "$helm_dir/$filename.yaml"
            fi
            if  [[ "$filename" == "keptnapps-crd" || "$filename" == "keptnappversion-crd" ]] ; then
           sed -i "0,/spec/{/spec/a\\
  conversion:\\
    strategy: Webhook\\
    webhook:\\
      clientConfig:\\
        service:\\
          name: 'lifecycle-webhook-service'\\
          namespace: '{{ .Release.Namespace }}'\\
          path: \/convert\\
      conversionReviewVersions:\\
      - v1
           }" "$helm_dir/$filename.yaml"
            fi
done