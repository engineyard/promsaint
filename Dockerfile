FROM golang:1.11-alpine

COPY . /go/src/github.com/engineyard/promsaint

COPY ./run.sh /run.sh

RUN go install \
  github.com/engineyard/promsaint/cmd/promsaint && \
  rm -rf /go/src

USER nobody

CMD ["/run.sh"]
