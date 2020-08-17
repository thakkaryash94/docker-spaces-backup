# Docker Spaces Backup ![Build and Push Docker](https://github.com/thakkaryash94/docker-spaces-backup/workflows/Build%20and%20Push%20Docker/badge.svg)

Docker image for Digital Ocean spaces backup

### Environment Variables

- ACCESS_KEY_ID: Spaces access key id
- BUCKET_NAME: Spaces bucket name
- CRON_SCHEDULE: Cron value in double quotes. https://godoc.org/github.com/robfig/cron
- S3_URL: Spaces url(nyc3.digitaloceanspaces.com)
- SECRET_ACCESS_KEY: Spaces secret access key

### Volumes

- mount backup folder with `/data` path

#### Example

```sh
docker run -v $(pwd)/data:/data:ro \
     -e BUCKET_NAME=YOUR_BUCKET_NAME \
     -e S3_URL=YOUR_S3_URL \
     -e ACCESS_KEY_ID=YOUR_ACCESS_KEY_ID \
     -e SECRET_ACCESS_KEY=YOUR_SECRET_ACCESS_KEY \
     -e CRON_SCHEDULE=YOUR_CRON_SCHEDULE \
     -d docker.pkg.github.com/thakkaryash94/docker-spaces-backup:latest
```
