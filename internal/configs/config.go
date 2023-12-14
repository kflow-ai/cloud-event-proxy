// Copyright 2023 Cake AI Technologies, Inc.
// SPDX-License-Identifier: Apache-2.0

package configs

import (
	"fmt"
	"os"
	"strconv"
)

const (
	// DefaultListenPort is the default port to listen on for incoming HTTP requests.
	DefaultListenPort = 8080

	// DefaultListenAddress is the default address to listen on for incoming HTTP requests.
	DefaultListenAddress = "0.0.0.0"

	// DefaultLogLevel is the default log level to use.
	DefaultLogLevel = "info"

	// ConfigEnvPrefix is the prefix for environment variables.
	ConfigEnvPrefix = "CEP_"

	// ConfigEnvPrefixPubSub is the prefix for Google PubSub environment variables.
	ConfigEnvPrefixPubSub = "CEP_PUBSUB_"

	// ConfigEnvPrefixFirehose is the prefix for Kinesis Data Firehose environment variables.
	ConfigEnvPrefixFirehose = "CEP_FIREHOSE_"
)

// DestinationAdapter is an enum of supported destination adapters.
type DestinationAdapter string

const (
	// DestinationAdapterPubSub is the PubSub destination adapter.
	DestinationAdapterPubSub DestinationAdapter = "pubsub"

	// DestinationAdapterFirehose is the Kinesis Data Firehose destination adapter.
	DestinationAdapterFirehose DestinationAdapter = "firehose"
)

// Config is the configuration for the CloudEventProxy.
type Config struct {
	// ListenPort is the port to listen on for incoming HTTP requests.
	ListenPort int

	// LogLevel is the log level to use.
	LogLevel string

	// DestinationAdapter is the destination adapter configuration.
	DestinationAdapter DestinationAdapter

	// PubSub is the PubSub configuration.
	PubSub *PubSubConfig

	// Firehose is the Kinesis Data Firehose configuration.
	Firehose *FirehoseConfig
}

// PubSubConfig is the configuration for the PubSub client.
type PubSubConfig struct {
	// ProjectID is the GCP project ID.
	ProjectID string

	// TopicID is the GCP PubSub topic ID.
	TopicID string
}

type FirehoseConfig struct {
	// DeliveryStreamARN is the ARN of the Kinesis Data Firehose delivery stream.
	DeliveryStreamARN string
}

func Load() (*Config, error) {
	listenPortStr := os.Getenv(ConfigEnvPrefix + "LISTEN_PORT")
	var listenPort int
	if listenPortStr == "" {
		listenPort = DefaultListenPort
	} else {
		var err error
		listenPort, err = strconv.Atoi(listenPortStr)
		if err != nil {
			return nil, fmt.Errorf("invalid listen port: %w", err)
		}
	}

	logLevel := os.Getenv(ConfigEnvPrefix + "LOG_LEVEL")
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}

	adapter := os.Getenv(ConfigEnvPrefix + "DESTINATION_ADAPTER")
	if adapter == "" {
		return nil, fmt.Errorf("missing %s", ConfigEnvPrefix+"DESTINATION_ADAPTER")
	}

	var config *Config

	switch adapter {
	case string(DestinationAdapterPubSub):
		pubSub, err := LoadPubSub()
		if err != nil {
			return nil, fmt.Errorf("failed to load PubSub configuration: %w", err)
		}
		config = &Config{
			ListenPort:         listenPort,
			LogLevel:           logLevel,
			DestinationAdapter: DestinationAdapterPubSub,
			PubSub:             pubSub,
		}
	case string(DestinationAdapterFirehose):
		firehose, err := LoadFirehose()
		if err != nil {
			return nil, fmt.Errorf("failed to load Kinesis Data Firehose configuration: %w", err)
		}
		config = &Config{
			ListenPort:         listenPort,
			LogLevel:           logLevel,
			DestinationAdapter: DestinationAdapterFirehose,
			Firehose:           firehose,
		}
	default:
		return nil, fmt.Errorf("invalid %s: %s", ConfigEnvPrefix+"DESTINATION_ADAPTER", adapter)
	}
	return config, nil
}

func LoadPubSub() (*PubSubConfig, error) {
	projectID := os.Getenv(ConfigEnvPrefixPubSub + "PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("missing %s", ConfigEnvPrefixPubSub+"PROJECT_ID")
	}

	topicID := os.Getenv(ConfigEnvPrefixPubSub + "TOPIC_ID")
	if topicID == "" {
		return nil, fmt.Errorf("missing %s", ConfigEnvPrefixPubSub+"TOPIC_ID")
	}

	return &PubSubConfig{
		ProjectID: projectID,
		TopicID:   topicID,
	}, nil
}

func LoadFirehose() (*FirehoseConfig, error) {
	deliveryStreamARN := os.Getenv(ConfigEnvPrefixFirehose + "DELIVERY_STREAM_ARN")
	if deliveryStreamARN == "" {
		return nil, fmt.Errorf("missing %s", ConfigEnvPrefixFirehose+"DELIVERY_STREAM_ARN")
	}

	return &FirehoseConfig{
		DeliveryStreamARN: deliveryStreamARN,
	}, nil
}
