package messagesToExecutionServer

import (
	"FenixGuiExecutionServer/common_config"
	"FenixGuiExecutionServer/grpcurl"
	"crypto/tls"
	"encoding/base64"
	fenixExecutionServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
	"log"
	"strings"
	"time"
)

// SetConnectionToExecutionServer - Set upp connection and Dial to  FenixExecutionServer
func (messagesToExecutionServerObject *MessagesToExecutionServerObjectStruct) SetConnectionToExecutionServer(ctx context.Context) (_ context.Context, err error) {

	// slice with sleep time, in milliseconds, between each attempt to Dial to ExecutionServer
	var sleepTimeBetweenDialAttempts []int
	sleepTimeBetweenDialAttempts = []int{100, 100, 200, 200, 300, 300, 500, 500, 600, 1000} // Total: 3.6 seconds

	//var opts []grpc.DialOption

	// Do multiple attempts to do connection to ExecutionServer
	var numberOfDialAttempts int
	var dialAttemptCounter int
	numberOfDialAttempts = len(sleepTimeBetweenDialAttempts)
	dialAttemptCounter = 0

	for {

		//When target is running on GCP then use credential otherwise not
		/*
			if common_config.ExecutionLocationForFenixExecutionServer == common_config.GCP {
				creds := credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true,
				})

				opts = []grpc.DialOption{
					grpc.WithTransportCredentials(creds),
				}
			}

		*/

		// Set up connection to Fenix ExecutionServer
		// When target is running on GCP, use credentials
		var newGrpcClientConnection *grpc.ClientConn
		if common_config.ExecutionLocationForFenixExecutionServer == common_config.GCP {
			// Run on GCP
			ctx, newGrpcClientConnection = dialFromGrpcurl(ctx)
			RemoteFenixExecutionServerConnection = newGrpcClientConnection
			//RemoteFenixExecutionServerConnection, err = grpc.Dial(FenixExecutionServerAddressToDial, opts...)
		} else {
			// Run Local
			RemoteFenixExecutionServerConnection, err = grpc.Dial(FenixExecutionServerAddressToDial, grpc.WithInsecure())
		}

		// Add to counter for how many Dial attempts that have been done
		dialAttemptCounter = dialAttemptCounter + 1

		if err != nil {
			messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
				"ID":                                "50b59b1b-57ce-4c27-aa84-617f0cde3100",
				"FenixExecutionServerAddressToDial": FenixExecutionServerAddressToDial,
				"error message":                     err,
				"dialAttemptCounter":                dialAttemptCounter,
			}).Error("Did not connect to FenixExecutionServer via gRPC")

			// Only return the error after last attempt
			if dialAttemptCounter >= numberOfDialAttempts {
				return nil, err
			}

		} else {
			messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
				"ID":                                "0c650bbc-45d0-4029-bd25-4ced9925a059",
				"FenixExecutionServerAddressToDial": FenixExecutionServerAddressToDial,
			}).Info("gRPC connection OK to FenixExecutionServer")

			// Creates a new Clients
			FenixExecutionServerGrpcClient = fenixExecutionServerGrpcApi.NewFenixExecutionServerGrpcServicesClient(RemoteFenixExecutionServerConnection)

			return nil, err
		}

		// Sleep for some time before retrying to connect
		time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenDialAttempts[dialAttemptCounter-1]))

	}

	return ctx, err
}

// Generate Google access token. Used when running in GCP
func (messagesToExecutionServerObject *MessagesToExecutionServerObjectStruct) generateGCPAccessToken(ctx context.Context) (appendedCtx context.Context, returnAckNack bool, returnMessage string) {

	// Only create the token if there is none, or it has expired
	if messagesToExecutionServerObject.gcpAccessToken == nil || messagesToExecutionServerObject.gcpAccessToken.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.
		tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+FenixExecutionServerAddressToUse)
		if err != nil {
			messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
				"ID":  "8ba622d8-b4cd-46c7-9f81-d9ade2568eca",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		if err != nil {
			messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
				"ID":  "0cf31da5-9e6b-41bc-96f1-6b78fb446194",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
				"ID":    "8b1ca089-0797-4ee6-bf9d-f9b06f606ae9",
				"token": token,
			}).Debug("Got Bearer Token")
		}

		messagesToExecutionServerObject.gcpAccessToken = token

	}

	messagesToExecutionServerObject.Logger.WithFields(logrus.Fields{
		"ID": "cd124ca3-87bb-431b-9e7f-e044c52b4960",
		"messagesToExecutionServerObject.gcpAccessToken": messagesToExecutionServerObject.gcpAccessToken,
	}).Debug("Will use Bearer Token")

	// Add token to GrpcServer Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+messagesToExecutionServerObject.gcpAccessToken.AccessToken)

	return appendedCtx, true, ""

}

// GetHighestFenixExecutionServerProtoFileVersion
// Get the highest FenixProtoFileVersionEnumeration for ExecutionServer-gRPC-api
func (messagesToExecutionServerObject *MessagesToExecutionServerObjectStruct) GetHighestFenixExecutionServerProtoFileVersion() int32 {

	// Check if there already is a 'highestFenixExecutionServerProtoFileVersion' saved, if so use that one
	if highestFenixExecutionServerProtoFileVersion != -1 {
		return highestFenixExecutionServerProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixExecutionServerGrpcApi.CurrentFenixExecutionServerProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestFenixExecutionServerProtoFileVersion = maxValue

	return highestFenixExecutionServerProtoFileVersion
}

var (
	isUnixSocket func() bool
)

func dialFromGrpcurl(ctx context.Context) (context.Context, *grpc.ClientConn) {

	target := FenixExecutionServerAddressToDial

	dialTime := 10 * time.Second

	ctx, cancel := context.WithTimeout(ctx, dialTime)
	defer cancel()
	var opts []grpc.DialOption

	var creds credentials.TransportCredentials

	var tlsConf *tls.Config

	creds = credentials.NewTLS(tlsConf)

	grpcurlUA := "FenixCAConnector"
	//if grpcurl.version == grpcurl.no_version {
	//	grpcurlUA = "grpcurl/dev-build (no version set)"
	//}

	opts = append(opts, grpc.WithUserAgent(grpcurlUA))
	//opts = append(opts, grpc.WithNoProxy())

	network := "tcp"
	if isUnixSocket != nil && isUnixSocket() {
		network = "unix"
	}

	cc, err := grpcurl.BlockingDial(ctx, network, target, creds, opts...)
	if err != nil {
		log.Panicln("Failed to Dial, ", target, err.Error())
	}
	return ctx, cc

}

// MetadataFromHeaders converts a list of header strings (each string in
// "Header-Name: Header-Value" form) into metadata. If a string has a header
// name without a value (e.g. does not contain a colon), the value is assumed
// to be blank. Binary headers (those whose names end in "-bin") should be
// base64-encoded. But if they cannot be base64-decoded, they will be assumed to
// be in raw form and used as is.
func MetadataFromHeaders(headers []string) grpcMetadata.MD {
	md := make(grpcMetadata.MD)
	for _, part := range headers {
		if part != "" {
			pieces := strings.SplitN(part, ":", 2)
			if len(pieces) == 1 {
				pieces = append(pieces, "") // if no value was specified, just make it "" (maybe the header value doesn't matter)
			}
			headerName := strings.ToLower(strings.TrimSpace(pieces[0]))
			val := strings.TrimSpace(pieces[1])
			if strings.HasSuffix(headerName, "-bin") {
				if v, err := decode(val); err == nil {
					val = v
				}
			}
			md[headerName] = append(md[headerName], val)
		}
	}
	return md
}

var base64Codecs = []*base64.Encoding{base64.StdEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.RawURLEncoding}

func decode(val string) (string, error) {
	var firstErr error
	var b []byte
	// we are lenient and can accept any of the flavors of base64 encoding
	for _, d := range base64Codecs {
		var err error
		b, err = d.DecodeString(val)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		return string(b), nil
	}
	return "", firstErr
}
