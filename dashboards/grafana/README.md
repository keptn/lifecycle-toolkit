# Grafana Dashboards for the Lifecycle Toolkit

This folder contains the Grafana dashboards for the Keptn Lifecycle Toolkit.

## Installing the dashboards
It is assumed, that there is a Grafana Instance available. In our provided examples, the dashboards are automatically provisioned. If you want to install the dashboards manually, you can use the following steps:

```sh
# This defaults to http://localhost:3000, but can be changed by setting the GRAFANA_SCHEME, GRAFANA_URL and GRAFANA_PORT environment variable
# The default credentials are admin:admin, but can be changed by setting the GRAFANA_USERNAME and GRAFANA_PASSWORD environment variable
make install
```

## Changing the dashboards
The dashboards can be changed in the Grafana UI. To export dashboards, export them using the share button and replace them in this folder.

## Exporting the dashboards for the Examples
You can prepare the dashboards for the examples and import using the following command:

```sh
make generate
```
