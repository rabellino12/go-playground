package ioclient

import (
	"errors"

	centrifuge "github.com/centrifugal/centrifuge-go"
)

// Handler is the match handler interface for the subscriber
type Handler interface {
	GetID() string
	OnPublish(*centrifuge.Subscription, centrifuge.PublishEvent)
}

// ListenMatch lobby io controller
func ListenMatch(c *centrifuge.Client, matchHandler Handler) error {
	id := matchHandler.GetID()
	sub, err := c.NewSubscription("$match:" + id)
	if err != nil {
		return errors.New("couldn't subscribe to " + id)
	}

	sub.OnPublish(matchHandler)
	sub.Subscribe()
	return nil
}
