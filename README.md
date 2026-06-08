# prominv

## Overview
prominv is an [Ansible Dynamic Inventory](https://docs.ansible.com/projects/ansible/latest/inventory_guide/intro_dynamic_inventory.html) that uses the [Prometheus API](https://prometheus.io/docs/prometheus/latest/querying/api/) as a data source.  The inventory use the labels associated with a given metric to construct variables that can be accessed from Ansible.  It also creates [Child Groups](https://docs.ansible.com/projects/ansible/latest/inventory_guide/intro_inventory.html#grouping-groups-parent-child-group-relationships) within the Inventory to group hosts according to label names and content.

## Getting Started

## Configuration
The prominv program will try to guess at required configuration options that have not been explicitly set.  It's probably best however to construct a fully-populated configuration file.  The configuration file can then be specified using the -config <location> flag.