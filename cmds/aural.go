package main

import (
	"flag"
	"log"

	"github.com/aural/aural"

	zmq "github.com/pebbe/zmq4"
)

func init() {
	flag.Parse()
}

func tracks(locations []string) (tracks []*aural.Track) {
	var track *aural.Track

	for _, location := range locations {
		track = new(aural.Track)
		track.Location = location
		tracks = append(tracks, track)
	}

	return tracks
}
func createServer(configuration *aural.Configuration, playstate *aural.Playstate) chan string {
	channel := make(chan string)
	server, err := zmq.NewSocket(zmq.REP)

	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		defer server.Close()

		log.Println("Listening at address:", configuration.Address)
		server.Bind(configuration.Address)

		for {
			request, err := server.RecvMessage(0)

			if err != nil {
				log.Fatalln(err)
			}

			var arguments []string

			if len(request) == 0 {
				continue
			} else if len(request) > 1 {
				arguments = request[1:]
			} else {
				arguments = []string{}
			}

			aural.HandleRequest(playstate, request[0], arguments)
			server.SendMessage("OK")
		}
	}()

	return channel
}

func main() {
	configuration := aural.GetConfiguration()

	log.Println("Starting aural daemon")
	defer aural.Terminate()

	playstate, err := aural.NewPlaystate()

	if err != nil {
		log.Fatalln("Could not acquire context for audio hardware:", err)
	}

	go createServer(&configuration, playstate)
	audio := playstate.MainLoop()

	for {
		<-audio
	}
}
