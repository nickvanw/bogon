package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/clbanning/x2j"
)

func init() {
	AddPlugin("NFLScores", "(?i)^\\.nfl$", MessageHandler(NFLScores), false, false)
}

/*
* To do when the NFL season starts again: Switch to other scoreboard, work on some postseason detection as well.
 */
const scores = "http://www.nfl.com/liveupdate/scorestrip/ss.xml"

var quarters = map[string]string{"1": "1st Q", "2": "2nd Q", "3": "3rd Q", "4": "4th Q", "5": "5th Q (OT)"}

func NFLScores(msg *Message) {
	data, _ := getSite(scores)
	jsondata, err := x2j.ByteDocToJson(data)
	if err != nil {
		msg.Return("Error parsing NFL data!")
		return
	}
	var resp NFL
	json.Unmarshal([]byte(jsondata), &resp)
	if len(resp.Ss.Gms.G) == 0 {
		msg.Return("The NFL did not return any games!")
	}
	var games []string
	for _, v := range resp.Ss.Gms.G {
		switch v.Q {
		case "1", "2", "3", "4", "5":
			gameString := fmt.Sprintf("%s [%v] vs %s [%v]; %s - %s", bold(upperFirst(v.Hnn)), v.Hs, bold(upperFirst(v.Vnn)), v.Vs, v.K, quarters[v.Q])
			games = append(games, gameString)
		case "H":
			gameString := fmt.Sprintf("%s [%v] vs %s [%v]; Half", bold(upperFirst(v.Hnn)), v.Hs, bold(upperFirst(v.Vnn)), v.Vs)
			games = append(games, gameString)
		}
	}
	if len(games) > 0 {
		msg.Return(fmt.Sprintf("NFL Games: %s", strings.Join(games, ", ")))
	} else {
		msg.Return("There are currently no NFL games.")
	}
}

type NFL struct {
	Ss struct {
		Gms struct {
			_Bph string `json:"-bph"`
			_Gd  string `json:"-gd"`
			_T   string `json:"-t"`
			_W   string `json:"-w"`
			_Y   string `json:"-y"`
			G    []struct {
				D    string `json:"-d"`
				Eid  string `json:"-eid"`
				Ga   string `json:"-ga"`
				Gsis string `json:"-gsis"`
				Gt   string `json:"-gt"`
				H    string `json:"-h"`
				Hnn  string `json:"-hnn"`
				Hs   string `json:"-hs"`
				Q    string `json:"-q"`
				Rz   string `json:"-rz"`
				T    string `json:"-t"`
				V    string `json:"-v"`
				Vnn  string `json:"-vnn"`
				Vs   string `json:"-vs"`
				K    string `json:"-k"`
			} `json:"g"`
		} `json:"gms"`
	} `json:"ss"`
}
