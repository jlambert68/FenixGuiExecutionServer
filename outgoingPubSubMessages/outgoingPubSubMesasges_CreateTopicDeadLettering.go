package outgoingPubSubMessages

import "cloud.google.com/go/pubsub"

// Creates a Topic
func createTopicDeadLettering(topicID string) (createdTopic *pubsub.Topic, err error) {

	// Get the Topic-name for DeadLettering
	var topicDeadLetteringName string
	topicDeadLetteringName = createDeadLetteringTopicName(topicID)

	// Create the DeadLettingTopic
	createdTopic, err = createTopic(topicDeadLetteringName)

	return createdTopic, err
}

// Creates a DeadLettering-TopicName
func createDeadLetteringTopicName(topicID string) (deadLetteringTopicName string) {

	const deadLetteringTopicPostfix string = "-DeadLettering"

	// Create the Topic-name for DeadLettering
	deadLetteringTopicName = topicID + deadLetteringTopicPostfix

	return deadLetteringTopicName
}
