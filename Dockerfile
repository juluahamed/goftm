FROM golang:1.16-alpine
ARG TOKEN
ENV token=$TOKEN


RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app/goftm
# We want to populate the module cache based on the go.{mod,sum} files.

COPY . .

RUN go get

RUN go mod download


# Build the Go app
RUN go build

RUN ls


# Run the binary program produced by `go install`
# CMD ["./goftm", "-t", "echo $TOKEN"]
# CMD ./goftm -t $TOKEN
CMD ["sh", "-c", "./goftm -t $token"]