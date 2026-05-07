// 浣滅敤锛氭秷鎭槦鍒楅€傞厤锛堝 Kafka / RabbitMQ锛夌ず渚嬶紝鐢ㄤ簬鍙戝竷鎴栬闃呬簨浠?
package messaging

type KafkaPublisher struct {}

func (p *KafkaPublisher) Publish(topic string, message []byte) error {
	return nil
}
