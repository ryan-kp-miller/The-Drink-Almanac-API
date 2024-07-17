go test ./go_api/... -v -covermode=atomic -coverprofile=./go_api/coverage_raw.txt -coverpkg=./go_api/...
RETVAL=$?

# if any of the tests fail, exit script with non-zero code
# to force the github actions step to fail
if [[ ${RETVAL} != 0 ]]; then
    echo "At least one test failed, so exiting the script"
    exit -1
fi

# join the expressions together and then remove the leading delimiter
file_path_regexp_to_ignore=("mocks" "Stub")
combined_regexp=$(printf "\|%s" "${file_path_regexp_to_ignore[@]}")
combined_regexp=${combined_regexp:2}
echo "combined expression ${combined_regexp}"

# remove any lines that match one of the expressions in them from the code coverage
cat ./go_api/coverage_raw.txt | grep -v $combined_regexp > ./go_api/coverage.txt