# Grafana Dashboards for the Lifecycle Toolkit

This folder contains the Grafana dashboards for the Keptn.

## Installing the dashboards

It is assumed, that there is a Grafana Instance available.
In our provided examples, the dashboards are automatically
provisioned.
If you want to install the dashboards manually, you can use the following steps:

```shell
# This defaults to http://localhost:3000, but can be changed by setting the GRAFANA_SCHEME, GRAFANA_URL and GRAFANA_PORT environment variable
# The default credentials are admin:admin, but can be changed by setting the GRAFANA_USERNAME and GRAFANA_PASSWORD environment variable
make install
```

## Changing the dashboards

The dashboards can be changed in the Grafana UI.
To export dashboards, export them using the share button and replace
them in this folder.

## Exporting the dashboards for the Examples

You can import the default dashboards by running:

```shell
make apply-configmaps
```

If you prefer to prepare the dashboards for the examples and importing them as json you can use:

```shell
make import-json
```

---

## Guide to Custom Dashboards with Grafana

This section provides links to the official Grafana documentation for creating and modifying custom dashboards using Grafana.

## How to Create Custom Dashboards

Craft personalized Grafana dashboards tailored to your specific needs by exploring the [official Grafana documentation](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/create-dashboard/).
Discover detailed steps, best practices, and tips for creating visualizations, adding panels,
utilizing various data sources, and more.

## Modifying the Dashboards

Refine and adapt visualizations in existing Grafana dashboards to meet evolving requirements.
Visit the
[relevant sections](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/modify-dashboard-settings/)
of the Grafana documentation to gain insights on adjusting, enhancing, or reconfiguring dashboards efficiently.
Learn about editing panels, incorporating new data sources, applying filters, and optimizing dashboard layouts.
