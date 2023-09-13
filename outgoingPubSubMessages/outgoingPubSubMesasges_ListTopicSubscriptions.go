package outgoingPubSubMessages

import (
	"FenixGuiExecutionServer/common_config"
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

// List Topics
func listSubscriptions(topicID string) (pubSubTopicSubscriptions []*pubsub.Subscription, err error) {

	ctx := context.Background()

	// Create a new PubSub-client
	var pubSubClient *pubsub.Client
	err = creatNewPubSubClient(ctx, pubSubClient)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "540c6a19-f205-4dc2-a422-fa2a72a46a50",
			"err": err,
		}).Error("Got some problem when creating 'pubsub.NewClient'")

		return nil, err
	}

	defer pubSubClient.Close()

	// Get TopicsSubscriptions
	var pubSubscriptionIterator *pubsub.SubscriptionIterator
	pubSubscriptionIterator = pubSubClient.Topic(topicID).Subscriptions(ctx)
	for {
		var pubSubTopicSubscription *pubsub.Subscription
		pubSubTopicSubscription, err = pubSubscriptionIterator.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"ID":  "c11524c6-0b6b-4525-a255-3f7dd9d87e4e",
				"err": err,
			}).Error("Got some problem iterating the topic-subscription-response")

			return nil, err
		}
		pubSubTopicSubscriptions = append(pubSubTopicSubscriptions, pubSubTopicSubscription)
	}

	return pubSubTopicSubscriptions, err

}
