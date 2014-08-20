package cmd

func init() {
	AddPlugin("AdminLogin", "(?i)^\\.login$", MessageHandler(Login), false, false)
	AddPlugin("AdminCheck", "(?i)^\\.admin$", MessageHandler(AdminCheck), false, false)
	AddPlugin("AdminOnly", "(?i)^\\.adminonly$", MessageHandler(AdminOnly), false, true)
	AddPlugin("ReloadData", "(?i)^\\.reloaddata$", MessageHandler(ReloadData), false, true)
}

func Login(msg *Message) {
	if msg.State.Password != "" {
		if msg.Params[1] == msg.State.Password && msg.User.String() != "" {
			msg.State.Admin = msg.User.String()
			msg.Return(msg.User.String() + " is now my admin")
		}
	}
}

func AdminCheck(msg *Message) {
	if msg.IsAdmin {
		msg.Return("My admin is " + msg.State.Admin + " and you are it!")
	} else {
		msg.Return("Sorry, you're not an admin")
	}
}

func AdminOnly(msg *Message) {
	msg.Return("You got it, bud!")
}

func ReloadData(msg *Message) {
	err := LoadConfig(myConfig)
	if err != nil {
		msg.Return("Unable to reload data: " + err.Error())
	} else {
		msg.Return("Reloaded the data successfully")
	}
}
