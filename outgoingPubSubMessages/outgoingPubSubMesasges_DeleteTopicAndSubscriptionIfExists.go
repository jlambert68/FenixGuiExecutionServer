package outgoingPubSubMessages

import (
	"cloud.google.com/go/pubsub"
	"context"
)

// DeleteTopicDeadLettingAndSubscriptionIfNotExists
// Delete Topic, TopicSubscription, DeadLetteringTopic and DeadLetteringTopicSubscription for 'Topic-name'
func DeleteTopicDeadLettingAndSubscriptionIfNotExists(pubSubTopicToDelete string) (err error) {

	var pubSubTopics []*pubsub.Topic
	var pubSubTopicSubscriptions []*pubsub.Subscription
	var pubSubDeadLetteringTopicSubscriptions []*pubsub.Subscription
	var topicExists bool
	var deadLetteringTopicExists bool

	ctx := context.Background()

	// Create DeadLetteringTopic-name
	var pubSubDeadLetteringTopicToDelete string
	pubSubDeadLetteringTopicToDelete = createDeadLetteringTopicName(pubSubTopicToDelete)

	// Get all topics
	pubSubTopics, err = listTopics()

	// Loop the slice with topics to find out if Topics already exists
	for _, tempTopic := range pubSubTopics {

		// Look if the Topic was found
		if tempTopic.ID() == pubSubTopicToDelete {
			topicExists = true

		}
		// If the DeadLettingTopic was found
		if tempTopic.ID() == pubSubDeadLetteringTopicToDelete {
			deadLetteringTopicExists = true

		}

		// If both Topic and DeadLettingTopic were found then exit for-loop
		if topicExists && deadLetteringTopicExists {
			break
		}
	}

	// Get all topic-subscriptions
	pubSubTopicSubscriptions, err = listSubscriptions(pubSubTopicToDelete)

	// Loop the slice with topic-subscriptions and delete them
	for _, tempTopicSubscription := range pubSubTopicSubscriptions {
		tempTopicSubscription.Delete(ctx)
	}

	// Get all DeadLettering-topic-subscriptions
	pubSubDeadLetteringTopicSubscriptions, err = listSubscriptions(pubSubDeadLetteringTopicToDelete)

	// Loop the slice with DeadLettering-topic-subscriptions and delete them
	for _, tempDeadLetteringTopicSubscription := range pubSubDeadLetteringTopicSubscriptions {
		tempDeadLetteringTopicSubscription.Delete(ctx)
	}

	// if the Topic was not found then delete the Topic
	if topicExists == false {
		err = deleteTopic(pubSubTopicToDelete)
		if err != nil {
			return err
		}
	}

	// if the DeadLettingTopic was not found then delete the Topic
	if deadLetteringTopicExists == false {
		err = deleteTopic(pubSubDeadLetteringTopicToDelete)
		if err != nil {
			return err
		}
	}

	return err
}
