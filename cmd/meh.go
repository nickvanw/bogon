package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	AddPlugin("Meh", "(?i)^\\.meh$", MessageHandler(Meh), false, false)
}

func Meh(msg *Message) {
	meh, err := GetMeh()
	if err != nil {
		msg.Return(fmt.Sprintf("Error: %s", err))
		return
	}
	data := fmt.Sprintf("Meh.com deal of the day: %s for %s", meh.Name, meh.Price)
	msg.Return(data)
}

type MehItemOfTheDay struct {
	// Basic
	Name  string
	Price string
}

func GetMeh() (*MehItemOfTheDay, error) {
	doc, e := goquery.NewDocument("https://meh.com/")

	if e != nil {
		return nil, errors.New("Unable to get to meh.com")
	}

	i_name := strings.TrimSpace(strings.Trim(doc.Find("section.features h2").Text(), "\n"))
	i_price := strings.TrimSpace(strings.Trim(doc.Find("button.buy-button").Text(), "\n"))

	data := &MehItemOfTheDay{Name: i_name, Price: i_price}
	return data, nil
}
