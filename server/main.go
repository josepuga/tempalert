package main

import (
	"fmt"
	"strconv"
	"time"

	pb "tempserver/proto/temperatures"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/josepuga/goini"
	"google.golang.org/protobuf/proto"
)

const configFile = "./config.ini"
const readDelay = 3
const topic = "sensor_alert"

var si SensorInfo
var producer *kafka.Producer

func main() {
	// Init Kafka Producer
	//
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		fmt.Printf("Error initializing Kafka Producer: %v\n", err)
		return
	}
	fmt.Println("Kafka Producer initialized sucessfully.")
	defer producer.Close()

	// Load SensorInfo data from ini file
	//
	if err := loadSensorInfo(); err != nil {
		fmt.Printf("%v \n", err)
		return
	}

	// Infinite loop to simulate the sensors monitoriz.
	for {
		time.Sleep(readDelay * time.Second)

		// Read the temps in the mock server
		if err := si.ReadTemps(); err != nil {
			fmt.Printf("Error reading temps: %v\n", err)
			continue
		}

		// Print temps to stdout
		fmt.Println(si)

		// Check for temps out of safe range and send notification to Kafka
		for i, temp := range si.Temps {
			if !si.SensorTempIsSafe(i) {
				if err := sendAlert(i, temp); err != nil {
					fmt.Printf("Error sending alert to Kafka: %v", err)
				}
			}
		}
	}
}

// sendAlert to Kafka
func sendAlert(sensor, temp int) error {

	// Set protobuf data
	alert := &pb.SensorAlert{
		SensorId:    int32(sensor),
		Temperature: int32(temp),
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	// Serialize protobuf to binary
	alertBytes, err := proto.Marshal(alert)
	if err != nil {
		return err
	}

	// Send message to Kafka
	alertTopic := topic
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &alertTopic,
			Partition: kafka.PartitionAny},
		Value: alertBytes,
	}

	if err = producer.Produce(message, nil); err != nil {
		return err
	}

	return nil
}

// loadSensorInfo load configuration about the fake sensors
func loadSensorInfo() error {
	ini := goini.NewIni()
	if err := ini.LoadFromFile(configFile); err != nil {
		return err
	}
	// Fake sensors to use
	si.SensorsCount = ini.GetInt("", "sensors", 1)

	// Set the temps slice size
	si.Temps = make([]int, si.SensorsCount)

	// Valid temps range
	temps := ini.GetStringSlice("", "temp range", "", ",")
	//Warning: No error check.
	si.MinReadableTemp, _ = strconv.Atoi(temps[0])
	si.MaxReadableTemp, _ = strconv.Atoi(temps[1])

	// Safety temps range
	temps = ini.GetStringSlice("", "safe range", "", ",")
	//Warning: No error check.
	si.MinSafeTemp, _ = strconv.Atoi(temps[0])
	si.MaxSafeTemp, _ = strconv.Atoi(temps[1])

	return nil
}
