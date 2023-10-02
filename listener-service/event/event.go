package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declaerExchange(ch *amqp.Channel) error {

	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      //type
		true,         //durable?
		false,        //autodeleted?
		false,        //internal?
		false,        //no-wait?
		nil,          //arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //name?
		false, //durable?
		false, //delete when unused?
		true,  //explusive?
		false, //no-wait?
		nil,   //arguments?
	)
}
