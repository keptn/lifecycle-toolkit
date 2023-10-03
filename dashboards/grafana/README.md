# Grafana Dashboards for the Lifecycle Toolkit

This folder contains the Grafana dashboards for the Keptn Lifecycle Toolkit.

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
To export dashboards, follow these steps:

- Access your Grafana instance.
- Open the dashboard you want to modify.
- Click on the "Share" button located in the top toolbar.
- From the dropdown menu, select "Export" to export the dashboard as a JSON file.
- Save the exported JSON file.
- Replace the corresponding dashboard file in this folder with your modified JSON file.

## Creating custom dashboards

To create your own custom dashboards in Grafana,
follow these steps:

- Log in to your Grafana instance.
- Click on the "+" icon in the left-hand menu and select "Dashboard" to create a new dashboard.
- On the new dashboard page, you'll see a toolbar at the top with various options.
- Click on the "Add panel" button to add panels to your dashboard.
  Panels are the individual visualizations or components that display data.
- Choose the type of panel you want to add, such as Graph, Singlestat, Table, etc.
  Each panel type has its own configuration options.
- Configure the panel by selecting the data source, specifying the query, and customizing the visualization
  settings as per your requirements.
- Repeat the above steps to add more panels to your dashboard.
- Customize the layout of your dashboard by dragging and resizing panels.
- Use the toolbar options to further customize your dashboard.
  You can add text, annotations, variables, and apply different themes to your dashboard.
- Once you've created your custom dashboard, click on the "Save" icon in the toolbar to save it.

## Exporting the dashboards for the Examples

You can import the default dashboards by running:

```shell
make apply-configmaps
```

If you prefere to prepare the dashboards for the examples and importing them as json you can use:

```shell
make import-json
```

<!-- markdownlint-disable-next-line MD033 MD013 -->
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />
