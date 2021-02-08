package kafka

import (
  "bcs/logcatdog/linebyline"
  "encoding/json"
  "gopkg.in/Shopify/sarama.v1"
  "log"
  "sync"
)


type Producer struct {
  DataStream chan linebyline.LineData
  done  chan struct{}
}

func NewProducer(dataStream chan linebyline.LineData) Producer {
  producer := new(Producer)

  producer.done = make(chan struct{})
  producer.DataStream = dataStream

  return *producer
}

func (p *Producer) Send() (chan struct{}, error)  {

  go func() {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true

    producer, err := sarama.NewAsyncProducer([]string{"sparksb6-wsl2.cdr0.net:9092"}, config)
    if err != nil {
      log.Fatal(err)
    }

    //signals := make(chan os.Signal, 1)
    //signal.Notify(signals, os.Interrupt)

    var (
      wg  sync.WaitGroup
      enqueued int
      successes int
      pErrors int
    )

    wg.Add(1)
    go func() {
      defer wg.Done()
      for s := range producer.Successes() {
        successes += 1
        log.Printf("Success: %v\n", s)
      }
    }()

    wg.Add(1)
    go func() {
      defer wg.Done()
      for err := range producer.Errors() {
        log.Println(err)
        pErrors += 1
      }
    }()

    for lineData := range p.DataStream {
      b, err := json.Marshal(lineData)
      if err != nil {
        log.Fatal(err)
      }

      message := &sarama.ProducerMessage{Topic: "logcattopic", Value: sarama.StringEncoder(b)}

      producer.Input() <- message
      enqueued += 1
    }

    producer.AsyncClose()

    wg.Wait()

    p.done <- struct{}{}

  }()

  return p.done, nil
}



