package rabbitmq

type UserReq struct {
	Username string `qs:"username"`
	Password string `qs:"password"`
}

type VhostReq struct {
	Username string `qs:"username"`
	Vhost    string `qs:"vhost"`
	Ip       string `qs:"ip"`
}
