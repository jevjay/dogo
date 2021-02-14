FROM alpine:3.10
# Install packages
RUN apk add --no-cache ca-certificates docker go git libc-dev
# Copy app source code
ADD . /app/
WORKDIR /app
# Build binary
RUN go build -o main .
# Entrypoint
CMD ["./main"]
