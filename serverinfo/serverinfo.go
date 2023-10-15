package serverinfo

import (
	"fmt"
	"strings"
)

var (
	ServerInfo string
	info       *Information
)

type Information struct {
	CommitID  string
	Author    string
	Branch    string
	BuildTime string
	GoVersion string
	BuildMode string
}

func InitInformation() {
	if ServerInfo == "" {
		info = &Information{}
	}
	si := strings.Split(ServerInfo, ";")
	if len(si) != 6 {
		info = &Information{}
	} else {
		info = &Information{
			CommitID:  si[0],
			Author:    si[1],
			Branch:    si[2],
			BuildTime: si[3],
			GoVersion: si[4],
			BuildMode: si[5],
		}
	}
}

func Get() *Information {
	if info == nil {
		panic("ServerInfo is not initialized")
	}
	return info
}

func (i *Information) String() string {
	return fmt.Sprintf(
		"\nCommitID: %s\n"+
			"Author: %s\n"+
			"Branch: %s\n"+
			"BuildTime: %s\n"+
			"GoVersion: %s\n"+
			"BuildMode: %s\n",
		i.CommitID,
		i.Author,
		i.Branch,
		i.BuildTime,
		i.GoVersion,
		i.BuildMode)
}
