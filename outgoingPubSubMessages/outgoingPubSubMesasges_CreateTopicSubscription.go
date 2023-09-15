package outgoingPubSubMessages

import (
	"FenixGuiExecutionServer/common_config"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/sirupsen/logrus"
)

// Creates a Topic
func createTopicSubscription(topicID string) (err error) {

	topicSubscriptionId := "topicID" + "-sub"

	ctx := context.Background()

	// Create a new PubSub-client
	var pubSubClient *pubsub.Client
	err = creatNewPubSubClient(ctx, pubSubClient)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "815eaa22-bbee-47e3-b83f-6374b587e691",
			"err": err,
		}).Error("Got some problem when creating 'pubsub.NewClient'")

		return err
	}

	defer pubSubClient.Close()

	// Get the Topic object
	var topic *pubsub.Topic
	topic = pubSubClient.Topic(topicID)

	// Set up Subscription parameters
	var subscriptionConfig pubsub.SubscriptionConfig
	subscriptionConfig = pubsub.SubscriptionConfig{
		Topic:                         topic,
		PushConfig:                    pubsub.PushConfig{},
		BigQueryConfig:                pubsub.BigQueryConfig{},
		CloudStorageConfig:            pubsub.CloudStorageConfig{},
		AckDeadline:                   0,
		RetainAckedMessages:           false,
		RetentionDuration:             0,
		ExpirationPolicy:              nil,
		Labels:                        nil,
		EnableMessageOrdering:         false,
		DeadLetterPolicy:              nil,
		Filter:                        "",
		RetryPolicy:                   nil,
		Detached:                      false,
		TopicMessageRetentionDuration: 0,
		EnableExactlyOnceDelivery:     false,
		State:                         0,
	}

	// Create a new Topic
	//var pubSubTopic *pubsub.Topic
	_, err = pubSubClient.CreateSubscription(ctx, topicSubscriptionId, subscriptionConfig)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "be22edc9-cfb8-45ff-b751-83c87bef56e4",
			"err": err,
		}).Error("Got some problem when creating a new PubSub Topic-Subscription")

		return err
	}

	return err
}
