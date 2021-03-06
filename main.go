package main

import (
	"fmt"
	"math/rand"
	_ "math/rand"
	"sync"
	"time"
)

type MessageEvent struct {
	Message   string
	TopicID   string
	MessageID int
}

type SubscriptionID chan MessageEvent

type SubscriptionIDChannelSlice []SubscriptionID

type PubSub struct {
	// map topic to subscriptionIDChannels
	topicToSubscriptionIDsMap map[string]SubscriptionIDChannelSlice
	// map subscriptionID to topic
	subscriptionIDtoTopicMap map[SubscriptionID]string
	rm                       sync.RWMutex
}

var pb = &PubSub{
	topicToSubscriptionIDsMap: map[string]SubscriptionIDChannelSlice{},
	subscriptionIDtoTopicMap:  map[SubscriptionID]string{},
}

func (pb *PubSub) CreateTopic(topicID string) {
	pb.rm.Lock()
	if _, found := pb.topicToSubscriptionIDsMap[topicID]; !found {
		pb.topicToSubscriptionIDsMap[topicID] = make([]SubscriptionID, 0)
	} else {
		fmt.Printf("Topic %s already exists", topicID)
	}
	pb.rm.Unlock()
}

func (pb *PubSub) DeleteTopic(topicID string) {
	pb.rm.Lock()
	if _, found := pb.topicToSubscriptionIDsMap[topicID]; found {
		delete(pb.topicToSubscriptionIDsMap, topicID)
	} else {
		fmt.Printf("Topic %s does not exist", topicID)
	}
	pb.rm.Unlock()
}
func contains(s []SubscriptionID, e SubscriptionID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (pb *PubSub) AddSubscription(topicID string, subscriptionId SubscriptionID) {
	pb.rm.Lock()
	if prev, found := pb.topicToSubscriptionIDsMap[topicID]; found {
		if contains(prev, subscriptionId) {
			fmt.Printf("SubscriptionID %s already exists", subscriptionId)
		} else {
			pb.topicToSubscriptionIDsMap[topicID] = append(prev, subscriptionId)
			if pb.subscriptionIDtoTopicMap != nil {
				pb.subscriptionIDtoTopicMap[subscriptionId] = topicID
			} else {
				pb.subscriptionIDtoTopicMap = make(map[SubscriptionID]string, 0)
				pb.subscriptionIDtoTopicMap[subscriptionId] = topicID
			}
		}
	} else {
		fmt.Printf("Topic %s does not exist", topicID)
	}
	pb.rm.Unlock()
}

func (pb *PubSub) DeleteSubscription(subscriptionId SubscriptionID) {
	pb.rm.Lock()
	if prev, found := pb.subscriptionIDtoTopicMap[subscriptionId]; found {
		delete(pb.topicToSubscriptionIDsMap, prev)
		delete(pb.subscriptionIDtoTopicMap, subscriptionId)
		delete(sb.subscriptionIDToSubscriberFunc, subscriptionId)
		close(subscriptionId)
	} else {
		fmt.Printf("SubscriptionID %s does not exist", subscriptionId)
	}
	pb.rm.Unlock()
}

func (pb *PubSub) Publish(topicID string, message string) {
	pb.rm.RLock()
	if chans, found := pb.topicToSubscriptionIDsMap[topicID]; found {
		channels := append(SubscriptionIDChannelSlice{}, chans...)
		go func(message MessageEvent, subscriptionIDChannelSlices SubscriptionIDChannelSlice) {
			for _, ch := range subscriptionIDChannelSlices {
				ch <- message
			}
		}(MessageEvent{Message: message, TopicID: topicID, MessageID: rand.Int()}, channels)
	}
	pb.rm.RUnlock()
}
func PublishIterate(topicID string, message string) {
	i := 0
	for {
		pb.Publish(topicID, fmt.Sprint(message, i))
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		i = i + 1
	}
}

func (sb *Subscriber) Subscribe(subscriptionID SubscriptionID, SubscriberFunc1 SubscriberFunc) {
	sb.rm.Lock()
	sb.subscriptionIDToSubscriberFunc[subscriptionID] = SubscriberFunc1
	sb.rm.Unlock()
}

func (sb *Subscriber) UnSubscribe(subscriptionID SubscriptionID) {
	pb.rm.Lock()
	if _, found := sb.subscriptionIDToSubscriberFunc[subscriptionID]; found {
		delete(sb.subscriptionIDToSubscriberFunc, subscriptionID)
	} else {
		fmt.Printf("SubscriptionID %s does not exist", subscriptionID)
	}
	pb.rm.Unlock()
}

func (pb *PubSub) Ack(messageID int) {
	pb.rm.Lock()
	//messageEvent := <-subscriptionID
	fmt.Println("MessageID: ", messageID, " has been received and processed")
	fmt.Println()
	pb.rm.Unlock()
}

type Subscriber struct {
	subscriptionIDToSubscriberFunc map[SubscriptionID]SubscriberFunc
	rm                             sync.RWMutex
}

type SubscriberFunc func(subscriptionID SubscriptionID)

func SubscriberFunc1(subscriptionID SubscriptionID) {
	messageEvent := <-subscriptionID
	message := messageEvent.Message
	messageID := messageEvent.MessageID
	topicID := messageEvent.TopicID
	fmt.Println("Message: %s", message)
	fmt.Println("MessageID: ", messageID)
	fmt.Println("TopicID: %s", topicID)
	pb.Ack(messageID)
}

var sb = &Subscriber{
	subscriptionIDToSubscriberFunc: map[SubscriptionID]SubscriberFunc{},
}

func main() {
	subscriptionID1 := make(SubscriptionID)

	pb.CreateTopic("topic1")
	pb.AddSubscription("topic1", subscriptionID1)
	sb.Subscribe(subscriptionID1, SubscriberFunc1)

	go PublishIterate("topic1", "message")
	for {
		go SubscriberFunc1(subscriptionID1)

	}
}
