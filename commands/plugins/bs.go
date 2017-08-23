package plugins

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/nickvanw/bogon/commands"
)

var bsCommand = func() (string, *regexp.Regexp, commands.CommandFunc, commands.Options) {
	out := regexp.MustCompile("(?i)^\\.bs$")
	return bsTitle, out, bsGenerate, defaultOptions
}

var adverbs = []string{"abstractly", "analytically", "disruptively", "efficiently", "financially", "holistically",
	"sufficiently", "sustainably", "quantifiably", "quickly"}
var verbs = []string{"analyzing", "encapsulating", "disrupting", "envisioning", "growing", "leveraging", "redefining",
	"reenvisioning", "sustaining", "testing", "unboxing"}
var nouns = []string{"alignment", "boundaries", "clutter", "density", "evolution", "horizons", "logistics", "missions",
	"paradigms", "potential", "sustainability", "transformation", "trust"}

func bsGenerate(msg commands.Message, ret commands.MessageFunc) string {
	// Choose a random adverb
	rand.Seed(time.Now().UTC().UnixNano())
	adverbIndex := rand.Intn(len(adverbs))
	adverb := adverbs[adverbIndex]

	// Choose a random verb
	verbIndex := rand.Intn(len(verbs))
	verb := verbs[verbIndex]

	// Choose a random
	nounIndex := rand.Intn(len(nouns))
	noun := nouns[nounIndex]

	// Puts them together
	return fmt.Sprintf("%s %s %s", adverb, verb, noun)
}
