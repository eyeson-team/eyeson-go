package eyeson

// RoomResponse holds available attributes to a room.
type RoomResponse struct {
	AccessKey string    `json:"access_key"`
	Links     RoomLinks `json:"links"`
	Room      Room      `json:"room"`
	User      User      `json:"user"`
	Ready     bool      `json:"ready"`
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
