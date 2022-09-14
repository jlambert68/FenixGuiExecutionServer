package main

import "time"

// All TestInstructions in CloudDB
//var cloudDBTestInstructionItems []cloudDBTestInstructionItem

// Type for holding one TestInstructionItem that is stored in CloudDB
type cloudDBTestInstructionItemStruct struct {
	domainUuid                   string    // The Domain, UUID, where the system resides
	domainName                   string    // The Domain, Name, where the system resides
	testInstructionUuid          string    // TestInstruction UUID
	testInstructionName          string    // TestInstruction Name
	testInstructionTypeUuid      string    // The Type(Group), Uuid, of TestInstruction
	testInstructionTypeName      string    // The Type(Group), Name, of TestInstruction
	testInstructionDescription   string    // The description of the TestInstruction
	testInstructionMouseOverText string    // The mouse over text when hovering over TestInstruction
	deprecated                   bool      // Indicates that this TestInstruction shouldn't be used anymore
	enabled                      bool      // TestInstruction can be disabled when the user shouldn't use it anymore
	majorVersionNumber           uint32    // Change in Major Version Number means that user must act on change
	minorVersionNumber           uint32    // Change in Minor Version Number means that user must NOT act on change
	updatedTimeStamp             time.Time // The TimeStamp when the TestInstruction was last updated
	updatedTimeStampAsString     string    // The TimeStamp, as a string, when the TestInstruction was last updated
}

// All Pre-Created TestInstructionsContainers in CloudDB
//var cloudDBTestInstructionsContainerItems []cloudDBTestInstructionContainerItem

// Type for holding one TestInstructionContainerItem that is stored in CloudDB
type cloudDBTestInstructionContainerItemStruct struct {
	domainUuid                            string    // The Domain, UUID, where the system resides
	domainName                            string    // The Domain, Name, where the system resides
	testInstructionContainerUuid          string    // TestInstructionContainer UUID
	testInstructionContainerName          string    // TestInstructionContainer Name
	testInstructionContainerTypeUuid      string    // The Type(Group), Uuid, of TestInstructionContainer
	testInstructionContainerTypeName      string    // The Type(Group), Name, of TestInstructionContainers
	testInstructionContainerDescription   string    // The description of the TestInstructionContainer
	testInstructionContainerMouseOverText string    // The mouse over text when hovering over TestInstructionContainer
	deprecated                            bool      // Indicates that this TestInstruction shouldn't be used anymore
	enabled                               bool      // TestInstruction can be disabled when the user shouldn't use it anymore
	majorVersionNumber                    uint32    // Change in Major Version Number means that user must act on change
	minorVersionNumber                    uint32    // Change in Minor Version Number means that user must NOT act on change
	updatedTimeStamp                      time.Time // The TimeStamp when the TestInstructionContainer was last updated
	updatedTimeStampAsString              string    // The TimeStamp, as a string, when the TestInstructionContainer was last updated
}
