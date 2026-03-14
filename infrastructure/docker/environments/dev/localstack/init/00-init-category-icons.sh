#!/usr/bin/env bash
set -euo pipefail

bucket="aion-assets"
assets_dir="/opt/localstack/assets"

awslocal s3api create-bucket --bucket "$bucket" >/dev/null 2>&1 || true

awslocal s3api put-bucket-cors --bucket "$bucket" --cors-configuration '{
  "CORSRules": [
    {
      "AllowedOrigins": ["*"],
      "AllowedMethods": ["GET", "HEAD"],
      "AllowedHeaders": ["*"]
    }
  ]
}' >/dev/null 2>&1 || true

awslocal s3api put-bucket-policy --bucket "$bucket" --policy '{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicRead",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::'"$bucket"'/*"
    }
  ]
}' >/dev/null 2>&1 || true

if [ -d "$assets_dir" ]; then
  awslocal s3 sync "$assets_dir" "s3://$bucket" --acl public-read >/dev/null
fi

echo "✅ Localstack assets seeded in s3://$bucket"
