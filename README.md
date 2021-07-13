# PubSub-System
Implementation of a PubSub system in Golang

There are two components - PubSub as pb, Subscriber as pb

Below methods are supported:

pb.CreateTopic(topicID)
pb.DeleteTopic(TopicID)
pb.AddSubscription(topicID,SubscriptionID); Creates and adds subscription with id SubscriptionID to topicName.
pb.DeleteSubscription(SubscriptionID)
pb.Publish(topicID, message); publishes the message on given topic

pb.Subscribe(SubscriptionID, SubscriberFunc); SubscriberFunc is the subscriber which is executed for each message of subscription.
pb.UnSubscribe(SubscriptionID)
