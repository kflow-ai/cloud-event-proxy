// Copyright 2023 Cake AI Technologies, Inc.
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"context"
	"log"

	"github.com/kflow-ai/cloud-event-proxy/internal/configs"

	"cloud.google.com/go/pubsub"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func Start() {
	ctx := context.Background()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Printf("starting CloudEventProxy")

	// Load the configuration
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Create the PubSub client
	pubsubClient, err := pubsub.NewClient(ctx, cfg.PubSub.ProjectID)
	if err != nil {
		log.Fatalf("failed to create PubSub client: %v", err)
	}

	topic := pubsubClient.Topic(cfg.PubSub.TopicID)

	log.Printf("using Pub/Sub project %s", cfg.PubSub.ProjectID)
	log.Printf("using Pub/Sub topic %s", topic.String())

	// Set up the CloudEvents HTTP receiver
	p, err := cloudevents.NewClientHTTP(
		cloudevents.WithPort(cfg.ListenPort),
	)
	if err != nil {
		log.Fatalf("failed to create protocol: %v", err)
	} else {
		log.Printf("listening on port %d", cfg.ListenPort)
	}

	log.Fatal(p.StartReceiver(ctx, func(event cloudevents.Event) {
		data, err := event.MarshalJSON()
		if err != nil {
			log.Printf("failed to marshal event: %v", err)
			return
		}

		msg := &pubsub.Message{
			Data: data,
		}

		if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
			log.Printf("failed to publish to Pub/Sub: %v", err)
			return
		}
	}))
}
