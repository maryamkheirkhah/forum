FROM golang:latest
LABEL project="Forum"
LABEL authors="Member"
LABEL version="1.0"
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod  download
COPY . ./
RUN go build -o server .
RUN rm Dockerfile
RUN rm go.mod
RUN rm main.go
RUN rm run-docker.sh
RUN ls -la
#Run command "go build -o main" to make app and name it "main"
EXPOSE 8080
CMD ["/app/server"]