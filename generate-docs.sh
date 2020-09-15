#!/bin/sh

data_source="$1"
stdlib_target=content/docs/stdlib/
stdlib_template=templates/docs/stdlib.template

if [ ! -f $data_source ]; then
  echo "Missing stdlib data source";
  exit 255;
fi

for module in $(yq read $data_source 'modules[*].name'); do
  safe_package="$(echo "${module}" | sed 's/\//-/g')"

  USING_KEY="$module" frep \
    --load "$data_source" \
    --overwrite \
    "$stdlib_template":"${stdlib_target}${safe_package}.md"
done
