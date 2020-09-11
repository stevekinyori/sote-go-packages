package sMessage

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"gitlab.com/soteapps/packages/v2020/sError"
	"gitlab.com/soteapps/packages/v2020/sLogger"
)

func SetAllOptions(streamName, streamCredentialFile, streamCredentialToken string, maxReconnect int,
	reconnectWait, timeOut time.Duration) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if opts, soteErr = SetStreamName(streamName); soteErr.ErrCode == nil {
		if opts, soteErr = SetCredentialsFile(streamCredentialFile); soteErr.ErrCode == nil {
			if opts, soteErr = SetCredentialsToken(streamCredentialToken); soteErr.ErrCode == nil {
				if opts, soteErr = SetReconnectOptions(maxReconnect, reconnectWait); soteErr.ErrCode == nil {
					opts, soteErr = SetTimeOut(timeOut)
				}
			}
		}
	}

	return
}

func SetStreamName(streamName string) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamName) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamName"}), nil)
	} else {
		opts = []nats.Option{nats.Name(streamName)}
	}

	return
}

/*
	streamCredentialToken will take precedence over streamCredentialFile
*/
func SetCredentialsFile(streamCredentialFile string) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamCredentialFile) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamCredentialFile"}), nil)
	} else if _, err := os.Stat(streamCredentialFile); err != nil {
		soteErr = sError.GetSError(600010, sError.BuildParams([]string{streamCredentialFile, err.Error()}), nil)
	} else {
		opts = []nats.Option{nats.UserCredentials(streamCredentialFile)}
	}

	return
}

/*
	streamCredentialToken will take precedence over streamCredentialFile
*/
func SetCredentialsToken(streamCredentialToken string) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(streamCredentialToken) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"streamCredentialToken"}), nil)
	} else {
		opts = []nats.Option{nats.Token(streamCredentialToken)}
	}

	return
}

func SetReconnectOptions(maxReconnect int, reconnectWait time.Duration) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if maxReconnect == 0 && reconnectWait == 0 {
		soteErr = sError.GetSError(200512, sError.BuildParams([]string{"maxReconnect", "reconnectWait"}), nil)
	} else {
		if reconnectWait == 0 {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"reconnectWait"}), nil)
		} else {
			opts = []nats.Option{nats.ReconnectWait(reconnectWait)}
		}
		if maxReconnect == 0 {
			soteErr = sError.GetSError(200513, sError.BuildParams([]string{"reconnectWait"}), nil)
		} else {
			opts = []nats.Option{nats.MaxReconnects(maxReconnect)}
		}
	}

	return
}

func SetTimeOut(timeOut time.Duration) (opts []nats.Option, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if timeOut == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"timeOut"}), nil)
	} else {
		opts = []nats.Option{nats.Timeout(timeOut)}
	}

	return
}

// This will connect to the NATS network
func Connect(url string, opts []nats.Option) (nc *nats.Conn, soteErr sError.SoteError) {
	sLogger.DebugMethod()

	if len(url) == 0 {
		soteErr = sError.GetSError(200513, sError.BuildParams([]string{"url"}), nil)
	} else {
		var err error
		// Connect to NATS  --  default "euwest1.aws.ngs.global"
		nc, err = nats.Connect(url, opts...)
		if err != nil {
			log.Fatal(err)
		}
		defer nc.Close()
	}

	return
}
