#first stage
FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o demoApp .

#second stage
FROM nginx:latest
WORKDIR /app/
COPY --from=builder /app .
EXPOSE 8080
CMD ["/app/demoApp"]