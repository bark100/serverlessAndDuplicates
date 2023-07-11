#!/bin/bash
set -e

if [ $# -ne 2 ]; then
    cat <<EOF
Usage: $0 <lambda_src> <event_json>
Example: $0 default calcOccurrences test.json
EOF
    exit 1
fi

rm -f ../"$2"/.build/main
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -C ../lambdas/"$1" -o main main.go

sam local invoke \
  --hook-name terraform \
  --beta-features "module.$1[0].aws_lambda_function.this[0]" \
  --event ../lambdas/"$1"/events/"$2"
