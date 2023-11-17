## This explains the flow of the code

## File structure
    - main code resides in /main/cmd.go
        - This is where the server and prometheus metrics are instanciated
    - helpers /services/stockApi.go
        - This program holds helperfunctions
            - function to ping to stock api provider and get prices
            - function to send http response
            - function of updating stock counter and gauge
    - poller/poller.go
        - This package holds a poller that pings to the api endpoint every 1m to get gain loss metrics and update them
    - /k8s
        - this folder holds yaml for creating a deployment for the server
    - prometh/current_config
        - This file contains current config in prometheus configmap and prometheus alertmanager configmap
        - This was prewritten and i have edited it to add
            - alerting rules in prometheus.yaml
            - alerting route in promethus_alertmanager.yaml
    

## How is this code working
    - Written a simple server and poller containerized it
    - Made a deployment in name of stock-server that runns the container
    - added annotations in the deployment for prometheus to automatically find and scrape for it
    - added rules file in the `my-prometheus-server` configmap
    - added route to the `my-prometheus-alertmanager` configmap

## Status of program
    - pods running properly
    - prometheus can identify and scrape for the metrics
    - alermanager alers are found and firing

## Work left
    - The alerts are firing but needs to be tuned to receive it in webex
    