package adapters

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/kflow-ai/cloud-event-proxy/internal/configs"
)

func StartFirehose(ctx context.Context, cfg *configs.Config) {
	// Parse from ARN
	arn, err := arn.Parse(cfg.Firehose.DeliveryStreamARN)
	if err != nil {
		log.Fatalf("failed to parse Firehose ARN: %v", err)
	}
	resource := arn.Resource
	streamName := resource[15:] // strip "deliverystream/"
	region := arn.Region

	// Create Kinesis Data Firehose client
	awsCfg := aws.NewConfig().WithRegion(region)
	firehoseClient := firehose.New(session.Must(session.NewSession()), awsCfg)

	log.Printf("using Firehose ARN %s", cfg.Firehose.DeliveryStreamARN)

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
		log.Printf("data %s", data)

		// TODO: batch? retry?
		resp, err := firehoseClient.PutRecord(&firehose.PutRecordInput{
			DeliveryStreamName: &streamName,
			Record: &firehose.Record{
				Data: data,
			},
		})
		log.Printf("PUT %v", resp)
		log.Printf("err %v", err)
	}))
}
