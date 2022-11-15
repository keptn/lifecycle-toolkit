#!/bin/bash

OUTDIR="$1"

for i in $(ls grafana*); do
  cat << EOF | jq '.dashboard.id = null | .dashboard.uid = null' > "$OUTDIR/$i"
  {
  "dashboard":
  $(cat "$i")
  }
EOF
done