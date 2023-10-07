package outgoingPubSubMessages

import "cloud.google.com/go/pubsub"

// Creates a Topic
func createTopicDeadLettering(topicID string) (createdTopic *pubsub.Topic, err error) {

	const deadLetteringTopicPostfix string = "-DeadLettering"

	// Create the Topic-name for DeadLettering
	var topicDeadLetteringName string
	topicDeadLetteringName = topicID + deadLetteringTopicPostfix

	// Create the DeadLettingTopic
	createdTopic, err = createTopic(topicDeadLetteringName)

	return createdTopic, err
}
