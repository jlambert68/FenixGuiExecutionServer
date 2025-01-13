package main

import (
	"FenixGuiExecutionServer/common_config"
	"context"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

// PrepareLoadUsersDomains
// Do initial preparations to be able to load all domains for a specific user
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) PrepareLoadUsersDomains(
	txn pgx.Tx,
	gCPAuthenticatedUser string) (
	domainAndAuthorizations []DomainAndAuthorizationsStruct,
	err error) {

	// Concatenate Users specific Domains and Domains open for every one to use
	domainAndAuthorizations, err = fenixGuiTestCaseBuilderServerObject.concatenateUsersDomainsAndDomainOpenToEveryOneToUse(
		txn, gCPAuthenticatedUser)
	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":                   "e86352e8-29f4-4a32-9e96-22f624254731",
			"error":                err,
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got problem extracting users Domains")

		return nil, err

	}

	return domainAndAuthorizations, err
}

// DomainAndAuthorizationsStruct
// Used for holding a Users domain and the Authorizations for that Domain
type DomainAndAuthorizationsStruct struct {
	GCPAuthenticatedUser                                       string
	DomainUuid                                                 string
	DomainName                                                 string
	CanListAndViewTestCaseOwnedByThisDomain                    int64
	CanBuildAndSaveTestCaseOwnedByThisDomain                   int64
	CanListAndViewTestCaseHavingTIandTICFromThisDomain         int64
	CanListAndViewTestCaseHavingTIandTICFromThisDomainExtended int64
	CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain        int64
}

// Concatenate Users specific Domains and Domains open for every one to use
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) concatenateUsersDomainsAndDomainOpenToEveryOneToUse(
	dbTransaction pgx.Tx,
	gCPAuthenticatedUser string) (
	domainAndAuthorizations []DomainAndAuthorizationsStruct,
	err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "80cdaa07-6c59-4107-a442-07be0cd826f2",
	}).Debug("Entering: concatenateUsersDomainsAndDomainOpenToEveryOneToUse()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":                      "5985d322-fb4c-4106-940b-2f53433f5bc0",
			"domainAndAuthorizations": domainAndAuthorizations,
		}).Debug("Exiting: concatenateUsersDomainsAndDomainOpenToEveryOneToUse()")
	}()

	// Load all domains open for every one to use in some way
	var domainsOpenForEveryOneToUse []DomainAndAuthorizationsStruct
	domainsOpenForEveryOneToUse, err = fenixGuiTestCaseBuilderServerObject.loadDomainsOpenForEveryOneToUse(dbTransaction, gCPAuthenticatedUser)
	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":                   "c3d160ea-9122-46fc-a483-0afa54ba45d2",
			"error":                err,
			"GCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Couldn't load all Domains open for every one to use from CloudDB")

		return nil, err
	}

	// Load all domains for a specific user
	domainAndAuthorizations, err = fenixGuiTestCaseBuilderServerObject.loadUsersDomains(dbTransaction, gCPAuthenticatedUser)
	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":                   "b743be52-700b-4d60-b12b-7b873cabac82",
			"error":                err,
			"GCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Couldn't load all Users Domains from CloudDB")

		return nil, err
	}

	// Concatenate Domains and Authorizations
	var domainMap map[string]DomainAndAuthorizationsStruct
	domainMap = make(map[string]DomainAndAuthorizationsStruct)

	// Loop 'domainAndAuthorizations' and add to Map
	for _, termpDomainAndAuthorization := range domainAndAuthorizations {
		domainMap[termpDomainAndAuthorization.DomainUuid] = termpDomainAndAuthorization
	}

	// Loop 'domainsOpenForEveryOneToUse' and add to Map if they don't already exist, if so then replace certain values
	var existsInMap bool
	var termpDomainAndAuthorization DomainAndAuthorizationsStruct
	for _, tempdomainOpenForEveryOneToUse := range domainsOpenForEveryOneToUse {

		termpDomainAndAuthorization, existsInMap = domainMap[tempdomainOpenForEveryOneToUse.DomainUuid]
		if existsInMap == false {
			// Add to Map
			domainMap[tempdomainOpenForEveryOneToUse.DomainUuid] = tempdomainOpenForEveryOneToUse

		} else {
			// Replace values
			termpDomainAndAuthorization.CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain =
				tempdomainOpenForEveryOneToUse.CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain
			termpDomainAndAuthorization.CanListAndViewTestCaseHavingTIandTICFromThisDomain =
				tempdomainOpenForEveryOneToUse.CanListAndViewTestCaseHavingTIandTICFromThisDomain
		}
	}

	// Clear and rebuild 'domainAndAuthorizations'
	domainAndAuthorizations = nil
	for _, tempDomainAndAuthorizations := range domainMap {
		domainAndAuthorizations = append(domainAndAuthorizations, tempDomainAndAuthorizations)
	}

	return domainAndAuthorizations, err

}

// Load all domains for a specific user
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadUsersDomains(
	dbTransaction pgx.Tx,
	gCPAuthenticatedUser string) (
	domainAndAuthorizations []DomainAndAuthorizationsStruct,
	err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "61715721-efc6-4c4e-8856-8157fb6911d5",
	}).Debug("Entering: loadUsersDomains()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":                      "2bb6f88e-892e-47a6-904a-b793ba47df71",
			"domainAndAuthorizations": domainAndAuthorizations,
		}).Debug("Exiting: loadUsersDomains()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT domainuuid, domainname, canlistandviewtestcaseownedbythisdomain, " +
		"canbuildandsavetestcaseownedbythisdomain, canlistandviewtestcasehavingtiandticfromthisdomain, " +
		"canlistandviewtestcasehavingtiandticfromthisdomainextended, canbuildandsavetestcasehavingtiandticfromthisdomain "
	sqlToExecute = sqlToExecute + "FROM \"FenixDomainAdministration\".\"allowedusers\" "
	sqlToExecute = sqlToExecute + "WHERE \"gcpauthenticateduser\" = '" + gCPAuthenticatedUser + "'"
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "2ac8169b-bb80-453c-b8db-0d5696459479",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadUsersDomains'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "68ca4ac8-b33d-4f83-8290-fc270c21a0ea",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	// Extract data from DB result set
	for rows.Next() {

		var tempDomainAndAuthorizations DomainAndAuthorizationsStruct

		err = rows.Scan(
			&tempDomainAndAuthorizations.DomainUuid,
			&tempDomainAndAuthorizations.DomainName,
			&tempDomainAndAuthorizations.CanListAndViewTestCaseOwnedByThisDomain,
			&tempDomainAndAuthorizations.CanBuildAndSaveTestCaseOwnedByThisDomain,
			&tempDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain,
			&tempDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomainExtended,
			&tempDomainAndAuthorizations.CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "5df1b095-dceb-48d7-81a4-30e285ad5b65",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Add user to the row-data
		tempDomainAndAuthorizations.GCPAuthenticatedUser = gCPAuthenticatedUser

		// Append DomainUuid to list of Domains
		domainAndAuthorizations = append(domainAndAuthorizations, tempDomainAndAuthorizations)

	}

	return domainAndAuthorizations, err
}

// Load all domains open for every one to use in some way
func (fenixGuiTestCaseBuilderServerObject *fenixGuiExecutionServerObjectStruct) loadDomainsOpenForEveryOneToUse(
	dbTransaction pgx.Tx,
	gCPAuthenticatedUser string) (
	domainsOpenForEveryOneToUse []DomainAndAuthorizationsStruct,
	err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "b157f2c9-b470-4b22-9acd-22bbeb8e70db",
	}).Debug("Entering: loadDomainsOpenForEveryOneToUse()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":                          "c4d893a3-b40d-4fdc-95b2-fbb13add118b",
			"domainsOpenForEveryOneToUse": domainsOpenForEveryOneToUse,
		}).Debug("Exiting: loadDomainsOpenForEveryOneToUse()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT dom.domain_uuid, dom.domain_name, " +
		"dom.\"AllUsersCanListAndViewTestCaseHavingTIandTICFromThisDomain\", " +
		"dom.\"AllUsersCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain\", dbpn.bitNumberValue "
	sqlToExecute = sqlToExecute + "FROM \"FenixDomainAdministration\".\"domains\" dom, " +
		"\"FenixDomainAdministration\".\"domainbitpositionenum\" dbpn "
	sqlToExecute = sqlToExecute + "WHERE (\"AllUsersCanListAndViewTestCaseHavingTIandTICFromThisDomain\" = true OR "
	sqlToExecute = sqlToExecute + "\"AllUsersCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain\" = true) AND "
	sqlToExecute = sqlToExecute + "dom.\"bitnumbername\" = dbpn.\"bitnumbername\" "
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "f2e26d20-a462-4690-9207-c305ba5e615f",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadDomainsOpenForEveryOneToUse'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "556956bf-90ec-4bfc-a09d-3a22db8c8029",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var tempCanListAndViewTestCaseHavingTIandTICFromThisDomain bool
	var tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain bool
	var bitNumberValue int64

	// Extract data from DB result set
	for rows.Next() {

		var tempDomainAndAuthorizations DomainAndAuthorizationsStruct

		err = rows.Scan(
			&tempDomainAndAuthorizations.DomainUuid,
			&tempDomainAndAuthorizations.DomainName,
			&tempCanListAndViewTestCaseHavingTIandTICFromThisDomain,
			&tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain,
			&bitNumberValue,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "ef811141-d16f-4e24-bcc2-4be4c6f06a7d",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert bool to int64 for 'tempCanListAndViewTestCaseHavingTIandTICFromThisDomain'
		if tempCanListAndViewTestCaseHavingTIandTICFromThisDomain == true {
			tempDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain = bitNumberValue
		}

		// Convert bool to int64 for 'tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain'
		if tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain == true {
			tempDomainAndAuthorizations.CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain = bitNumberValue
		}

		// Add user to the row-data
		tempDomainAndAuthorizations.GCPAuthenticatedUser = gCPAuthenticatedUser

		// Append DomainUuid to list of Domains
		domainsOpenForEveryOneToUse = append(domainsOpenForEveryOneToUse, tempDomainAndAuthorizations)

	}

	return domainsOpenForEveryOneToUse, err
}
