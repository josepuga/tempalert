use std::process;

// By José Puga 2025. GPL3
// Based on the tutorial from https://www.influxdata.com/blog/building-simple-pure-rust-async-apache-kafka-client/
use rdkafka::consumer::{Consumer, StreamConsumer};
use rdkafka::message::Message;
use rdkafka::ClientConfig;
use futures::StreamExt;
use tokio;



// Trait Message to decode(). Alias because conflict with rdkafka message
use prost::Message as ProstMessage;

// Proto generated module
mod temperatures;
use temperatures::SensorAlert;

const TOPIC : &str = "sensor_alert";

#[tokio::main]
async fn main() {
    // Setup Kafka consumer
    let consumer: StreamConsumer = match ClientConfig::new()
        .set("bootstrap.servers", "localhost:9092")
        .set("group.id", "consumer_group")
        .set("auto.offset.reset", "earliest")
        .create()
    {
        Ok(c) => c,
        Err(e) => {
            eprintln!("Unable to create Kafka consumer: {}",e);
            process::exit(1)
        }
    };

    // Suscribe to the alert topic
    consumer
        .subscribe(&[TOPIC])
        .expect("Unable to suscribe to the topic");

    println!("Waiting for alerts...");

    // Consume kafka messages
    let mut stream = consumer.stream();
    while let Some(result) = stream.next().await {
        match result {
            Ok(message) => {
                if let Some(payload) = message.payload() {
                    // Deserializa the message
                    match SensorAlert::decode(payload) {

                        Ok(alert) => {
                            // More "polite" time format.
                            println!("Alert: Sensor {:>2}. Temp {:>4} @ {}",
                            alert.sensor_id, alert.temperature, alert.timestamp);}
                        Err(e) => {println!("Error unserializing alert: {}", e);}
                    }
                }
            }
            Err(e) => eprintln!("Error receiving message: {}", e),
        }
    }
}


// Unit test for serialize/deseralize message
#[cfg(test)]
mod tests {
    use super::*;
    use prost::Message;

    #[test]
    fn test_sensor_alert_serialization() {
        let alert = SensorAlert {
            sensor_id: 1,
            temperature: 200,
            timestamp: "2025-03-09T19:48:20+01:00".to_string(),
        };

        //alert.enconde_to_vec() is the same as proto.Marshal(alert) in the Go Server
        let encoded = alert.encode_to_vec();
        let decoded = SensorAlert::decode(&*encoded).expect("Failed to decode alert");

        assert_eq!(alert.sensor_id, decoded.sensor_id);
        assert_eq!(alert.temperature, decoded.temperature);
        assert_eq!(alert.timestamp, decoded.timestamp);
    }
}
