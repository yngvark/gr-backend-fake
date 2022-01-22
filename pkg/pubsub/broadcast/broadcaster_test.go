package broadcast_test

import (
	"fmt"
	"testing"

	"github.com/yngvark/gr-zombie/pkg/pubsub/broadcast"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	t.Run("Should send message to listeners", func(t *testing.T) {
		// Given
		var broadcaster *broadcast.Broadcaster = broadcast.New()
		testSubscriber := make(chan string)

		broadcaster.AddSubscriber(testSubscriber)

		// When
		fmt.Println("sending")
		go func() {
			err := broadcaster.BroadCast("YO")
			require.NoError(t, err)
		}()

		// Then
		fmt.Println("receiving")
		lastMsgReceived := <-testSubscriber
		assert.Equal(t, "YO", lastMsgReceived)
	})
}
