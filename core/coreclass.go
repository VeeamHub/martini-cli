package core

type HeartbeatMessage struct {
	Status string `json:"status"`
}
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Token struct {
	Token      string `json:"token"`
	Renew      string `json:"renew"`
	Lifetime   int64  `json:"lifetime"`
	Status     string `json:"status"`
	ServerTime int64  `json:"now"`
}
type ClientConfig struct {
	Server     string
	Token      string
	Renew      string
	Lifetime   int64
	Username   string
	ServerSkew int64
}

//very generic return
//can be used for error handling
type ReturnStatus struct {
	Status string `json:"status"`
	Id     string `json:"id,omitempty"`
	SubId  string `json:"subid,omitempty"`
}

//very generic send option
type SendID struct {
	Id     string `json:"id,omitempty"`
	Action string `json:"action,omitempty"`
	Data   string `json:"data,omitempty"`
}
