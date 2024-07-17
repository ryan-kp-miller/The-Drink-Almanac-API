aws lambda invoke \
  --function-name favorites-lambda \
  --endpoint-url http://localhost:4566 \
  --payload '{"path": "test"}' \
  stdout