FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

WORKDIR /

COPY prestop.sh .
RUN chmod -R 777 prestop.sh

ENTRYPOINT ["sleep 10000"]
