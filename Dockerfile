FROM golang:buster AS builder

RUN mkdir /usr/local/kumquat
WORKDIR /usr/local/kumquat

COPY . .

RUN go get
RUN go build -o kumquat *.go

FROM debian:stable
RUN mkdir /usr/local/kumquat
WORKDIR /usr/local/kumquat
COPY --from=builder /usr/local/kumquat/kumquat .
COPY --from=builder /usr/local/kumquat/pass.sh .

RUN addgroup --gid 1000 kumquat && \
    adduser --home /usr/local/kumquat -u 1000 --gid 1000 kumquat && \
    chown -R kumquat:kumquat /usr/local/kumquat

RUN chmod a+x /usr/local/kumquat/pass.sh

RUN apt-get update -y
RUN apt-get install curl -y
RUN apt-get install apache2-utils -y
RUN apt-get install -y python3 && apt-get install -y git
RUN apt-get update && apt-get install -y procps
RUN apt-get clean -y

USER kumquat
WORKDIR /usr/local/kumquat
ENTRYPOINT [ "/usr/local/kumquat/kumquat" ]