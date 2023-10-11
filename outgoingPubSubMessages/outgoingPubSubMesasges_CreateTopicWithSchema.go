package outgoingPubSubMessages

import (
	"FenixGuiExecutionServer/common_config"
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

// Creates a Topic
func createTopicWithSchema(topicID string, topicSchemaId string) (createdTopic *pubsub.Topic, err error) {

	// Set context to be used
	var ctx context.Context
	ctx = context.Background()

	// Create a new PubSub-client
	var pubSubClient *pubsub.Client
	pubSubClient, err = creatNewPubSubClient(ctx)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":           "b5c955cb-2b2b-47e0-a908-1294da40c930",
			"err":          err,
			"pubSubClient": pubSubClient,
		}).Error("Got some problem when creating 'pubsub.NewClient'")

		return nil, err
	}

	if pubSubClient == nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":           "43a196d4-182f-4899-a96a-9e2947da7b79",
			"pubSubClient": pubSubClient,
		}).Error("Got some problem when creating 'pubsub.NewClient'")

		return nil, errors.New("got some problem when creating 'pubsub.NewClient'")
	}

	defer pubSubClient.Close()

	// Create Schema long name
	var schemaLongName string
	schemaLongName = fmt.Sprintf("projects/%s/schemas/%s", gcpProject, topicSchemaId)

	// Create the Topic Config
	var topicConfig *pubsub.TopicConfig
	topicConfig = &pubsub.TopicConfig{
		Labels: nil,
		MessageStoragePolicy: pubsub.MessageStoragePolicy{
			AllowedPersistenceRegions: nil,
		},
		KMSKeyName: "",
		SchemaSettings: &pubsub.SchemaSettings{
			Schema:          schemaLongName,
			Encoding:        pubsub.EncodingJSON,
			FirstRevisionID: "",
			LastRevisionID:  "",
		},
		RetentionDuration: nil,
	}

	// Create a new Topic
	//var pubSubTopic *pubsub.Topic
	createdTopic, err = pubSubClient.CreateTopicWithConfig(ctx, topicID, topicConfig)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "437a8d8a-6e84-4542-8fad-5ff9be240d4a",
			"err": err,
		}).Error("Got some problem when creating a new PubSub Topic")

		return nil, err
	}

	return createdTopic, err
}
