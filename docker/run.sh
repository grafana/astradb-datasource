# This image runs stargate and cassandra combined

# https://hub.docker.com/r/stargateio/stargate-4_0
# docker pull stargateio/stargate-4_0:v1.0.64

docker run --name stargate \
  -p 8081:8081 \
  -p 8082:8082 \
  -p 8090:8090 \
  -p 127.0.0.1:9042:9042 \
  -d \
  -e CLUSTER_NAME=stargate \
  -e CLUSTER_VERSION=3.11 \
  -e DEVELOPER_MODE=true \
  stargateio/stargate-3_11:v1.0.40
