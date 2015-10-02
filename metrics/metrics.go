package metrics

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// DataPoint is a single json value.
type DataPoint struct {
	Metric      string      `json:"metric"`
	Value       interface{} `json:"value"`
	Description string      `json:"description"`
}

// MarshalJSON verifies d is valid and converts it to JSON.
func (d *DataPoint) String() ([]byte, error) {
	if err := d.clean(); err != nil {
		return nil, err
	}
	return json.Marshal(struct {
		Metric string      `json:"metric"`
		Value  interface{} `json:"value"`
	}{
		d.Metric,
		d.Value,
	})
}

func (d *DataPoint) clean() error {

	switch v := d.Value.(type) {
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			d.Value = i
		} else if f, err := strconv.ParseFloat(v, 64); err == nil {
			d.Value = f
		} else {
			return fmt.Errorf("Unparseable number %v", v)
		}
	case uint64:
		if v > math.MaxInt64 {
			d.Value = float64(v)
		}

	}
	return nil
}

// MultiDataPoint is a list with data points
type MultiDataPoint []DataPoint

// MarshalJSON verifies d is valid and converts it to JSON.
func (md *MultiDataPoint) String() ([]byte, error) {
	b, err := json.Marshal(md)
	return b, err

}

// Group is a group of multiple data points.
type Group struct {
	Metrics MultiDataPoint `json:"metrics"`
	Name    string         `json:"name"`
}

// MarshalJSON verifies d is valid and converts it to JSON.
func (g *Group) String() ([]byte, error) {
	fmt.Println(g)
	return json.Marshal(g)
}

// Block is a group of multiple groups.
type Block []Group

// MarshalJSON verifies d is valid and converts it to JSON.
func (b *Block) String() ([]byte, error) {
	return json.Marshal(b)
}

// Add creates a new datapoint
func (md *MultiDataPoint) Add(metric string, value interface{}) {
	e := DataPoint{
		Metric: metric,
		Value:  value,
	}
	e.clean()

	*md = append(*md, e)
}
