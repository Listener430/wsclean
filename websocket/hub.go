package websocket

func NewOnlineHub() *OnlineHub {
	return &OnlineHub{
		Register:   make(chan *OnlineClient),
		Unregister: make(chan *OnlineClient),
		Users:      map[string]*OnlineClient{},
	}
}

func (ohb *OnlineHub) Run() {
	for {
		select {
		case client := <-ohb.Register:
			HandleOnlineUserRegistration(ohb, client)
		case client := <-ohb.Unregister:
			HandleUserDisconnect(ohb, client)
			//case <-ohb.Quit:
			return
		}
	}
}

func HandleOnlineUserRegistration(ohb *OnlineHub, cl *OnlineClient) {
	ohb.Users[cl.UserId] = cl
}

func HandleUserDisconnect(ohb *OnlineHub, cl *OnlineClient) {
	delete(ohb.Users, cl.UserId)
}
