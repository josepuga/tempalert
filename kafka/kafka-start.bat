@echo off
rem By José Puga 2025. GPL3 License

rem The "virtual network" to connect Zookeeper and Kafka
set NETWORK_NAME=kafka-net
set PORT_ZOOKEEPER=2181
set PORT_KAFKA=9092
rem Can be more secure with "SSH://...".
set LISTENER=PLAINTEXT://localhost:%PORT_KAFKA%
rem set TOPIC_NAME=sensor_alerts

rem Removes any running container
docker rm -f zookeeper >nul 2>&1
docker rm -f kafka >nul 2>&1

rem Create network if not exists
docker network ls --format "{{.Name}}" | findstr /C:"%NETWORK_NAME%" >nul
IF %ERRORLEVEL% NEQ 0 (
    echo Creating network %NETWORK_NAME%...
    docker network create %NETWORK_NAME%
)

echo Starting Zookeeper..
docker run --name zookeeper -d --rm ^
    --network "%NETWORK_NAME%" ^
    -p "%PORT_ZOOKEEPER%:%PORT_ZOOKEEPER%" ^
    -e ZOOKEEPER_CLIENT_PORT="%PORT_ZOOKEEPER%" ^
    docker.io/confluentinc/cp-zookeeper:latest

rem Wait few seconds to start..
timeout /t 2 /nobreak >nul
echo Started.

rem Autocreate topics enable is an easy way, if you want to test.
rem For production better use -e KAFKA_CREATE_TOPICS="$TOPIC_NAME:$PARTITIONS:$REPLICAS"
echo Starting Kafka in PLAINTEXT...
docker run --name kafka -d --rm ^
    --network "%NETWORK_NAME%" ^
    -p "%PORT_KAFKA%:%PORT_KAFKA%" ^
    -e KAFKA_BROKER_ID=1 ^
    -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 ^
    -e KAFKA_ZOOKEEPER_CONNECT="zookeeper:%PORT_ZOOKEEPER%" ^
    -e KAFKA_ADVERTISED_LISTENERS="%LISTENER%" ^
    -e KAFKA_AUTO_CREATE_TOPICS_ENABLE=true ^
    docker.io/confluentinc/cp-kafka:latest

echo Started.
