# Deployment Guide

## Local Development
1. Run `./setup.sh` to initialize.
2. Run `go run main.go` to test locally.

## Google Cloud Deployment

### 1. Build the Container Image
Use Google Cloud Build to create a container image from the source code.
`gcloud builds submit --tag gcr.io/[PROJECT-ID]/[APP-NAME]`

### 2. Deploy to Cloud Run
Deploy the container image to Cloud Run and connect it to a Cloud SQL database.
`gcloud run deploy [SERVICE-NAME] --image gcr.io/[PROJECT-ID]/[APP-NAME] --platform managed --region [REGION] --allow-unauthenticated --add-cloudsql-instances="[INSTANCE_CONNECTION_NAME]" --set-env-vars="DATABASE_URL=postgres://[DB_USER]:[DB_PASSWORD]@/[DB_NAME]?host=/cloudsql/[INSTANCE_CONNECTION_NAME]"`
