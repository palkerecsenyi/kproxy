FROM golang:1.16beta1-buster
WORKDIR /go/src/app
COPY . .
RUN go get -d .
EXPOSE 80
CMD ["go", "run", "main.go"]