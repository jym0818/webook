package ioc

import (
	"github.com/IBM/sarama"
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

func InitKafkaProducer(client sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return producer
}
