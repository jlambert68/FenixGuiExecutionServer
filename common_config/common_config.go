package common_config

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"time"
)

// Logrus debug level

//const LoggingLevel = logrus.DebugLevel
// const LoggingLevel = logrus.InfoLevel

var LoggingLevel logrus.Level

var Logger *logrus.Logger

var highestFenixGuiExecutionServerProtoFileVersion int32 = -1

var gcpAccessToken *oauth2.Token

// ApplicationRunTimeUuid
// Unique 'Uuid' for this running instance. Created at start up. Used as identification
var ApplicationRunTimeUuid string

// ApplicationRunTimeStartUpTime
// Startup time for this running instance
var ApplicationRunTimeStartUpTime time.Time

const ZeroUuid = "00000000-0000-0000-0000-000000000000"
