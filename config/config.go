package config

import (
	"flag"

	"github.com/fsnotify/fsnotify"
	"github.com/kerti/cloudflare-ddns/logger"
	"github.com/spf13/viper"
)

// Load pre-sets and loads the config
func Load(optConfig *string) error {
	/* parse flags first */
	flag.Parse()

	/* set config locations */
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigFile(*optConfig)

	/* set default config values */
	viper.SetDefault("loglevel", 3)
	viper.SetDefault("resolver.noVerify", true)
	viper.SetDefault("worker.checkInterval", "auto")

	viper.SetDefault("notifier.ifttt.webhook.active", false)
	viper.SetDefault("notifier.ifttt.webhook.eventName", "cf_ddns_update")

	/* read the config file */
	err := viper.ReadInConfig()
	if err != nil {
		logger.Warn("[cf-ddns] failed to read config file(%s), running on defaults", *optConfig)
	} else {
		logger.Info("[cf-ddns] config file (%s) successfully read", *optConfig)
	}

	logger.ResetLogLevel()

	/* and watch it */
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Print("[cf-ddns] configuration reloaded...")
		logger.ResetLogLevel()
	})

	return nil
}
