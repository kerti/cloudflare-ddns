package worker

import (
	"testing"

	"github.com/kerti/cloudflare-ddns/config"
	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/kerti/cloudflare-ddns/resolver"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	initloglevel uint8 = uint8(0)

	initIntervalTestCases = []struct {
		name          string
		configValue   string
		resolverCount int
		intervalValue int
	}{
		{
			name:          "emptyConfig",
			configValue:   "",
			resolverCount: 6,
			intervalValue: 50,
		},
		{
			name:          "randomString",
			configValue:   "randomString",
			resolverCount: 3,
			intervalValue: 100,
		},
		{
			name:          "auto1",
			configValue:   "auto",
			resolverCount: 1,
			intervalValue: 300,
		},
		{
			name:          "auto2",
			configValue:   "auto",
			resolverCount: 2,
			intervalValue: 150,
		},
		{
			name:          "auto3",
			configValue:   "auto",
			resolverCount: 3,
			intervalValue: 100,
		},
		{
			name:          "auto4",
			configValue:   "auto",
			resolverCount: 4,
			intervalValue: 75,
		},
		{
			name:          "auto5",
			configValue:   "auto",
			resolverCount: 5,
			intervalValue: 60,
		},
		{
			name:          "auto6",
			configValue:   "auto",
			resolverCount: 6,
			intervalValue: 50,
		},
		{
			name:          "auto7",
			configValue:   "auto",
			resolverCount: 7,
			intervalValue: 43,
		},
		{
			name:          "auto8",
			configValue:   "auto",
			resolverCount: 8,
			intervalValue: 38,
		},
		{
			name:          "auto9",
			configValue:   "auto",
			resolverCount: 9,
			intervalValue: 33,
		},
		{
			name:          "auto10",
			configValue:   "auto",
			resolverCount: 10,
			intervalValue: 30,
		},
		{
			name:          "auto",
			configValue:   "auto11",
			resolverCount: 11,
			intervalValue: 30,
		},
		{
			name:          "manual-10-30",
			configValue:   "30",
			resolverCount: 10,
			intervalValue: 30,
		},
		{
			name:          "manual-10-40",
			configValue:   "40",
			resolverCount: 10,
			intervalValue: 40,
		},
		{
			name:          "manual-10-20",
			configValue:   "20",
			resolverCount: 10,
			intervalValue: 30,
		},
		{
			name:          "zeroResolvers",
			configValue:   "",
			resolverCount: 0,
			intervalValue: 300,
		},
	}
)

func TestWorker(t *testing.T) {

	logger.InitLogger(&initloglevel)
	configFile := "../config.yaml"
	config.Load(&configFile)

	t.Run("getInterval", func(t *testing.T) {

		worker := Worker{}

		for _, tc := range initIntervalTestCases {
			t.Run(tc.name, func(t *testing.T) {
				viper.Set("worker.checkInterval", tc.configValue)
				worker.Resolvers = make([]resolver.Resolver, tc.resolverCount)
				worker.initInterval()
				assert.Equal(t, tc.intervalValue, worker.Interval)
			})
		}

	})
}
