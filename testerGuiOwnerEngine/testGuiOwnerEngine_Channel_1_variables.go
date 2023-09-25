package testerGuiOwnerEngine

import (
	"FenixGuiExecutionServer/common_config"
	"time"
)

// SomeoneIsClosingDownStruct
// The following message is sent over Postgres Broadcast system, TesterGuiOwner-Channel 1
// Used to specify an GuiExecutionServer or a TesterGui that is Closing Down
type SomeoneIsClosingDownStruct struct {
	WhoISClosingDown common_config.WhoISClosingDownType `json:"whoisclosingdown"`
	ApplicationId    string                             `json:"applicationid"`
	UserId           string                             `json:"userid"`
	MessageTimeStamp time.Time                          `json:"messagetimestamp"`
}
