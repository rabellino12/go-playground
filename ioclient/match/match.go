package ioclient

import (
	"fmt"

	centrifuge "github.com/centrifugal/centrifuge-go"
)

// Handler is the match handler interface for the subscriber
type Handler interface {
	GetID() string
	OnPublish(*centrifuge.Subscription, centrifuge.PublishEvent)
}

// ListenMatch lobby io controller
func ListenMatch(c *centrifuge.Client, matchHandler Handler) {
	id := matchHandler.GetID()
	fmt.Println("Listening to match " + id)
	sub, _ := c.NewSubscription("$match:" + id)
	sub.OnPublish(matchHandler)
	err := sub.Subscribe()
	if err != nil {
		fmt.Println(err.Error())
	}
	select {}
}
