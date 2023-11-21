# Setting up and running prometheus

## Steps being followed to set up and run prometheus in local

### Cluster config path
  - `export KUBECONFIG="/Users/vijaysek/observability/dopemeth/clusterConf/devcluster_config.yaml"`
### Running kind
  - It is a tool used to create a local kubenetes cluster using docker (need docker daemon to be running)
  - It makes a kubenetes docker container 
  - installation
    - `brew install kind`
  - create a k8s cluster
    - `kind create cluster`
  - Delete a cluster after work done (unless killed the docker image will be running)
    - `kind delete cluster`
### Setting up prometheus in the cluster
  - Installing prometheus in the cluster
    - We use helm to install prometheus in the cluster
    - https://artifacthub.io/packages/helm/prometheus-community/prometheus --> For more info
    - The following commands can be used to install prometheus in the cluster
    - `helm repo add prometheus-community https://prometheus-community.github.io/helm-charts`
    - `helm install my-prometheus prometheus-community/prometheus --version 25.4.0`
    - `helm upgrade --install my-prometheus prometheus-community/prometheus --set prometheus.prometheusSpec.configMaps.prometheus=custom-prometheus-config`
## Accessing prometheus metrics from the cluster
  - installation of prometheus will spawn up some pods and services that can be used to assess the metrics
    - Portforwarded from a service named 'my-prometheus-server '
    - `kubectl port-forward services/my-prometheus-server 8080:80`
  - This way the endpoint of prometheus
    - `http:localhost:8080/metrics`
    - `http:localhost:8080`
## Setting up Grafana
  - installation of grafana in cluster
    - `helm repo add grafana https://grafana.github.io/helm-charts`
    - `helm install my-grafana grafana/grafana --version 7.0.3`
  - website followed
    - https://artifacthub.io/packages/helm/grafana/grafana
  - After setting up of grafana we need to get secrets that can be used to access the grafana ui
    - Secret can be obtained from the following command
    - `kubectl get secret --namespace default my-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo`
    - login creds
      - username : admin
      - password : qyELVamWKl16nGPEcPXPRk4W3K6ni7GYeZjVeVMo
    - Data will be lost when the grafana pod is terminated (no idea why is here but worth noting)
## Setting up of influx db
  - installation of influx db in cluster
    - `helm repo add bitnami https://charts.bitnami.com/bitnami`
    - `helm install my-influxdb bitnami/influxdb --version 5.10.0`
  - Website followed
    - https://artifacthub.io/packages/helm/bitnami/influxdb
  - connect to influx db outside cluster
    - `kubectl port-forward --namespace default svc/my-influxdb 8086:8086 & influx -host 127.0.0.1 -port 8086`

## This will make a docker image running from what is pulled form prometheus
docker run -d -p 9090:9090 \
  -v /Users/vijaysek/observability/dopemeth/config/prometheus.yaml:/etc/prometheus/prometheus.yaml \
  --name prometheus \
  prom/prometheus

## docker build command from scratch with logs
  - `docker build --platform linux/amd64 -t containers.cisco.com/vijaysek/server-app1:v8 --no-cache --progress plain .`

## docker build 
  - `docker build --platform linux/amd64 -t containers.cisco.com/vijaysek/server-app1:v8 --progress plain .`


## Setup vscode in default for vi when editing something in the terminal
  - `export EDITOR='code --wait'`

## Editing a kuebctl config file 
  - `k edit <resource-type> <resource-name>`

## updating the changes in the configmap in the deployment 
  - `k rollout restart <resourcetype> <reource-name>`

## Settingup opentelemetry collector
  - `kubectl apply -f https://raw.githubusercontent.com/open-telemetry/opentelemetry-collector/main/examples/k8s/otel-config.yaml`

## Loading local config file into otel container
  - `docker run -v $(pwd)/config.yaml:/etc/otelcol-contrib/config.yaml otel/opentelemetry-collector-contrib:0.88.0`

## Fun kubernetes command to play around
  -  `kubectl get configmap my-prometheus-server -o go-template='{{.data.prometheus}}`
  
## Webex access token

## Setting up alert manager
  - First prometheus alerting rules are set in the prometheus.yaml file
  - Alertmanager config files is configured with appropriate route to send trigger
  
## Sample curl to test the testfire alert
  - `curl -XPOST http://[localhost:[port where prometheus ui running]]/api/v1/query -d 'query=vector(2)'`

  webexteams://im?space=3141ae60-7d39-11ee-90ad-45d6ddfdea08

  ## app password for google gmail alert channel
   - lstj tyma nksz ulbs

  ## Team creation webex 
    
    ```
    {
      "id": "Y2lzY29zcGFyazovL3VzL1RFQU0vMGExYjZiMDAtODc3MC0xMWVlLTgwNzItZTcxZTdlMzVhNjc0",
      "name": "monitoringTeam",
      "description": "For receiving notification from alertmanager",
      "creatorId": "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9jNjFjNjJhZC1mYjFiLTQ5OTItYmE5Mi0xZjJiNWQyYmE5NGY",
      "created": "2023-11-20T06:42:57.072Z"
    }
    ```

## webex room 
    ```
        {
      "id": "Y2lzY29zcGFyazovL3VzL1JPT00vY2UyMjVlZjAtODc3MC0xMWVlLWEzYjEtZmI3ZmU5MWQ3NmZk",
      "title": "prometh-alertmanager",
      "type": "group",
      "isLocked": false,
      "lastActivity": "2023-11-20T06:48:25.961Z",
      "teamId": "Y2lzY29zcGFyazovL3VzL1RFQU0vMGExYjZiMDAtODc3MC0xMWVlLTgwNzItZTcxZTdlMzVhNjc0",
      "creatorId": "Y2lzY29zcGFyazovL3VzL1BFT1BMRS9jNjFjNjJhZC1mYjFiLTQ5OTItYmE5Mi0xZjJiNWQyYmE5NGY",
      "created": "2023-11-20T06:48:25.961Z",
      "ownerId": "Y2lzY29zcGFyazovL3VzL09SR0FOSVpBVElPTi8xZWI2NWZkZi05NjQzLTQxN2YtOTk3NC1hZDcyY2FlMGUxMGY",
      "description": "room for getting notifications",
      "isPublic": false
    }
    ```