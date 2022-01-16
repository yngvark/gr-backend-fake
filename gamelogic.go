package main

import (
	gamelogicPkg "github.com/yngvark/gr-zombie/pkg/gamelogic"
)

func runGameLogic(o *GameOpts) error {
	// Close producer and consumer when done
	defer func() {
		err := o.publisher.Close()
		if err != nil {
			o.logger.Errorf("error closing publisher: %s", err.Error())
		}
	}()

	defer func() {
		err := o.consumer.Close()
		if err != nil {
			o.logger.Errorf("error closing consumer: %s", err.Error())
		}
	}()

	// Create game
	gameLogic := gamelogicPkg.NewGameLogic(o.context, o.logger, o.publisher)

	// Wait until some external orchestrator sends a "start" message
	go func() {
		err := o.consumer.ListenForMessages()
		if err != nil {
			o.logger.Errorf("Error listening for messages: %s", err.Error())
		}
	}()

	o.logger.Info("Waiting for start message...")

listenForStartMsg:
	for {
		select {
		case msg := <-o.consumer.SubscriberChannel():
			o.logger.Infof("Waiting for start message... Received: %s", msg)

			if msg == "start" {
				break listenForStartMsg
			}
		case <-o.context.Done():
			o.logger.Info("Aborted waiting for game to start")
			return nil
		}
	}

	o.logger.Info("Running game")
	gameLogic.Run()
	o.logger.Info("Done running game")

	return nil
}
