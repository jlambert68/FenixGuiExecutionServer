package outgoingPubSubMessages

import (
	"FenixGuiExecutionServer/common_config"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

// Creates a Topic
func createTopicSubscription(topicID string, deadLetteringTopicID string) (err error) {

	// Get the Topic-Subscription-name
	var topicSubscriptionId string
	topicSubscriptionId = createTopicSubscriptionName(topicID)

	ctx := context.Background()

	// Create a new PubSub-client
	var pubSubClient *pubsub.Client
	pubSubClient, err = creatNewPubSubClient(ctx)

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

	// Get the DeadLettering-Topic object if an incoming name was supplied
	var deadLetteringTopic *pubsub.Topic
	var deadLetterPolicy *pubsub.DeadLetterPolicy
	if len(deadLetteringTopicID) > 0 {
		deadLetteringTopic = pubSubClient.Topic(deadLetteringTopicID)

		deadLetterPolicy = &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     deadLetteringTopic.String(),
			MaxDeliveryAttempts: 5,
		}
	}

	// Set up Subscription parameters
	var subscriptionConfig pubsub.SubscriptionConfig
	subscriptionConfig = pubsub.SubscriptionConfig{
		Topic:                 topic,
		PushConfig:            pubsub.PushConfig{},
		BigQueryConfig:        pubsub.BigQueryConfig{},
		CloudStorageConfig:    pubsub.CloudStorageConfig{},
		AckDeadline:           time.Duration(time.Second * 60),
		RetainAckedMessages:   false,
		RetentionDuration:     0,
		ExpirationPolicy:      nil,
		Labels:                nil,
		EnableMessageOrdering: false,
		DeadLetterPolicy:      deadLetterPolicy,
		Filter:                "",
		RetryPolicy: &pubsub.RetryPolicy{
			MinimumBackoff: nil,
			MaximumBackoff: nil,
		},
		Detached:                      false,
		TopicMessageRetentionDuration: 0,
		EnableExactlyOnceDelivery:     true,
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

// Creates a Topic-Subscription-Name
func createTopicSubscriptionName(topicID string) (topicSubscriptionName string) {

	const topicSubscriptionPostfix string = "-sub"

	// Create the Topic-Subscription-name
	topicSubscriptionName = topicID + topicSubscriptionPostfix

	return topicSubscriptionName
}
