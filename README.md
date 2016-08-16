# kullect
Kapacitor UDF's to Analyze Service Costs in Kubernetes

Add to Kapacitor Config File
```
[udf]
  [udf.functions]
  [udf.functions.kullect]
            prog = "./kullect"
            args = []
            timeout = "10s"
```
Create a Kapacitor Tick File
```
var uptime = stream
    |from()
        .measurement('uptime')
        .groupBy('pod_name','namespace_name','labels')

var usage = stream
    |from()
        .measurement('cpu/usage_rate')
        .groupBy('pod_name','namespace_name','labels')
usage
    |join(uptime)
        .as('cpu','uptime')

    @kullect()
        .as('value')
    |influxDBOut()
        .database('k8s')
        .retentionPolicy('default')
        .measurement('cpu/cost')
        .tag('kapacitor', 'true')
        .tag('version', '0.2')
```
