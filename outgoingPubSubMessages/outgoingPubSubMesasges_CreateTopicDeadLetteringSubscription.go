package outgoingPubSubMessages

// Creates a DeadLettering-TopicSubscription
func createDeadLetteringTopicSubscription(topicID string) (err error) {

	// Create the Topic-name for DeadLettering
	var deadLetteringTopicName string
	deadLetteringTopicName = createDeadLetteringTopicName(topicID)

	// Create the DeadLettingTopic
	err = createTopicSubscription(deadLetteringTopicName, "")

	return err
}

// Creates a DeadLettering-Topic-Subscription-Name
func createDeadLetteringTopicSubscriptionName(topicID string) (deadLetteringTopicSubscriptionName string) {

	// Create The DeadLettering-Name for the Topic
	var deadLetteringTopicName string
	deadLetteringTopicName = createDeadLetteringTopicName(topicID)

	// Create the DeadLettering-Topic-Subscription-name
	deadLetteringTopicSubscriptionName = createTopicSubscriptionName(deadLetteringTopicName)

	return deadLetteringTopicSubscriptionName
}
