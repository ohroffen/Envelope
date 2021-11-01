FROM centos:7
COPY main /root/server
EXPOSE 9090
CMD /root/server
