package kafka

import (
  lbl "bcs/logcatdog/linebyline"
  "encoding/json"
  "fmt"
  "gopkg.in/Shopify/sarama.v1"
  "log"
  "sync"
)

func Send(done <-chan struct{}, in <-chan lbl.LineData, topic string)  {
  wgOuter := sync.WaitGroup{}

  wgOuter.Add(1)
  go func() {
    //defer wgOuter.Done()

    config := sarama.NewConfig()
    config.Producer.Return.Successes = true

    //pUrl := "sparksb6-wsl2.cdr0.net:9092"
    pUrl := "127.0.0.1:9092"
    producer, err := sarama.NewAsyncProducer([]string{pUrl}, config)
    if err != nil {
      fmt.Printf("Is the IP of url (%s) right?", pUrl)
      log.Fatal(err)
    }

    var (
      wg  sync.WaitGroup
      enqueued int
      successes int
      pErrors int
    )

    wg.Add(1)
    go func() {
      defer wg.Done()
      for _ = range producer.Successes() {
        successes += 1
        //log.Printf("Success: %v\n", s)
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

    for lineData := range in {
      b, err := json.Marshal(lineData)
      if err != nil {
        log.Fatal(err)
      }
      //fmt.Printf("line: %s\n", b)

      message := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(b)}

      producer.Input() <- message
      enqueued += 1
    }

    producer.AsyncClose()

    wg.Wait()

    fmt.Printf("enqueued: %d, successes: %d, errors: %d\n", enqueued, successes, pErrors)

    wgOuter.Done()
  }()

  wgOuter.Wait()
}
