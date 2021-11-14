FROM centos:7
ENV GO_ENV product
WORKDIR /root
COPY main ./server
COPY configs/ ./configs/
EXPOSE 9090
CMD /root/server