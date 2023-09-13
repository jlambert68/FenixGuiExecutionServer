package outgoingPubSubMessages

import (
	"FenixGuiExecutionServer/common_config"
	"cloud.google.com/go/pubsub"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"strings"
)

func Publish(msg string) (returnMessageAckNack bool, returnMessageString string, err error) {
	projectID := common_config.GcpProject
	topicID := common_config.TestExecutionStatusPubSubTopicBase

	// Remove any unwanted characters
	// Remove '\n'
	var cleanedMessage string
	cleanedMessage = strings.Replace(msg, "\n", "", -1)

	// Replace '\"' with '"'
	cleanedMessage = strings.ReplaceAll(cleanedMessage, "\\\"", "\"")

	var pubSubClient *pubsub.Client
	var opts []grpc.DialOption

	ctx := context.Background()

	// PubSub is handled within GCP so add TLS
	var creds credentials.TransportCredentials
	creds = credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})

	opts = []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	pubSubClient, err = pubsub.NewClient(ctx, projectID, option.WithGRPCDialOption(opts[0]))

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "19951388-cad6-4f4f-b1d3-a1e8cf758fb4",
			"err": err,
		}).Error("Got some problem when creating 'pubsub.NewClient'")

		return
	}

	defer pubSubClient.Close()

	var pubSubTopic *pubsub.Topic
	var pubSubResult *pubsub.PublishResult
	pubSubTopic = pubSubClient.Topic(topicID)
	pubSubResult = pubSubTopic.Publish(ctx, &pubsub.Message{
		Data: []byte(cleanedMessage),
	})
	// Block until the pubSubResult is returned and a server-generated
	// ID is returned for the published message.
	id, err := pubSubResult.Get(ctx)
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"ID":  "dc8bb67a-2caf-4a46-8a5c-598e253515c5",
			"msg": msg,
		}).Error(fmt.Errorf("pubsub: pubSubResult.Get: %w", err))

		return false, "", err

	}

	common_config.Logger.WithFields(logrus.Fields{
		"ID": "8da81faa-a2a9-4130-83c8-e90b8fbbb955",
		//"token": token,
	}).Debug(fmt.Sprintf("Published a message; msg ID: %v", id))

	return true, "", err
}
