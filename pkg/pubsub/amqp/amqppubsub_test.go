/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package amqp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	dctest "github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	"github.com/trustbloc/orb/pkg/lifecycle"
	"github.com/trustbloc/orb/pkg/pubsub/spi"
)

const (
	dockerImage = "rabbitmq"
	dockerTag   = "3.8.16"
)

func TestAMQP(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		const topic = "some-topic"

		p := New(Config{URI: "amqp://guest:guest@localhost:5672/"})
		require.NotNil(t, p)

		msgChan, err := p.Subscribe(context.Background(), topic)
		require.NoError(t, err)

		msg := message.NewMessage(watermill.NewUUID(), []byte("some payload"))
		require.NoError(t, p.Publish(topic, msg))

		select {
		case m := <-msgChan:
			require.Equal(t, msg.UUID, m.UUID)
		case <-time.After(200 * time.Millisecond):
			t.Fatal("timed out waiting for message")
		}

		require.NoError(t, p.Close())

		_, err = p.Subscribe(context.Background(), topic)
		require.True(t, errors.Is(err, lifecycle.ErrNotStarted))
		require.True(t, errors.Is(p.Publish(topic, msg), lifecycle.ErrNotStarted))
	})

	t.Run("Connection failure", func(t *testing.T) {
		require.Panics(t, func() {
			p := New(Config{URI: "amqp://guest:guest@localhost:9999/", MaxConnectRetries: 3})
			require.NotNil(t, p)
		})
	})

	t.Run("Pooled subscriber -> success", func(t *testing.T) {
		const (
			n     = 100
			topic = "pooled"
		)

		publishedMessages := &sync.Map{}
		receivedMessages := &sync.Map{}

		p := New(Config{
			URI: "amqp://guest:guest@localhost:5672/",
		})
		require.NotNil(t, p)
		defer func() {
			require.NoError(t, p.Close())
		}()

		msgChan, err := p.SubscribeWithOpts(context.Background(), topic, spi.WithPool(10))
		require.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(n)

		go func(msgChan <-chan *message.Message) {
			for m := range msgChan {
				go func(msg *message.Message) {
					receivedMessages.Store(msg.UUID, msg)

					// Add a delay to simulate processing.
					time.Sleep(100 * time.Millisecond)

					msg.Ack()

					wg.Done()
				}(m)
			}
		}(msgChan)

		for i := 0; i < n; i++ {
			go func() {
				msg := message.NewMessage(watermill.NewUUID(), []byte("some payload"))
				publishedMessages.Store(msg.UUID, msg)

				require.NoError(t, p.Publish(topic, msg))
			}()
		}

		wg.Wait()

		publishedMessages.Range(func(msgID, _ interface{}) bool {
			_, ok := receivedMessages.Load(msgID)
			require.Truef(t, ok, "message not received: %s", msgID)

			return true
		})
	})
}

func TestMain(m *testing.M) {
	code := 1

	defer func() { os.Exit(code) }()

	pool, err := dctest.NewPool("")
	if err != nil {
		panic(fmt.Sprintf("pool: %v", err))
	}

	resource, err := pool.RunWithOptions(&dctest.RunOptions{
		Repository: dockerImage,
		Tag:        dockerTag,
		PortBindings: map[dc.Port][]dc.PortBinding{
			"5672/tcp": {{HostIP: "", HostPort: "5672"}},
		},
	})
	if err != nil {
		logger.Errorf(`Failed to start RabbitMQ Docker image.`)

		panic(fmt.Sprintf("run with options: %v", err))
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			panic(fmt.Sprintf("purge: %v", err))
		}
	}()

	code = m.Run()
}
