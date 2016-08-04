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
stream
    // Select just the cpu measurement from our k8s database.
    |from()
        .measurement('cpu/usage')
    @kullect()
        .field('value')
        .size(100)
        .as('cost')
    |httpOut('cpu_cost')
    |influxDBOut()
        .database('k8s')
        .retentionPolicy('default')
        .measurement('cpu/usage')
        .tag('kapacitor', 'true')
```
