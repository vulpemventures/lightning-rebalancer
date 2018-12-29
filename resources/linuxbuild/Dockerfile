FROM golang:latest
WORKDIR /
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN mkdir -p /go/src/github.com/vulpemventures/lightning-rebalancer
WORKDIR /go/src/github.com/vulpemventures/lightning-rebalancer
COPY . /go/src/github.com/vulpemventures/lightning-rebalancer
RUN chmod +x scripts/install
RUN ./scripts/install
RUN chmod +x ./scripts/build
RUN scripts/build linux arm 
