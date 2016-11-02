//
//  2016 October 17
//  John Gilliland [john.gilliland@rndgroup.com]
//

package svclog

import (
	"os"

	"v.io/x/lib/vlog"
)

// LogPath contains constant path string for logging from go rpc server
const LogPath string = "/Logs/Grail/Hamilton/rpc_log"

// Start method configures and starts the logger
func Start() {
	os.RemoveAll(LogPath)
	os.MkdirAll(LogPath, 0777)

	vlog.Configure(vlog.LogDir(LogPath), vlog.AlsoLogToStderr(false), vlog.AutoFlush(true))
	vlog.Info("RPC Log Service started in ", LogPath)
}
