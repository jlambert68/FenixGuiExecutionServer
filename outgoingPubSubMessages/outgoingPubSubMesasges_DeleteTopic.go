package outgoingPubSubMessages

import (
	"FenixGuiExecutionServer/common_config"
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
)

// Delete a Topic
func deleteTopic(topicID string) (err error) {

	ctx := context.Background()

	// Create a new PubSub-client
	var pubSubClient *pubsub.Client
	pubSubClient, err = creatNewPubSubClient(ctx)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "b5c955cb-2b2b-47e0-a908-1294da40c930",
			"err": err,
		}).Error("Got some problem when creating 'pubsub.NewClient'")

		return err
	}

	defer pubSubClient.Close()

	// Get PubSub Topic
	var pubSubTopic *pubsub.Topic
	pubSubTopic = pubSubClient.Topic(topicID)
	if pubSubTopic == nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID": "82207085-7229-4824-bf05-b10a34ab1a0f",
		}).Error("Got some problem when getting PubSub Topic")

		return errors.New("got some problem when getting PubSub Topic")
	}

	// Delete the Topic
	err = pubSubTopic.Delete(ctx)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID": "3430202d-496e-4c75-9c6a-3a67f79b6855",
		}).Error("Got some problem when deleting PubSub Topic")

		return err
	}

	return err
}
