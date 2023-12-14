// Copyright 2023 Cake AI Technologies, Inc.
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"context"
	"log"

	"github.com/kflow-ai/cloud-event-proxy/internal/adapters"
	"github.com/kflow-ai/cloud-event-proxy/internal/configs"
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

	switch cfg.DestinationAdapter {
	case configs.DestinationAdapterPubSub:
		adapters.StartPubSub(ctx, cfg)
	case configs.DestinationAdapterFirehose:
		adapters.StartFirehose(ctx, cfg)
	default:
		log.Fatalf("unknown destination adapter: %s", cfg.DestinationAdapter)
	}
}
