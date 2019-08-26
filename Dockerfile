FROM alpine:3.10
# Install packages
RUN apk add --no-cache ca-certificates docker go git libc-dev
# Install go packages
RUN go get -v -u github.com/nlopes/slack
RUN go get -v -u gopkg.in/yaml.v2
RUN go get -v -u github.com/docker/docker/api
RUN go get -v -u github.com/docker/docker/client
# Copy app source code
ADD . /app/
WORKDIR /app
# Build binary
RUN go build -o main .
# Entrypoint
CMD ["./main"]