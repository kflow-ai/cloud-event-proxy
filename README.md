# cloud-event-proxy

A simple web app to proxy CloudEvents to GCP PubSub.

In the future, we will also support AWS Kinesis Data Firehose

## Setup

Set up the Application Default Credentials on your local CLI to the target
Google Cloud project:

```bash
gcloud auth application-default login
```

Ensure your PubSub topic exists and set the following env vars:
```bash
export CEP_PUBSUB_PROJECT_ID=<your-project-id>
export CEP_PUBSUB_TOPIC_ID=<your-topic-id>
```

Then either run in Docker:

```bash
docker compose up
```

Or run using your local Go development environment:

```bash
go run cmd/cloudeventproxy/main.go
```

## Generate License

Install `addlicense`:

```
go install github.com/google/addlicense@v1.0.0
```

Make sure all files contain a license:

```
addlicense -c "Cake AI Technologies, Inc." -y $(date +"%Y") -l apache -s=only ./**/*.go
```
