# prominv

## Overview
prominv is an [Ansible Dynamic Inventory](https://docs.ansible.com/projects/ansible/latest/inventory_guide/intro_dynamic_inventory.html) that uses the [Prometheus API](https://prometheus.io/docs/prometheus/latest/querying/api/) as a data source.  The inventory use the labels associated with a given metric to construct variables that can be accessed from Ansible.  It also creates [Child Groups](https://docs.ansible.com/projects/ansible/latest/inventory_guide/intro_inventory.html#grouping-groups-parent-child-group-relationships) within the Inventory to group hosts according to label names and content.

## Getting Started
Start by cloning this Git repository.  Once you have the code, it can be installed using **go build**.  Providing no issues are encountered, you should now be able to run *prominv*.  Without any configuration, it won't do anything beyond completeing without error if the build is valid.  Read on to create a valid configuration.

## Configuration
The prominv program will try to guess at required configuration options that have not been explicitly set.  It's probably best however to construct a fully-populated configuration file.  The configuration file can then be specified using the *prominv -config <location>* flag.  Below is a complete example of a configuration file (in yaml format).

    ---
    prometheus_url: http://prometheus.mydomain.org:9090
    promql_query: up{job="node"}
    labels:
      group_by:
        - env
        - site

The options in the above example are described below:-
* prometheus_url: The URL of the Prometheus instance to be queried.
* promql_query: The query, in PromQL format, used to generate the desired inventory.
* labels: A category of things relating to Prometheus labels.
  * group_by: A list of labels return by the promql_query that should be used to generate inventory Child Groups.