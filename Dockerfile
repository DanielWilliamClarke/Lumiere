#DEPENDENCIES -----------------------------------
FROM golang:1.13-alpine as builder

LABEL maintainer="Daniel Clarke (clarkit@gmail.com)"

COPY /src /lumiere
WORKDIR /lumiere

RUN go mod tidy && go build

#RELEASE -----------------------------------
FROM alpine:3.11.6

RUN mkdir ./lumiere
WORKDIR /lumiere

COPY --from=builder /lumiere/lumiere lumiere

## Add the wait script to the image
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

# Create user
RUN adduser -D user && \
  chown -R user /lumiere
USER user

# Set ports
ENV PORT 5000
EXPOSE 5000

ENTRYPOINT /wait && ./lumiere