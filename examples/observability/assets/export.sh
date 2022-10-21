#!/usr/bin/env bash

# SOURCE VARIABLES
export grafanaurl=http://admin:admin@localhost:3000/api


rm -rf backups
mkdir -p backups
cd backups

datasources=$(curl -s -X GET $grafanaurl/datasources)
for uid in $(echo $datasources | jq -r '.[] | .uid'); do
  uid="${uid/$'\r'/}" # remove the trailing '/r'
  curl -s -X GET "$grafanaurl/datasources/uid/$uid" | jq > grafana-datasource-$uid.json
  slug=$(cat grafana-datasource-$uid.json | jq -r '.name')
  mv grafana-datasource-$uid.json grafana-datasource-$uid-$slug.json # rename with datasource name and id
  echo "Datasource $uid exported"
done



dashboards=$(curl -s -X GET $grafanaurl/search?folderIds=0&query=&starred=false)
for uid in $(echo $dashboards | jq -r '.[] | .uid'); do
  uid="${uid/$'\r'/}" # remove the trailing '/r'
  curl -s -X GET "$grafanaurl/dashboards/uid/$uid" | jq > grafana-dashboard-$uid.json
  slug=$(cat grafana-dashboard-$uid.json | jq -r '.meta.slug')
  mv grafana-dashboard-$uid.json grafana-dashboard-$uid-$slug.json # rename with dashboard name and id
  echo "Dashboard $uid exported"
done
