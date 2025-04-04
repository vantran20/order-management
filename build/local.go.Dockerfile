FROM golang:1.23

RUN apt-get update

RUN GO111MODULE=on go install golang.org/x/tools/cmd/goimports@v0.25.0

RUN GO111MODULE=on go install github.com/volatiletech/sqlboiler/v4@v4.16.2 && \
    GO111MODULE=on go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@v4.16.2 \