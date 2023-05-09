FROM golang:1.20

WORKDIR /chaincode
COPY . .

RUN go mod tidy
#RUN go build -o ivs_chaincode
RUN go build -v -o ivs_chaincode

CMD ["./ivs_chaincode"]

