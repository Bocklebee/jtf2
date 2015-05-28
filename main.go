package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Member struct {
	Id   string `xml:"id,attr"`
	Nick string `xml:"nick,attr"`
}

type Squad struct {
	Nick    string   `xml:"nick,attr"`
	Name    string   `xml:"name"`
	Email   string   `xml:"email"`
	Web     string   `xml:"web"`
	Picture string   `xml:"picture"`
	Title   string   `xml:"title"`
	Members []Member `xml:"member"`
}

const URL = "https://dl.dropboxusercontent.com/u/88240903/squad.xml"
const FILE = "fn_isPlayerAuthorizedForZeus.sqf"

func main() {
	squad := getSquad()
	err := squad.WriteToFile()
	if err != nil {
		panic(err)
	}
}

func (squad *Squad) WriteToFile() error {
	file, err := os.OpenFile(FILE, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("// Data from: " + URL + "\n")
	file.WriteString("// fetched: " + time.Now().Format(time.RFC1123) + "\n")
	file.WriteString("// https://fortkickass.co\n\n")
	file.WriteString("_authorizedUsers = [" + "\n")
	for i, member := range squad.Members {
		file.WriteString(fmt.Sprintf("    \"%s\"", member.Id))
		if i < len(squad.Members)-1 {
			file.WriteString(fmt.Sprintf(", // %s", member.Nick) + "\n")
		} else {
			file.WriteString(fmt.Sprintf("  // %s", member.Nick) + "\n")
		}
	}
	file.WriteString("];" + "\n")
	file.WriteString("\n")
	file.WriteString("_currentPlayerUid = getPlayerUID (_this select 0);" + "\n")
	file.WriteString("_currentPlayerUid in _authorizedUsers;" + "\n")
	return nil
}

func getSquad() *Squad {
	squad := Squad{}
	response, err := http.Get(URL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	decoder := xml.NewDecoder(response.Body)
	err = decoder.Decode(&squad)
	if err != nil {
		panic(err)
	}
	return &squad
}
