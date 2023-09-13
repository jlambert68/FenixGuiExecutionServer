package outgoingPubSubMessages

import (
	"FenixGuiExecutionServer/common_config"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

func publishMessage(topicID string, msg string) (err error) {

	ctx := context.Background()

	// Create a new PubSub-client
	var pubSubClient *pubsub.Client
	err = creatNewPubSubClient(ctx, pubSubClient)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "358e33c0-f993-4ce6-95e1-538bd14c466b",
			"err": err,
		}).Error("Got some problem when creating 'pubsub.NewClient'")

		return err
	}

	defer pubSubClient.Close()

	var pubSubTopic *pubsub.Topic
	var pubSubResult *pubsub.PublishResult
	pubSubTopic = pubSubClient.Topic(topicID)
	pubSubResult = pubSubTopic.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the pubSubResult is returned and a server-generated
	// ID is returned for the published message.
	var messageId string
	messageId, err = pubSubResult.Get(ctx)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":        "48b9a5cf-b76a-4ddc-aa2e-97c2d0126ca8",
			"msg":       msg,
			"messageId": messageId,
		}).Error(fmt.Errorf("pubsub: pubSubResult.Get: %w", err))

		return err

	}

	return err
}
