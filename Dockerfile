FROM centos:7
WORKDIR /root
COPY main ./server
COPY configs/ ./configs/
EXPOSE 9090
RUN /root/server
