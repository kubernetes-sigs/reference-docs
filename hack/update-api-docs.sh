#!/usr/bin/env bash
set -euo pipefail

if [ ! -d "$K8S_WEBROOT" ]; then
  echo "Error: missing $K8S_WEBROOT directory. Exiting..."
  exit 1
fi

if [ ! -d "gen-apidocs" ]; then
  echo "Error: gen-apidocs directory not found. Exiting..."
  exit 1
else
  cd "gen-apidocs"
fi

CONFIG_DIR="config"
GITHUB_RAW_URL="https://raw.githubusercontent.com/kubernetes/kubernetes"

if ! mapfile -t latest_versions < <(curl -s "https://api.github.com/repos/kubernetes/kubernetes/tags?per_page=100" |
  jq -r '.[].name' | grep -v -e alpha -e beta -e rc |
  sort -V |
  awk -F. '!seen[$1"."$2]++ {print $0}'); then
  echo "Error: Failed to fetch or process tags. Exiting..."
  exit 1
fi

for version in "${latest_versions[@]}"; do
  directory_version=${version}
  directory_version=${directory_version%.*}
  directory_version=${directory_version//./_}
  if [ -d "$CONFIG_DIR/$directory_version" ]; then
    echo "Processing tag $version in directory $directory_version"
    if ! curl -s --output "$CONFIG_DIR/$directory_version/swagger.json" "$GITHUB_RAW_URL/$version/api/openapi-spec/swagger.json"; then
      echo "Error: Failed to download swagger.json for $version"
      continue
    fi
    if ! go run main.go --kubernetes-release="$(echo $version | sed 's/^v//' | cut -d'.' -f1-2)" --work-dir=.; then
      echo "Error: Failed to run main.go for $version"
      continue
    fi
    cp ./build/index.html "$K8S_WEBROOT/static/docs/reference/generated/kubernetes-api/${version%.*}"
  else
    echo "$directory_version NOT FOUND"
  fi
done
