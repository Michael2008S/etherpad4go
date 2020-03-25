package poker

import "github.com/Michael2008S/etherpad4go/api"

var sessionInfo map[string]*authInfo

type authInfo struct {
	sessionID string
	padID     string
	token     string
	password  string
	author    string
	rev       int
}

func init() {
	sessionInfo = make(map[string]*authInfo, 1)
}

func createSessionInfo(client *Client, req api.ClientReadyReq) {
	// Remember this information since we won't
	// have the cookie in further socket.io messages.
	// This information will be used to check if
	// the sessionId of this connection is still valid
	// since it could have been deleted by the API.
	//auth :=
	sessionInfo[client.ID] = &authInfo{
		sessionID: req.SessionID,
		padID:     req.PadID,
		token:     req.Token,
		password:  req.Password,
	}
}
