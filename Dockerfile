
FROM golang:1.17-buster AS builder

ADD . /src

RUN cd /src \
  && go get -u -v golang.org/x/lint/golint \
  && go mod tidy \
  && go get -u -v \
  && go mod download \
  && golint . \
  && export CI=1 \
  && go test -short -covermode=count -coverprofile=coverage.out \
  && cat coverage.out | grep -v "main.go" | grep -v "clickhouse.go" > coverage.txt \
  && TOTAL_COVERAGE_FOR_CI_F=$(go tool cover -func coverage.txt | grep total | grep -Eo '[0-9]+.[0-9]+') \
  && echo "TOTAL_COVERAGE_FOR_CI_F: $TOTAL_COVERAGE_FOR_CI_F" \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o analytics-collector \
  && ls -lah /src/analytics-collector

FROM scratch

COPY --from=builder /src/analytics-collector /usr/bin/analytics-collector
ADD tmp/GeoLite2-ASN.mmdb /GeoLite2-ASN.mmdb
ADD tmp/GeoLite2-City.mmdb /GeoLite2-City.mmdb

ENTRYPOINT ["/usr/bin/analytics-collector"]
