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

	// ConfigEnvPrefixPubSub is the prefix for PubSub environment variables.
	ConfigEnvPrefixPubSub = "CEP_PUBSUB_"
)

// Config is the configuration for the CloudEventProxy.
type Config struct {
	// ListenPort is the port to listen on for incoming HTTP requests.
	ListenPort int

	// LogLevel is the log level to use.
	LogLevel string

	// PubSub is the PubSub configuration.
	PubSub *PubSubConfig
}

// PubSubConfig is the configuration for the PubSub client.
type PubSubConfig struct {
	// ProjectID is the GCP project ID.
	ProjectID string

	// TopicID is the GCP PubSub topic ID.
	TopicID string
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

	pubSub, err := LoadPubSub()
	if err != nil {
		return nil, fmt.Errorf("failed to load PubSub configuration: %w", err)
	}

	return &Config{
		ListenPort: listenPort,
		LogLevel:   logLevel,
		PubSub:     pubSub,
	}, nil
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
