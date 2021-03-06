FROM golang:bullseye

LABEL maintainer="eric@tedor.org"

# install streamlink
RUN echo "deb http://deb.debian.org/debian buster-backports main" | tee "/etc/apt/sources.list.d/streamlink.list"
RUN apt-get update
RUN apt-get install -y --no-install-recommends apt-utils
RUN apt-get -t buster-backports install -y streamlink

# install quadvision
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/quadvision/ ./cmd/quadvision
COPY internal/ ./internal

RUN go build -o ./quadvision ./cmd/quadvision

CMD [ "./quadvision" ]
