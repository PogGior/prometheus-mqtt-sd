package util_test

import (
	"testing"

	"prometheus-mqtt-sd/cmd/prometheus-mqtt-sd/util"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/stretchr/testify/assert"
)

func TestProcessMessage(t *testing.T) {
	payload := []byte(`[
		{
			"targets": ["localhost:9090", "localhost:9091"],
			"labels": {
				"job": "info",
				"instance": "instance-1"
			},
			"source": "info"
		},
		{
			"targets": ["localhost:8080"],
			"labels": {
				"job": "debug",
				"instance": "instance-2"
			},
			"source": "debug"
		}
	]`)

	expectedTargetGroups := []*targetgroup.Group{
		{
			Targets: []model.LabelSet{
				{
					model.AddressLabel: model.LabelValue("localhost:9090"),
				},
				{
					model.AddressLabel: model.LabelValue("localhost:9091"),
				},
			},
			Labels: model.LabelSet{
				model.LabelName("job"):      model.LabelValue("info"),
				model.LabelName("instance"): model.LabelValue("instance-1"),
			},
			Source: "info",
		},
		{
			Targets: []model.LabelSet{
				{
					model.AddressLabel: model.LabelValue("localhost:8080"),
				},
			},
			Labels: model.LabelSet{
				model.LabelName("job"):      model.LabelValue("debug"),
				model.LabelName("instance"): model.LabelValue("instance-2"),
			},
			Source: "debug",
		},
	}

	targetGroups, err := util.ProcessMessage(payload)
	assert.Nil(t, err)
	assert.Equal(t, expectedTargetGroups, targetGroups)
}
