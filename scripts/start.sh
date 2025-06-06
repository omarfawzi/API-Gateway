#!/bin/sh
set -e

gomplate -f "${CONFIG_TEMPLATE_FILE:-config/config.yaml}" -o config/rendered_config.yaml

yq eval -o=json config/rendered_config.yaml > config/config.json

rm config/rendered_config.yaml

exec ./app -p "${APP_PORT:-8080}"
