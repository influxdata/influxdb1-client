# influxdb1-clientv2
influxdb1-clientv2 is the current Go client API for InfluxDB 1.x. For connecting to InfluxDB 2.x see the [influxdb-client-go](https://github.com/influxdata/influxdb-client-go) client library.

InfluxDB is an open-source distributed time series database, find more about [InfluxDB](https://www.influxdata.com/time-series-platform/influxdb/) at https://docs.influxdata.com/influxdb/latest

## Usage
To import into your Go project, run the following command in your terminal:
`go get github.com/influxdata/influxdb1-client/v2`
Then, in your import declaration section of your Go file, paste the following:
`import "github.com/influxdata/influxdb1-client/v2"`

If you get the error `build github.com/user/influx: cannot find module for path github.com/influxdata/influxdb1-client/v2` when trying to build:
change your import to:
```go
import(
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)
```
## Example
The following example creates a new client to the InfluxDB host on localhost:8086 and runs a query for the measuring the average value of `cpu_load` from the `mydb` database over the last 10 timesteps. 
``` go
func ExampleClient_query() {
	var (
		numTimesteps = 10
		cpu_load float64
	)

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()
	// build the query
	query := fmt.Sprintf("SELECT value FROM cpu_load limit %d", numTimesteps)
	q := client.NewQuery(query, "mydb", "")
	response, err := c.Query(q)
	if err != nil {
		fmt.Println("Error while querying the database: ", err.Error())
	}
	rows := response.Results[0].Series[0]
	for j := 0; j < numTimesteps; j++ {
		// while the first column(0) is the timestamp,
		// we get the second column(1)
		val, err := rows.Values[j][1].(json.Number).Float64()
		if err != nil {
			fmt.Printf("Error while calculating %s", rows.Columns[1])
			return
		}
		cpu_load += val
	}
	fmt.Printf("The average value of the last %d measurements is %f\n", numTimesteps, cpu_load/float64(numTimesteps))
}
```

