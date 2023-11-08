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