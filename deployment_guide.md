# Deployment Guide

This guide provides instructions for setting up the project for local development and deploying it to Google Cloud.

## Prerequisites

Ensure you have the following tools installed on your system:
- [Go](https://golang.org/doc/install) (check version compatibility)
- [Google Cloud SDK](https://cloud.google.com/sdk/install) (`gcloud` CLI)
- Docker (for local container testing)

## Local Development

1.  **Run the setup script** to initialize your environment. This script will install dependencies and set up a local PostgreSQL database.
    ```bash
    ./setup.sh
    ```
2. Run `go run main.go` to test locally.

## Google Cloud Deployment

Follow these steps to deploy the application to Google Cloud Run.

### 1. (Optional) Store Database Credentials in Secret Manager

For better security, store your database password in Google Secret Manager instead of passing it as a plain text environment variable.

1.  Create a new secret:
    `gcloud secrets create [SECRET_NAME] --replication-policy="automatic"`

2.  Add the password as a secret version:
    `printf "[DB_PASSWORD]" | gcloud secrets versions add [SECRET_NAME] --data-file=-`

3.  Grant your Cloud Run service account access to the secret.

### 1. Build the Container Image

Use Google Cloud Build to create a container image from the source code.

`gcloud builds submit --tag gcr.io/[PROJECT-ID]/[APP-NAME]`

### 2. Deploy to Cloud Run

Deploy the container image to Cloud Run and connect it to a Cloud SQL database.

`gcloud run deploy [SERVICE-NAME] \
    --image gcr.io/[PROJECT-ID]/[APP-NAME] \
    --platform managed \
    --region [REGION] \
    --allow-unauthenticated \
    --add-cloudsql-instances="[INSTANCE_CONNECTION_NAME]" \
    --set-env-vars="DB_USER=[DB_USER],DB_NAME=[DB_NAME],INSTANCE_CONNECTION_NAME=[INSTANCE_CONNECTION_NAME]" \
    --set-secrets="DB_PASSWORD=[SECRET_NAME]:latest"`

**Note on `--allow-unauthenticated`**: This flag makes your service publicly accessible. For production environments, it is highly recommended to use an authentication method to secure your service.

## Firebase Deployment (using Cloud Run)

You can deploy your containerized application as a service on Firebase, which uses Cloud Run behind the scenes. This is ideal if you are using Firebase Hosting for a frontend.

### Prerequisites

- Install Node.js and npm.
- Install the Firebase CLI: `npm install -g firebase-tools`
- Log in to the Firebase CLI: `firebase login`

### 1. Initialize Firebase

Run `firebase init` in your project root.

- Select **Hosting: Configure files for Firebase Hosting...**.
- Choose a Firebase project.
- When asked for your public directory, you can use `frontend`.
- When asked to **Configure as a single-page app**, say **No**.
- When asked to **Set up automatic builds and deploys with GitHub?**, say **No** for now.

This will create `firebase.json` and `.firebaserc` files.

### 2. Build and Deploy

The Firebase CLI can build your container and deploy it to Cloud Run in a single step.

Run the following command:

`firebase deploy --only hosting,run`

The CLI will prompt you for the service name, region, and other configuration details based on your `gcloud` settings. It will use your `Dockerfile` to build the image, push it to Artifact Registry, and deploy it to Cloud Run, automatically connecting it to your Firebase Hosting URL.
