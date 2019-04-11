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
	Status     int    `json:"status"`
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
