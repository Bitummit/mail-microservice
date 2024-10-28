package my_kafka

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

type Kafka struct {
	Conn *kafka.Conn
	LeaderAddress string
	Brokers []string
	Topic string
	ConsumerGroup *kafka.ConsumerGroup
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
		return nil, fmt.Errorf("faield to create consumer group")
	}

	return &Kafka{
		Conn: conn,
		LeaderAddress: leaderAddress,
		Brokers: brokers,
		Topic: topic,
		ConsumerGroup: consumerGroup,
	}, nil

}

func (k *Kafka) RunConsumerWithGroup(ctx context.Context, group_id string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:	k.Brokers,
		GroupID:	group_id,
		Topic:		k.Topic,
	})
	slog.Info("%v", r)
	go func() {
		for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			break
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