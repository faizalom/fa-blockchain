package grpcclient

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"

	"fa-blockchain/src/utils"

	pb "github.com/faizalom/grpc-mq/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc/status"
)

var Client pb.MessageBrokerClient

func init() {
	// Initialize the gRPC client
	conn, err := DialServer()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	Client = pb.NewMessageBrokerClient(conn)
}

func DialServer() (*grpc.ClientConn, error) {
	tls := false // change that to true if needed
	opts := []grpc.DialOption{}

	if tls {
		certFile := "../../ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatalf("Error while loading CA trust certificate: %v\n", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		creds := grpc.WithTransportCredentials(insecure.NewCredentials())
		opts = append(opts, creds)
	}

	// conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	return grpc.NewClient("localhost:50051", opts...)
}

func NewSubscription(topic string) (pb.MessageBroker_SubscribeClient, error) {
	// Use a context without timeout for long-lived streaming
	ctx := context.Background()

	// Subscribe to the topic and receive a stream
	stream, err := Client.Subscribe(ctx, &pb.SubscriptionRequest{
		Topic:        topic,
		SubscriberId: os.Getenv("DEVICE_ID"),
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			log.Printf("Error message from server: %d %v\n", respErr.Code(), respErr.Message())
		} else {
			log.Printf("Error subscribing to topic: %v", err)
		}
	} else {
		fmt.Println("Successfully subscribed to topic:", topic)
	}

	return stream, err
}

func NewPublish(topic string, contentByte []byte, eventID string) (string, error) {
	if eventID == "" {
		// Generate a new event ID if not provided
		eventID = utils.NewEventId()
	}
	ctx := context.Background()

	pbm := &pb.Message{
		Topic:     topic,
		SenderId:  os.Getenv("DEVICE_ID"),
		EventId:   &eventID,
		Content:   &pb.Message_Binary{Binary: contentByte},
		Timestamp: time.Now().Unix(),
	}

	_, err := Client.Publish(ctx, pbm)
	if err == nil {
		log.Println("Message published successfully: ", pbm.Topic)
	}

	return eventID, err
}

func PublishVerifyTransaction(payload []byte) (string, error) {
	eventID, err := NewPublish("verify_trans", payload, "")
	if err == nil {
		return eventID, nil
	}

	respErr, ok := status.FromError(err)
	if ok {
		if respErr.Code() == codes.NotFound { // respErr.Message()
			return eventID, utils.ErrCreateTransTopicNotFound
		} else if respErr.Code() == codes.Unavailable { // respErr.Message()
			return eventID, utils.ErrMQBrokerUnavailable
		} else {
			log.Println(respErr.Err())
		}
	}

	log.Println("Error calling Publish: ", err)
	return eventID, err
}

func PublishCreateBlock(payload []byte, eventID string) (string, error) {
	eventID, err := NewPublish("create_block", payload, eventID)
	if err == nil {
		return eventID, nil
	}

	respErr, ok := status.FromError(err)
	if ok {
		if respErr.Code() == codes.NotFound { // respErr.Message()
			return eventID, fmt.Errorf("create_block: %v", utils.ErrCreateTransTopicNotFound)
		} else if respErr.Code() == codes.Unavailable { // respErr.Message()
			return eventID, utils.ErrMQBrokerUnavailable
		} else {
			log.Println(respErr.Err())
		}
	}

	log.Println("Error calling Publish: ", err)
	return eventID, err
}

func PublishNewMine(payload any) (string, error) {
	eventID, err := NewPublish("new_mining", Serialize(payload), "")
	if err == nil {
		return eventID, nil
	}

	respErr, ok := status.FromError(err)
	if ok {
		if respErr.Code() == codes.NotFound { // respErr.Message()
			return eventID, utils.ErrCreateTransTopicNotFound
		} else if respErr.Code() == codes.Unavailable { // respErr.Message()
			return eventID, utils.ErrMQBrokerUnavailable
		} else {
			log.Println(respErr.Err())
		}
	}

	log.Println("Error calling Publish: ", err)
	return eventID, err
}

func PublishVerifyNewBlock(payload any) (string, error) {
	eventID, err := NewPublish("verify_block", Serialize(payload), "")
	if err == nil {
		return eventID, nil
	}

	respErr, ok := status.FromError(err)
	if ok {
		if respErr.Code() == codes.NotFound { // respErr.Message()
			return eventID, utils.ErrCreateTransTopicNotFound
		} else if respErr.Code() == codes.Unavailable { // respErr.Message()
			return eventID, utils.ErrMQBrokerUnavailable
		} else {
			log.Println(respErr.Err())
		}
	}

	log.Println("Error calling Publish: ", err)
	return eventID, err
}

// func PublishVerifiedBlock(block blockchain.Block, blockStatus bool) (string, error) {
// 	date := struct {
// 		blockchain.Block
// 		Status bool
// 	}{
// 		block,
// 		blockStatus,
// 	}

// 	eventID, err := NewPublish("verified_block", Serialize(date), "")
// 	if err == nil {
// 		return eventID, nil
// 	}

// 	respErr, ok := status.FromError(err)
// 	if ok {
// 		if respErr.Code() == codes.NotFound { // respErr.Message()
// 			return eventID, utils.ErrCreateTransTopicNotFound
// 		} else if respErr.Code() == codes.Unavailable { // respErr.Message()
// 			return eventID, utils.ErrMQBrokerUnavailable
// 		} else {
// 			log.Println(respErr.Err())
// 		}
// 	}

// 	log.Println("Error calling Publish: ", err)
// 	return eventID, err
// }

func Serialize(outs any) []byte {
	var buffer bytes.Buffer

	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	if err != nil {
		// WIP
		log.Fatalf("Failed to encode data: %v", err)
	}
	return buffer.Bytes()
}
