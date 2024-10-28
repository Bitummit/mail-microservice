package my_kafka

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Bitummit/mail-microservice/pkg/proto"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

type Kafka struct {
	Conn *kafka.Conn
	Brokers []string
	Topic string
	ConsumerGroup *kafka.ConsumerGroup
	Reader *kafka.Reader
	Server HTTPServer
}

type HTTPServer interface {
	Send(ctx context.Context, req *proto.EmailRequest) (*proto.EmailResponse, error)
}

func New(ctx context.Context, leaderAddress, topic, group_id string, partition int, brokers []string) (*Kafka, error){
	// Create consumergroup
	// Создавать топик с партициями тут вручную или самой кафкой?
	// На одном брокере несколько топиков норм?
	conn, err := kafka.DialLeader(ctx, "tcp", leaderAddress, topic, partition)
	if err != nil {
		return nil, fmt.Errorf("failed to dial leader: %w", err)
	}
	conn.SetReadDeadline(time.Now().Add(10*time.Second))

	groupConfig := kafka.ConsumerGroupConfig{
		ID: group_id,
		Brokers: brokers,
		Topics: []string{topic},
	}
	consumerGroup, err := kafka.NewConsumerGroup(groupConfig)
	if err != nil {
		return nil, fmt.Errorf("faield to create consumer group: %w", err)
	}

	return &Kafka{
		Conn: conn,
		Brokers: brokers,
		Topic: topic,
		ConsumerGroup: consumerGroup,
	}, nil

}

func (k *Kafka) RunConsumerWithGroup(ctx context.Context, group_id string) {
	// с группой не работает
	// после перезапуска читает все евенты
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:	k.Brokers,
		// GroupID:	group_id,
		// Partition: 0,
		Topic:		k.Topic,
	})
	k.Reader = r

	go func() {
		for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			r.Close()
			break
		}
		
		req := &proto.EmailRequest{
			To: []string{string(m.Value)},
			Subject: "You are registered",
			Body: "Thanks for registration",
		}
		_, err = k.Server.Send(ctx, req)
		if err != nil {
			slog.Error("send email: %w", err)
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
	}()
}


// conn, err := kafka.Dial("tcp", "localhost:9092")
// 	if err != nil {
// 		return nil, fmt.Errorf("dialing kafka: %w", err)
// 	}
// 	defer conn.Close()

// 	controller, err := conn.Controller()
// 	if err != nil {
// 		return nil, fmt.Errorf("creating controller: %w", err)
// 	}
// 	var controllerConn *kafka.Conn
// 	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
// 	if err != nil {
// 		return nil, fmt.Errorf("dialing controller: %w", err)
// 	}
// 	defer controllerConn.Close()

// 	topicConfigs := []kafka.TopicConfig{
// 		{
// 			Topic:             topic,
// 			NumPartitions:     1,
// 			ReplicationFactor: 1,
// 		},
// 	}