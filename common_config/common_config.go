package common_config

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Logrus debug level

//const LoggingLevel = logrus.DebugLevel
//const LoggingLevel = logrus.InfoLevel
const LoggingLevel = logrus.DebugLevel // InfoLevel

var Logger *logrus.Logger

var highestFenixGuiExecutionServerProtoFileVersion int32 = -1

var gcpAccessToken *oauth2.Token

// Unique 'Uuid' for this running instance. Created at start up. Used as identification
var ApplicationRunTimeUuid string

const ZeroUuid = "00000000-0000-0000-0000-000000000000"
