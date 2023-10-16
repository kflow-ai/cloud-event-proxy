# cloud-event-proxy

A simple web app to proxy CloudEvents to GCP PubSub or AWS Kinesis Data Firehose

## Generate License

Install `addlicense`:

```
go install github.com/google/addlicense@v1.0.0
```

Make sure all files contain a license:

```
addlicense -c "Cake AI Technologies, Inc." -y $(date +"%Y") -l apache -s=only ./**/*.go
```
