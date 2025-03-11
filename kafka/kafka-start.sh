#!/bin/bash
# By José Puga 2025. GPL3 License


# The "virtual network" to connects zookeeper and kafka
NETWORK_NAME="kafka-net"
PORT_ZOOKEEPER=2181
PORT_KAFKA=9092
# Can be more secure with "SSH://...".
LISTENER="PLAINTEXT://localhost:$PORT_KAFKA"
# TOPIC_NAME="sensor_alerts"

# Removes any running container
docker rm -f zookeeper &>/dev/null
docker rm -f kafka &>/dev/null

set -e

# Create net if not exits
if ! docker network ls --format "{{.Name}}" | grep -q -w "$NETWORK_NAME"; then
    echo "Creating network $NETWORK_NAME..."   
    docker network create "$NETWORK_NAME"
fi

echo "Starting Zookeeper.."
docker run --name zookeeper -d --rm \
    --network "$NETWORK_NAME" \
    -p "$PORT_ZOOKEEPER:$PORT_ZOOKEEPER" \
    -e ZOOKEEPER_CLIENT_PORT="$PORT_ZOOKEEPER" \
    docker.io/confluentinc/cp-zookeeper:latest

# Wait few seconds to start..
sleep 2
echo "Started."

# Autocreate topics enable is an easy way, if you wan to test.
# For production better use -e KAFKA_CREATE_TOPICS="$TOPIC_NAME:$PARTITIONS:$REPLICAS"
echo "Starting Kafka in PLAINTEXT..."
docker run --name kafka -d --rm \
    --network "$NETWORK_NAME" \
    -p "$PORT_KAFKA:$PORT_KAFKA" \
    -e KAFKA_BROKER_ID=1 \
    -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
    -e KAFKA_ZOOKEEPER_CONNECT=zookeeper:"$PORT_ZOOKEEPER" \
    -e KAFKA_ADVERTISED_LISTENERS="$LISTENER" \
    -e KAFKA_AUTO_CREATE_TOPICS_ENABLE=true \
    docker.io/confluentinc/cp-kafka:latest

echo "Started."



