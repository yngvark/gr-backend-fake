package main

import (
	"encoding/json"
	"fmt"
	"github.com/yngvark/gr-zombie/pkg/connectors"
	gamelogicPkg "github.com/yngvark/gr-zombie/pkg/gamelogic"
	"github.com/yngvark/gr-zombie/pkg/worldmap"
)

func runGameLogic(o *GameOpts) error {
	// Close producer and consumer when done
	defer func() {
		err := o.connector.Close()
		if err != nil {
			o.log.Errorf("error closing consumer: %s", err.Error())
		}
	}()

	// Create game
	gameLogic := gamelogicPkg.NewGameLogic(o.context, o.log, o.broadcaster)

	err := o.connector.ListenForConnections(createOnConnect(o))
	if err != nil {
		o.log.Errorf("Error listening for connections: %s", err.Error())
		o.cancelFn()
	}

	go func() {
		for {
			select {
			case msgFromClient := <-o.subscriber:
				o.log.Debugf("Msg from client: %s", msgFromClient)
			case <-o.context.Done():
				o.log.Debug("gamelogic: Stop listnening for messages from client")
			}
		}
	}()

	o.log.Info("Running game")
	gameLogic.Run()
	o.log.Info("Done running game")

	o.log.Debug("runGameLogic: cancelFn")
	o.cancelFn()

	return nil
}

func createOnConnect(o *GameOpts) connectors.OnConnect {
	return func(messagesToClientChannel chan string) error {
		o.log.Debug("Client connected. Sending world map.")

		wmap := worldmap.New(30, 30)

		wmapJSON, err := json.Marshal(wmap)
		if err != nil {
			return fmt.Errorf("could not marshal world map: %w", err)
		}

		messagesToClientChannel <- string(wmapJSON)

		return nil
	}
}
