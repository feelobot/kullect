package main

import (
	"errors"
	"github.com/influxdata/kapacitor/udf"
	"github.com/influxdata/kapacitor/udf/agent"
	"log"
	"os"
	"strconv"
)

// An Agent.Handler
type costHandler struct {
	field  string
	uptime string
	as     string
	size   int
	agent  *agent.Agent
}

func newCostHandler(a *agent.Agent) *costHandler {
	return &costHandler{agent: a}
}

// Return the InfoResponse. Describing the properties of this UDF agent.
func (a *costHandler) Info() (*udf.InfoResponse, error) {
	info := &udf.InfoResponse{
		Wants:    udf.EdgeType_STREAM,
		Provides: udf.EdgeType_STREAM,
		Options: map[string]*udf.OptionInfo{
			"field":  {ValueTypes: []udf.ValueType{udf.ValueType_STRING}},
			"uptime": {ValueTypes: []udf.ValueType{udf.ValueType_STRING}},
			"size":   {ValueTypes: []udf.ValueType{udf.ValueType_INT}},
			"as":     {ValueTypes: []udf.ValueType{udf.ValueType_STRING}},
		},
	}
	return info, nil
}

// Initialze the handler based of the provided options.
func (a *costHandler) Init(r *udf.InitRequest) (*udf.InitResponse, error) {
	init := &udf.InitResponse{
		Success: true,
		Error:   "",
	}
	for _, opt := range r.Options {
		switch opt.Name {
		case "field":
			a.field = opt.Values[0].Value.(*udf.OptionValue_StringValue).StringValue
		case "uptime":
			a.uptime = opt.Values[0].Value.(*udf.OptionValue_StringValue).StringValue
		case "size":
			a.size = int(opt.Values[0].Value.(*udf.OptionValue_IntValue).IntValue)
		case "as":
			a.as = opt.Values[0].Value.(*udf.OptionValue_StringValue).StringValue
		}
	}

	if a.field == "" {
		init.Success = false
		init.Error += " must supply field"
	}
	if a.uptime == "" {
		init.Success = false
		init.Error += " must supply field uptime"
	}
	if a.size == 0 {
		init.Success = false
		init.Error += " must supply window size"
	}
	if a.as == "" {
		init.Success = false
		init.Error += " invalid as name provided"
	}

	return init, nil
}

// Create a snapshot of the running state of the process.
func (a *costHandler) Snaphost() (*udf.SnapshotResponse, error) {
	return &udf.SnapshotResponse{}, nil
}

// Restore a previous snapshot.
func (a *costHandler) Restore(req *udf.RestoreRequest) (*udf.RestoreResponse, error) {
	return &udf.RestoreResponse{
		Success: true,
	}, nil
}

// This handler does not do batching
func (a *costHandler) BeginBatch(*udf.BeginBatch) error {
	return errors.New("batching not supported")
}

// Compute Cost
func (a *costHandler) Point(p *udf.Point) error {
	// Update the moving average.
	// Re-use the existing point so we keep the same tags etc.
	total_millicores := 544000.00
	hourly_rate := 33.46
	cpu_usage := p.FieldsDouble[a.field]
	uptime, _ := strconv.ParseFloat(p.FieldsString["uptime"], 64)
	hourly_uptime := uptime / 36000000
	cost := (hourly_uptime * hourly_rate) * (cpu_usage / total_millicores)
	log.Println("CPU: ", cpu_usage)
	p.FieldsDouble = map[string]float64{a.as: cost}
	p.FieldsInt = nil
	p.FieldsString = nil
	// Send point with average value.
	a.agent.Responses <- &udf.Response{
		Message: &udf.Response_Point{
			Point: p,
		},
	}
	return nil
}

// This handler does not do batching
func (a *costHandler) EndBatch(*udf.EndBatch) error {
	return errors.New("batching not supported")
}

// Stop the handler gracefully.
func (a *costHandler) Stop() {
	close(a.agent.Responses)
}

func main() {
	a := agent.New(os.Stdin, os.Stdout)
	h := newCostHandler(a)
	a.Handler = h

	log.Println("Starting agent")
	a.Start()
	err := a.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
