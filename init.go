package main

import (
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/messagesToExecutionServer"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strconv"
)

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func init() {
	//executionLocationForGuiExecutionServer := flag.String("startupType", "0", "The application should be started with one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
	//flag.Parse()

	var err error

	// Where is GuiExecutionServer started
	var executionLocationForGuiExecutionServer = mustGetenv("ExecutionLocationForFenixGuiExecutionServer")

	switch executionLocationForGuiExecutionServer {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForFenixGuiExecutionServer = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForFenixGuiExecutionServer = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForFenixGuiExecutionServer = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for FenixGuiExecutionServer: " + executionLocationForGuiExecutionServer + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	// Address to GuiExecutionServer
	common_config.FenixGuiExecutionServerAddress = mustGetenv("FenixGuiExecutionServerAddress")

	// Port for GuiExecutionServer
	common_config.FenixGuiExecutionServerPort, err = strconv.Atoi(mustGetenv("FenixGuiExecutionServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixGuiExecutionServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Where is ExecutionServer started
	var executionLocationForExecutionServer = mustGetenv("ExecutionLocationForFenixExecutionServer")

	switch executionLocationForExecutionServer {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForFenixExecutionServer = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForFenixExecutionServer = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForFenixExecutionServer = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for FenixExecutionServer: " + executionLocationForExecutionServer + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	// Address to ExecutionServer
	common_config.FenixExecutionServerAddress = mustGetenv("FenixExecutionServerAddress")

	// Port for ExecutionServer
	common_config.FenixExecutionServerPort, err = strconv.Atoi(mustGetenv("FenixExecutionServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixExecutionServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Save the Dial-string to use for connecting to ExecutionServer
	messagesToExecutionServer.FenixExecutionServerAddressToDial = common_config.FenixExecutionServerAddress + ":" + strconv.Itoa(common_config.FenixExecutionServerPort)

	// Save the address to use for getting access token
	messagesToExecutionServer.FenixExecutionServerAddressToUse = common_config.FenixExecutionServerAddress

	// Extract Debug level
	var loggingLevel = mustGetenv("LoggingLevel")

	switch loggingLevel {

	case "DebugLevel":
		common_config.LoggingLevel = logrus.DebugLevel

	case "InfoLevel":
		common_config.LoggingLevel = logrus.InfoLevel

	default:
		fmt.Println("Unknown loggingLevel '" + loggingLevel + "'. Expected one of the following: 'DebugLevel', 'InfoLevel'")
		os.Exit(0)

	}

	fmt.Printf("%s", common_config.LoggingLevel)

	// Extract OAuth 2.0 Client ID
	common_config.AuthClientId = mustGetenv("AuthClientId")

	// Extract OAuth 2.0 Client Secret
	common_config.AuthClientSecret = mustGetenv("AuthClientSecret")

	// Should all SQL-queries be logged before executed
	var tempBoolAsString string
	var tempBool bool
	tempBoolAsString = mustGetenv("LogAllSQLs")
	tempBool, err = strconv.ParseBool(tempBoolAsString)
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'LogAllSQLs' to a boolean, error: ", err)
		os.Exit(0)
	}
	common_config.LogAllSQLs = tempBool

	// Max number of DB-connection from Pool. Not stored because it is re-read when connecting the DB-pool
	_ = mustGetenv("DB_POOL_MAX_CONNECTIONS")

	_, err = strconv.Atoi(mustGetenv("DB_POOL_MAX_CONNECTIONS"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'DB_POOL_MAX_CONNECTIONS' to an integer, error: ", err)
		os.Exit(0)

	}

	// Extract the GCP-project
	common_config.GcpProject = mustGetenv("GcpProject")

	// Should PubSub be used for sending 'TestExecutionsStatus' to TesterGui
	common_config.UsePubSubWhenSendingExecutionStatus, err = strconv.ParseBool(mustGetenv("UsePubSubWhenSendingExecutionStatus"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'UsePubSubWhenSendingExecutionStatus' to a boolean, error: ", err)
		os.Exit(0)
	}

	// Extract PubSub-Topic-base for where to send the 'TestExecutionsStatus'
	common_config.TestExecutionStatusPubSubTopicBase = mustGetenv("TestExecutionStatusPubSubTopicBase")

	// Extract local path to Service-Account file
	common_config.LocalServiceAccountPath = mustGetenv("LocalServiceAccountPath")
	// The only way have an OK space is to replace an existing character
	if common_config.LocalServiceAccountPath == "#" {
		common_config.LocalServiceAccountPath = ""
	}

	// Extract the topic-schema name to be used when sending 'TestExecutionsStatus' to TesterGui
	common_config.TestExecutionStatusPubSubTopicSchema = mustGetenv("TestExecutionStatusPubSubTopicSchema")

	// Set the environment variable that Google-client-libraries look for
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", common_config.LocalServiceAccountPath)

}
