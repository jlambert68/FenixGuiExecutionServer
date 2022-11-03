# Compile stage
FROM golang:1.18 AS build-env

ADD . /dockerdev
WORKDIR /dockerdev

COPY go.* ./
RUN go mod tidy

RUN go build -o /fenixGuiExecutionServer .


# Final stage
FROM debian:buster
#FROM golang:1.13.8

EXPOSE 6670
#FROM golang:1.13.8
WORKDIR /
COPY --from=build-env /fenixGuiExecutionServer /
#Add data/ data/

#CMD ["/fenixClientServer"]
ENTRYPOINT ["/fenixGuiExecutionServer"]



#// docker build -t  fenix-client-server .
#// docker run -p 5998:5998 -it  fenix-client-server
#// docker run -p 5998:5998 -it --env StartupType=LOCALHOST_DOCKER fenix-client-server

#//docker run --name fenix-client-server --rm -i -t fenix-client-server  bash