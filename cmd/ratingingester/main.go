package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/TylerAldrich814/MetaReviews/rating/pkg/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main(){
  fmt.Println("Creating a Kafka Producer")

  producer, err := kafka.NewProducer(
    &kafka.ConfigMap{
      "bootstrp.servers":"localhost",
    },
  )
  if err != nil {
    panic(err)
  }
  defer producer.Close()
  
  const fileName = "ratingsdata.json"
  fmt.Println("Reading rating events from file " + fileName)

  ratingEvents, err := readRatingEvents(fileName)
  if err != nil {
    panic(err)
  }

  const topic = "ratings"
  if err := produceRatingEvents(
    topic, 
    producer,
    ratingEvents,
  ); err != nil {
    panic(err)
  }

  const timeout = 10 * time.Second
  fmt.Println("Waiting " + timeout.String() + " until all events get produced.")

  producer.Flush(int(timeout.Milliseconds()))
}

func readRatingEvents(fileName string)( []model.RatingEvent, error ){
  f, err := os.Open(fileName)
  if err != nil {
    return nil, err
  }
  defer f.Close()

  var ratings []model.RatingEvent
  if err := json.NewDecoder(f).Decode(&ratings); err != nil {
    return nil, err
  }
  return ratings, nil
}

func produceRatingEvents(
  topic        string,
  producer     *kafka.Producer,
  ratingEvents []model.RatingEvent,
) error {
  for _, ratingEvent := range ratingEvents {
    encodedEvent, err := json.Marshal(ratingEvent)
    if err != nil {
      return err
    }

    if err := producer.Produce(&kafka.Message{
      TopicPartition: kafka.TopicPartition{
        Topic: &topic,
      },
      Value: encodedEvent,
    }, nil); err != nil {
      return err
    }
  }

  return nil
}
