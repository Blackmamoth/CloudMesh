#!/bin/sh

set -e

for secret in /run/secrets/*; do
  if [ -f "$secret" ]; then
    secret_name=$(basename "$secret")
    secret_value=$(cat "$secret")
    export "$secret_name"="$secret_value"
  fi
done

sh ./setup_oauth_config.sh oauth.config.sample.yml oauth.config.yml

exec /bin/cloudmesh
