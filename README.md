# kullect
Kapacitor UDF's to Analyze Service Costs in Kubernetes

### Setup

1. Download binary from releases

2. Add to Kapacitor Config File
  ```
  [udf]
    [udf.functions]
    [udf.functions.kullect]
            prog = "path/to/kullect/binary"
            args = []
            timeout = "10s"
  ```

3. Add the tick file included here
  * Change the values for **total_cpu**, **hourly_cost**, and **total_memory**

4. Define and enable the tick file
  ```
  kapacitor define kullector \
    -tick kullect.tick \
    -type stream \
    -dbrp k8s.default 
  kapacitor enable kullector
  ```

5. Validate the tick is working
  ```
  kapacitor show kullector
  ```
