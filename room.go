package eyeson

// RoomResponse holds available attributes to a room.
type RoomResponse struct {
	AccessKey string    `json:"access_key"`
	Links     RoomLinks `json:"links"`
	Room      Room      `json:"room"`
	User      User      `json:"user"`
	Ready     bool      `json:"ready"`
	Signaling Signaling `json:"signaling"`
}

// Signaling base container for signaling options.
// So far only type "sepp" is allowed.
type Signaling struct {
	Type    string        `json:"type"`
	SigSepp SeppSignaling `json:"options"`
}

// SeppSignaling holds information required by the
// gosepp signaling interface.
type SeppSignaling struct {
	ClientID    string       `json:"client_id"`
	AuthToken   string       `json:"auth_token"`
	ConfID      string       `json:"conf_id"`
	Endpoint    string       `json:"endpoint"`
	StunServers []string     `json:"stun_servers"`
	TurnServer  []TurnServer `json:"turn_servers"`
}

// GetSigEndpoint returns the signaling endpoint
func (rr *RoomResponse) GetSigEndpoint() string {
	return rr.Signaling.SigSepp.Endpoint
}

// GetAuthToken returns the JWT-Authtoken for
// authenticating to the sig-service.
func (rr *RoomResponse) GetAuthToken() string {
	return rr.Signaling.SigSepp.AuthToken
}

// GetClientID returns the client-id of this
// signaling entity.
func (rr *RoomResponse) GetClientID() string {
	return rr.Signaling.SigSepp.ClientID
}

// GetConfID returns the conf-id to connect to.
func (rr *RoomResponse) GetConfID() string {
	return rr.Signaling.SigSepp.ConfID
}

// GetStunServers returns stun info
func (rr *RoomResponse) GetStunServers() []string {
	return rr.Signaling.SigSepp.StunServers
}

// GetTurnServerURLs returns turn info
func (rr *RoomResponse) GetTurnServerURLs() []string {
	sepp := rr.Signaling.SigSepp
	if len(sepp.TurnServer) > 0 {
		return sepp.TurnServer[0].URLs
	}
	return []string{}
}

// GetTurnServerPassword returns turn credentials
func (rr *RoomResponse) GetTurnServerPassword() string {
	sepp := rr.Signaling.SigSepp
	if len(sepp.TurnServer) > 0 {
		return sepp.TurnServer[0].Password
	}
	return ""
}

// GetTurnServerUsername return turn credentials
func (rr *RoomResponse) GetTurnServerUsername() string {
	sepp := rr.Signaling.SigSepp
	if len(sepp.TurnServer) > 0 {
		return sepp.TurnServer[0].Username
	}
	return ""
}

func (rr *RoomResponse) GetDisplayname() string {
	return rr.User.Name
}

// TurnServer provides connection info for ICE-Servers
type TurnServer struct {
	URLs     []string `json:"urls"`
	Username string   `json:"username"`
	Password string   `json:"password"`
}

// Room has attributes for SIP details and GuestToken
type Room struct {
	GuestToken string `json:"guest_token"`
	SIP        SIP    `json:"sip"`
	Shutdown   bool   `json:"shutdown"`
}

// RoomLinks provide all public web URLs for a room.
type RoomLinks struct {
	Gui       string `json:"gui"`
	GuestJoin string `json:"guest_join"`
	Websocket string `json:"websocket"`
}

// SIP contains access details for the protocol used to establish a connection.
type SIP struct {
	URI               string `json:"uri"`
	Domain            string `json:"domain"`
	AuthorizationUser string `json:"authorizationUser"`
	Password          string `json:"password"`
	WSServers         string `json:"wsServers"`
	DisplayName       string `json:"displayName"`
}

// User has information on the current participant that has joined the meeting.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	SIP  SIP    `json:"sip"`
}
