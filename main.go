package main

import (
	"flag"
	"fmt"
	"github.com/bigdatadev/goryman"
	"github.com/hpcloud/tail"
	"strings"
)

func parseEvent(line string) (*goryman.Event, error) {
	if strings.Contains(line, "SERVICE CHECK") {
		splitLine := strings.SplitN(line, ": ", 2)
		splitAlert := strings.Split(splitLine[1], ";")
		var state = ""
		switch splitAlert[2] {
		case "0":
			state = "ok"
		case "1":
			state = "warning"
		case "2":
			state = "critical"
		case "3":
			state = "unknown"
		default:
			state = splitAlert[2]
		}
		var event = &goryman.Event{
			Host:        splitAlert[0],
			Service:     splitAlert[1],
			State:       state,
			Description: splitAlert[3],
			Tags:        []string{"nagios"},
			Ttl:         600,
		}
		return event, nil
	}
	return nil, fmt.Errorf("Event not found")

}

func sendToRiemann(events chan *goryman.Event, riemannHost *string) {
	client := goryman.NewGorymanClient(*riemannHost)
	err := client.Connect()
	defer client.Close()
	if err != nil {
		panic(err)
	}
	for {
		event := <-events
		client.SendEvent(event)
	}
}

func tailFile(events chan *goryman.Event, fileName string, poll bool) {
	t, _ := tail.TailFile(fileName, tail.Config{Follow: true, ReOpen: true, Poll: poll})
	for line := range t.Lines {
		event, err := parseEvent(line.Text)
		// Swollow parse errrors.
		if err == nil {
			events <- event
		}
	}

}
func main() {
	var events = make(chan *goryman.Event)

	fileName := flag.String("nagioslog", "/var/log/icinga/icinga.log", "Nagios log file.")
	poll := flag.Bool("poll", false, "Poll instead of use inotify")
	riemannHost := flag.String("server", "SEVER:5555", "Rieman host and port")
	flag.Parse()
	go sendToRiemann(events, riemannHost)
	tailFile(events, *fileName, *poll)
}
