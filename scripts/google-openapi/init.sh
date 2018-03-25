#! /bin/bash -eu

# example.
# $ ./scripts/google-openapi/init.sh your-gcp-project-name

cat scripts/google-openapi/openapi.yaml \
  | sed "s/GCP_PROJECT_ID/$1/g" \
  > ./openapi.yaml
gcloud endpoints services deploy openapi.yaml
rm ./openapi.yaml
