package sapi

import (
	"strings"

	"github.com/happyhippyhippo/slate"
)

// ----------------------------------------------------------------------------
// defs
// ----------------------------------------------------------------------------

const (
	// RestLogMwContainerID defines the default id used to register
	// the application log middleware and related services.
	RestLogMwContainerID = RestContainerID + ".log.mw"

	// RestLogMwEnvID defines the log middleware module base
	// environment variable name.
	RestLogMwEnvID = RestEnvID + "_LOG_MW"
)

var (
	// RestLogMwRequestChannel defines the channel id to be used when
	// the log middleware sends the request logging signal to the logger
	// instance.
	RestLogMwRequestChannel = slate.EnvString(RestLogMwEnvID+"_REQUEST_CHANNEL", "rest")

	// RestLogMwRequestLevel defines the logging level to be used when
	// the log middleware sends the request logging signal to the logger
	// instance.
	RestLogMwRequestLevel = envToLogLevel(RestLogMwEnvID+"_REQUEST_LEVEL", slate.DEBUG)

	// RestLogMwRequestMessage defines the request event logging message to
	// be used when the log middleware sends the logging signal to the logger
	// instance.
	RestLogMwRequestMessage = slate.EnvString(RestLogMwEnvID+"_REQUEST_MESSAGE", "Request")

	// RestLogMwResponseChannel defines the channel id to be used when the
	// log middleware sends the response logging signal to the logger instance.
	RestLogMwResponseChannel = slate.EnvString(RestLogMwEnvID+"_RESPONSE_CHANNEL", "rest")

	// RestLogMwResponseLevel defines the logging level to be used when the
	// log middleware sends the response logging signal to the logger instance.
	RestLogMwResponseLevel = envToLogLevel(RestLogMwEnvID+"_RESPONSE_LEVEL", slate.INFO)

	// RestLogMwResponseMessage defines the response event logging message
	// to be used when the log middleware sends the logging signal to the
	// logger instance.
	RestLogMwResponseMessage = slate.EnvString(RestLogMwEnvID+"_RESPONSE_MESSAGE", "Response")
)

func envToLogLevel(ev string, def slate.LogLevel) slate.LogLevel {
	v, ok := slate.LogLevelMap[strings.ToLower(ev)]
	if !ok {
		return def
	}
	return v
}
