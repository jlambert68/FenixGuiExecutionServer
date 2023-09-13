package outgoingPubSubMessages

import (
	"FenixGuiExecutionServer/common_config"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/sirupsen/logrus"
)

// Creates a Topic
func createTopic(topicID string) (err error) {

	ctx := context.Background()

	// Create a new PubSub-client
	var pubSubClient *pubsub.Client
	err = creatNewPubSubClient(ctx, pubSubClient)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "b5c955cb-2b2b-47e0-a908-1294da40c930",
			"err": err,
		}).Error("Got some problem when creating 'pubsub.NewClient'")

		return err
	}

	defer pubSubClient.Close()

	// Create a new Topic
	//var pubSubTopic *pubsub.Topic
	_, err = pubSubClient.CreateTopic(ctx, topicID)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "1ce8e7f5-bbf6-4c9e-9e52-04d292ae0147",
			"err": err,
		}).Error("Got some problem when creating a new PubSub Topic")

		return err
	}

	return err
}
