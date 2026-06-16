package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Routes struct {
	Queue string
	Route string
}

func New(
	username string,
	password string,
	host string,
	port string,
) (*amqp.Channel, error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		username,
		password,
		host,
		port,
	)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func Exchange(
	channel *amqp.Channel,
	exchangeName string,
	exchangeType string,
) error {
	if err := channel.ExchangeDeclare(
		exchangeName, // exchange name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,
	); err != nil {
		return err
	}
	return nil
}

func Queue(
	channel *amqp.Channel,
	exchangeName string,
	routes []Routes,
) error {
	for _, v := range routes {
		q, err := channel.QueueDeclare(
			v.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		if err := channel.QueueBind(
			q.Name,
			v.Route,
			exchangeName,
			false,
			nil,
		); err != nil {
			return err
		}
	}
	return nil
}
