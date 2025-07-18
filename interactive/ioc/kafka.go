package ioc

import (
	"github.com/IBM/sarama"
	"github.com/jym0818/webook/pkg/saramax"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Cfg struct {
		Addrs []string `yaml:"addrs"`
	}
	var cfg Cfg
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addrs, config)
	if err != nil {
		panic(err)
	}
	return client
}

func NewConsumers(c1 saramax.Consumer) []saramax.Consumer {
	return []saramax.Consumer{c1}
}
