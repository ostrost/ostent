package internal_models

import (
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/filter"
)

// Filter containing drop/pass and tagdrop/tagpass rules
type Filter struct {
	NameDrop []string
	nameDrop filter.Filter

	IsActive bool
}

// Compile all Filter lists into filter.Filter objects.
func (f *Filter) CompileFilter() error {
	var err error
	f.nameDrop, err = filter.CompileFilter(f.NameDrop)
	if err != nil {
		return fmt.Errorf("Error compiling 'namedrop', %s", err)
	}
	return nil
}

func (f *Filter) ShouldMetricPass(metric telegraf.Metric) bool {
	// TODO if f.ShouldNamePass(metric.Name()) && f.ShouldTagsPass(metric.Tags()) {
	if f.ShouldNamePass(metric.Name()) {
		return true
	}
	return false
}

// ShouldFieldsPass returns true if the metric should pass, false if should drop
// based on the drop/pass filter parameters
func (f *Filter) ShouldNamePass(key string) bool {
	/* TODO
	if f.namePass != nil {
		if f.namePass.Match(key) {
			return true
		}
		return false
	}
	*/

	if f.nameDrop != nil {
		if f.nameDrop.Match(key) {
			return false
		}
	}
	return true
}
