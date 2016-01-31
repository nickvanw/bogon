package plugins

import (
	"fmt"
	"regexp"

	"github.com/nickvanw/bogon/commands"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/meh"
)

var mehCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.(meh)$")
	return mehTitle, out, mehLookup, defaultOptions
}

func mehLookup(_ commands.Message, ret commands.MessageFunc) string {
	key, avail := config.Get("MEH_API")
	if !avail {
		return ""
	}

	client := meh.NewClient(meh.WithKey(key))

	current, err := client.Current()
	if err != nil {
		return "Error contacting Meh API!"
	}

	title := current.Deal.Title
	price := current.Deal.Items[0].Price

	for i := 0; i < len(current.Deal.Items); i++ {
		if current.Deal.Items[i].Price != price {
			// At least one of the prices vary
			return fmt.Sprintf("Meh.com deal of the day: %s, varying prices", title)
		}
	}

	return fmt.Sprintf("Meh.com deal of the day: %s for $%d", title, price)
}
