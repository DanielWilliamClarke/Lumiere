#!/bin/bash

go get github.com/jstemmer/go-junit-report
go get github.com/axw/gocov/gocov
go get github.com/AlekSi/gocov-xml
go get -u github.com/matm/gocov-html

OUTPUT_DIR=./../test_results

# Run tests and generate coverage
go test -v -coverpkg=$(go list ./... | grep -v "test" | grep -v "mocks" | awk '{print $1}' | paste -s -d ,) -coverprofile=$OUTPUT_DIR/coverage.txt -covermode=atomic ./test/... > $OUTPUT_DIR/results.txt

cat $OUTPUT_DIR/results.txt

echo "Generate Junit Report"
# Generate Junit
go-junit-report < $OUTPUT_DIR/results.txt > $OUTPUT_DIR/TEST_ancho_report.xml

echo "Generate Coverage Report"
# Generate Cobertura
gocov convert $OUTPUT_DIR/coverage.txt > $OUTPUT_DIR/coverage.json
gocov-xml < $OUTPUT_DIR/coverage.json > $OUTPUT_DIR/TEST_ancho_coverage.xml

echo "Generate Coverage HTML Report"
# Generate HTML
gocov-html < $OUTPUT_DIR/coverage.json > $OUTPUT_DIR/index.html