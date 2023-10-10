package outgoingPubSubMessages

import (
	"cloud.google.com/go/pubsub"
)

// CreateTopicDeadLettingAndSubscriptionIfNotExists
// Create Topic, TopicSubscription, DeadLetteringTopic and DeadLetteringTopicSubscription based on 'Topic-name'
func CreateTopicDeadLettingAndSubscriptionIfNotExists(pubSubTopicToVerify string) (err error) {

	var foundDeadLettingTopics *pubsub.Topic
	var pubSubTopics []*pubsub.Topic
	var pubSubTopicSubscriptions []*pubsub.Subscription
	var pubSubDeadLetteringTopicSubscriptions []*pubsub.Subscription
	var topicExists bool
	var deadLetteringTopicExists bool
	var topicSubscriptionExists bool
	var deadLetteringTopicSubscriptionExists bool

	// Create DeadLetteringTopic-name
	var pubSubDeadLetteringTopicToVerify string
	pubSubDeadLetteringTopicToVerify = createDeadLetteringTopicName(pubSubTopicToVerify)

	// Get all topics
	pubSubTopics, err = listTopics()

	// Loop the slice with topics to find out if Topics already exists
	for _, tempTopic := range pubSubTopics {

		// Look if the Topic was found
		if tempTopic.ID() == pubSubTopicToVerify {
			topicExists = true

		}
		// If the DeadLettingTopic was found
		if tempTopic.ID() == pubSubDeadLetteringTopicToVerify {
			deadLetteringTopicExists = true
			foundDeadLettingTopics = tempTopic

		}

		// If both Topic and DeadLettingTopic were found then exit for-loop
		if topicExists && deadLetteringTopicExists {
			break
		}
	}

	// Create Subscription Name
	var topicSubscriptionNameToVerify string
	topicSubscriptionNameToVerify = createTopicSubscriptionName(pubSubTopicToVerify)

	// Create DeadLettering-Subscription Name
	var topicDeadLetteringSubscriptionNameToVerify string
	topicDeadLetteringSubscriptionNameToVerify = createDeadLetteringTopicSubscriptionName(pubSubTopicToVerify)

	// Get all topic-subscription when the Topic existed
	if topicExists == true {
		pubSubTopicSubscriptions, err = listSubscriptions(pubSubTopicToVerify)

		// Loop the slice with topic-subscriptions to find out if subscriptions already exists
		for _, tempTopicSubscription := range pubSubTopicSubscriptions {

			// If the TopicSubscription was found then exit for loop
			if tempTopicSubscription.ID() == topicSubscriptionNameToVerify {
				topicSubscriptionExists = true
				break
			}
		}
	}

	// Get all DeadLettering-topic-subscriptions when the Topic existed
	if deadLetteringTopicExists == true {
		pubSubDeadLetteringTopicSubscriptions, err = listSubscriptions(pubSubDeadLetteringTopicToVerify)

		// Loop the slice with DeadLettering-topic-subscriptions to find out if subscriptions already exists
		for _, tempDeadLetterTopicSubscription := range pubSubDeadLetteringTopicSubscriptions {

			// If the DeadLetteringTopicSubscription was found then exit for loop
			if tempDeadLetterTopicSubscription.ID() == topicDeadLetteringSubscriptionNameToVerify {
				deadLetteringTopicSubscriptionExists = true
				break
			}
		}
	}

	// if the Topic was not found then create the Topic
	if topicExists == false {
		_, err = createTopic(pubSubTopicToVerify)
		if err != nil {
			return err
		}
	}

	// if the DeadLettingTopic was not found then create the Topic
	if deadLetteringTopicExists == false {
		foundDeadLettingTopics, err = createTopicDeadLettering(pubSubTopicToVerify)
		if err != nil {
			return err
		}
	}

	// if the TopicSubscription was not found then create the TopicSubscription
	if topicSubscriptionExists == false {
		err = createTopicSubscription(pubSubTopicToVerify, foundDeadLettingTopics.ID())
		if err != nil {
			return err
		}
	}

	// if the DeadLetteringTopicSubscription was not found then create the DeadLetteringTopicSubscription
	if deadLetteringTopicSubscriptionExists == false {
		err = createDeadLetteringTopicSubscription(pubSubTopicToVerify)
		if err != nil {
			return err
		}
	}

	return err
}
