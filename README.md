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
stream
    // Select just the cpu measurement from our example database.
    |from()
        .measurement('cpu/usage')
    @kullect()
        .field('value')
        .uptime('uptime.value')
        .size(100)
        .as('cost')
    |influxDBOut()
        .database('k8s')
        .retentionPolicy('default')
        .measurement('cpu/usage')
        .tag('kapacitor', 'true')
        .tag('version', '0.2')
```
