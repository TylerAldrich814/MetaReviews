package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/TylerAldrich814/MetaReviews/rating/pkg/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Ingester defines a Kafka Message Ingester.
type Ingester struct {
  consumer  *kafka.Consumer
  topic     string
}

// New creates a new Kafka Ingester.
func New(
  addr    string,
  groupID string,
  topic   string,
)( *Ingester,error ){
  consumer, err := kafka.NewConsumer(
    &kafka.ConfigMap{
      "bootstrap.servers" : addr,
      "group.id"          : groupID,
      "auto.offset.reset" : "earliest",
    },
  )
  if err != nil {
    return nil, err
  }

  return &Ingester{ consumer,topic },nil
}

// Ingest starts ingestion from Kafka and retuns a RatingEvent Channel.
// Which represents the data consumed from the topic.
func(i *Ingester) Ingest(
  ctx  context.Context,
)( chan model.RatingEvent,error ){
  if err := i.consumer.SubscribeTopics(
    []string{i.topic},
    nil,
  ); err != nil {
    return nil, err
  }
  ch := make(chan model.RatingEvent, 1)
  go func(){
    for {
      select {
      case <-ctx.Done():
        close(ch)
        i.consumer.Close()
      default:
      }
      msg, err := i.consumer.ReadMessage(-1)
      if err != nil {
        log.Printf("Consumer Error: %s\n", err.Error())
        continue
      }
      var event model.RatingEvent
      if err := json.Unmarshal(msg.Value, &event); err != nil {
        log.Printf("Unmarshal Error: %s\n", err.Error())
        continue
      }
      ch<-event
    }
  }()

  return ch,nil
}
