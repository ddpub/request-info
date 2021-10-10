FROM golang:1.17 as build
USER root

WORKDIR /root
COPY . ./src

RUN cd src; go mod download
RUN cd ~/src; go test --race --cover -v ./...

RUN mkdir bin
RUN cd src; go build -o ../bin/request-info


FROM ubuntu:20.04 as release

ENV TZ 'Europe/Moscow'
RUN apt-get update; echo $TZ > /etc/timezone \
    && apt-get install -y tzdata \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && dpkg-reconfigure -f noninteractive tzdata \
    && apt-get install curl -y \
    && apt-get clean


COPY --from=build /root/bin/request-info /bin/request-info
RUN chmod +x /bin/request-info

USER 65534:65534

CMD ["/bin/request-info", "-srv.addr=:7777"]

