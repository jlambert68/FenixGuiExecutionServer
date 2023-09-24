package outgoingPubSubMessages

import (
	"fmt"
	"log"
	"os"
)

func MyTestPubSubFunctions() {

	fmt.Println("START")

	cloudProject := "mycloud-run-project"

	InitiatePubSubFunctionality(cloudProject)

	//_, _, _ = Publish("msg string")

	// -- createTopic 'myTestTopic'
	fmt.Println("createTopic")
	testTopic := "myTestTopic"
	err := createTopic(testTopic)
	if err != nil {
		log.Println("'createTopic' - err: %s ", err.Error())
		os.Exit(0)
	}
	fmt.Println("createTopic - SUCCESS")

	// -- createTopic 'myTestTopic-DeadLettering'
	fmt.Println("createTopic")
	testTopicDeadLettering := "myTestTopic-DeadLettering"
	err = createTopic(testTopicDeadLettering)
	if err != nil {
		log.Println("'createTopic' - err: %s ", err.Error())
		os.Exit(0)
	}
	fmt.Println("createTopic - SUCCESS")

	// -- createTopic
	fmt.Println("createTopic")
	testTopic = "myTestTopic"
	err = createTopicSubscription(testTopic, testTopicDeadLettering)
	if err != nil {
		log.Println("'createTopicSubscription' - err: %s ", err.Error())
		os.Exit(0)
	}
	fmt.Println("createTopicSubscription - SUCCESS")

	fmt.Println("FINISH")
	os.Exit(0)
}
