package cmd

import (
	"fmt"

	"github.com/nickvanw/meh"
)

func init() {
	AddPlugin("Meh", "(?i)^\\.meh$", MessageHandler(Meh), false, false)
}

func Meh(msg *Message) {
	key, avail := GetConfig("Meh")
	if !avail {
		fmt.Println("Meh API key not found")
		return
	}

	client := meh.NewClient(meh.WithKey(key))

	current, err := client.Current()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error: %s", err))
		msg.Return("Error contacting Meh API!")
		return
	}

	title := current.Deal.Title
	price := current.Deal.Items[0].Price

	for i := 0; i < len(current.Deal.Items); i++ {
		if current.Deal.Items[i].Price != price {
			// At least one of the prices vary
			data := fmt.Sprintf("Meh.com deal of the day: %s, varying prices", title)
			msg.Return(data)
			return
		}
	}

	data := fmt.Sprintf("Meh.com deal of the day: %s for $%d", title, price)
	msg.Return(data)
}
