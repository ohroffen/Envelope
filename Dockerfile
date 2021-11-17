FROM centos:7
ENV KAFKA_HOST kafka-ymqfnev51bodns.cn-beijing.kafka-internal.ivolces.com:9092
ENV KAFKA_TOPIC wri
ENV REDIS_HOST redis-cn02nwh9l8hdm1gmx.redis.ivolces.com:6379
ENV REDIS_PASSWORD Group12345678
ENV GIN_MODE release
WORKDIR /root
COPY main ./server
EXPOSE 9090
CMD /root/server