package agent

import (
	"log"
	"math"
	"time"

	"github.com/influxdata/telegraf"
)

func NewAccumulator(
	metrics chan telegraf.Metric,
) *accumulator {
	acc := accumulator{}
	acc.metrics = metrics
	return &acc
}

type accumulator struct {
	metrics chan telegraf.Metric

	debug bool
}

// Add of telegraf.Accumulator interface.
func (ac *accumulator) Add(measurement string, value interface{}, tags map[string]string, t ...time.Time) {
	ac.AddFields(measurement, map[string]interface{}{"value": value}, tags, t...)
}

// AddFields of telegraf.Accumulator interface.
func (ac *accumulator) AddFields(measurement string, fields map[string]interface{}, tags map[string]string, t ...time.Time) {
	if measurement == "" || len(fields) == 0 {
		return
	}
	if tags == nil {
		tags = make(map[string]string)
	}
	values := make(map[string]interface{})
	for k, v := range fields {

		// Validate uint64 and float64 fields
		switch val := v.(type) {
		case uint64:
			// InfluxDB does not support writing uint64
			if val < uint64(9223372036854775808) {
				values[k] = int64(val)
			} else {
				values[k] = int64(9223372036854775807)
			}
			continue
		case float64:
			// NaNs are invalid values in influxdb, skip measurement
			if math.IsNaN(val) || math.IsInf(val, 0) {
				if false { // if ac.debug TODO
					log.Printf(
						"Measurement [%s] field [%s] has a NaN or Inf field, skipping",
						measurement, k)
				}
				continue
			}
		}

		values[k] = v
	}
	if len(values) == 0 {
		return
	}

	var ts time.Time
	if len(t) > 0 {
		ts = t[0]
	} else {
		ts = time.Now()
	}
	// timestamp = timestamp.Round(ac.precision) // TODO

	m, err := telegraf.NewMetric(measurement, tags, values, ts)
	if err != nil {
		log.Printf("Error adding point [%s]: %s\n", measurement, err.Error())
		return
	}
	// if ac.trace { fmt.Println("> " + m.String()) } // TODO
	ac.metrics <- m
}

func (ac accumulator) Debug() bool {
	return ac.debug
}

func (ac *accumulator) SetDebug(on bool) {
	ac.debug = on
}

func (ac *accumulator) SetPrecision(precision, interval time.Duration) {
	// TODO
}

func (ac *accumulator) DisablePrecision() {
	// TODO
}

// TODO AddError is not part of telegraf.Accumulator interface as of 1.0.0-beta3
func (ac *accumulator) AddError(err error) { log.Printf("Error in input: %s", err) }
