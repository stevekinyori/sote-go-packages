package TestingNATSFetch

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2023/sMessage"
)

func run() {
	mmPtr, soteErr := sMessage.New(context.Background(), "synadia", "staging", "", "west.eu.geo.ngs.global", "test", false, 1, 250*time.Millisecond,
		false)
	if soteErr.ErrCode != nil {
		panic(soteErr.FmtErrMsg)
	}

	js, errJetStream := mmPtr.NatsConnectionPtr.JetStream()
	if errJetStream != nil {
		panic(errJetStream.Error())
	}

	_ = js.DeleteStream("Delete-Me-Stream")

	// _, soteErr = mmPtr.CreateLimitsStreamWithFileStorage("Delete-Me-Stream", []string{"my-test-subject"}, 1, false)
	// if soteErr.ErrCode != nil {
	// 	panic(soteErr.FmtErrMsg)
	// }

	_, errAddStream := js.AddStream(&nats.StreamConfig{
		Name:     "Delete-Me-Stream",
		Subjects: []string{"my-test-subject"},
	})
	if errAddStream != nil {
		panic(errAddStream.Error())
	}

	// sConsumerConfig := &nats.ConsumerConfig{
	// 	Durable:         "Delete-Me-Durable-Name",
	// 	DeliverSubject:  "",
	// 	DeliverPolicy:   nats.DeliverAllPolicy,
	// 	OptStartSeq:     0,
	// 	OptStartTime:    nil,
	// 	AckPolicy:       nats.AckExplicitPolicy,
	// 	AckWait:         0,
	// 	MaxDeliver:      1,
	// 	FilterSubject:   "my-test-subject",
	// 	ReplayPolicy:    nats.ReplayInstantPolicy,
	// 	RateLimit:       0,
	// 	SampleFrequency: "",
	// 	MaxWaiting:      0,
	// 	MaxAckPending:   0,
	// }

	// _, errConsumer := js.AddConsumer("Delete-Me-Stream", sConsumerConfig)
	// if errConsumer != nil {
	// 	panic(errConsumer.Error())
	// }
	//
	_, errPublish := js.Publish("my-test-subject", []byte("Hello world message"))
	if errPublish != nil {
		panic(errPublish.Error())
	}

	myPullSub, errPullSub := js.PullSubscribe("my-test-subject", "Delete-Me-Durable-Name")
	if errPullSub != nil {
		panic(errPullSub.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	start := time.Now()
	messages, errFetch := myPullSub.Fetch(1, nats.Context(ctx))
	fmt.Printf("Took %v\n", time.Since(start))
	if errFetch != nil {
		panic(errFetch.Error())
	}
	for _, message := range messages {
		fmt.Println(string(message.Data))
	}

	// _ = js.DeleteStream("Delete-Me-Stream")

	mmPtr.NatsConnectionPtr.Close()
}
