FROM golang:latest

WORKDIR /app
COPY ./ ./
LABEL maintainer = "dkalykov"
LABEL description="This Dockerfile is compiled by Forum"

RUN go build -o main .
EXPOSE 8080


CMD [ "./main" ]