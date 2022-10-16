go test ./go_api/... -v -covermode=atomic -coverprofile=./go_api/coverage_raw.txt -coverpkg=./go_api/...

# remove any lines with mocks in them to exclude mock files from the code coverage
cat ./go_api/coverage_raw.txt | grep -v "mocks" > ./go_api/coverage.txt