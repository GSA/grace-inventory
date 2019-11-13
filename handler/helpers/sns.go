package helpers

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

// SnsTopic ... struct definition for Attributes map in GetTopicAttributesOutput
type SnsTopic struct {
	DisplayName             *string
	TopicArn                *string
	Owner                   *string
	SubscriptionsPending    *string
	SubscriptionsConfirmed  *string
	SubscriptionsDeleted    *string
	DeliveryPolicy          *string
	EffectiveDeliveryPolicy *string
}

// Subscriptions ... pages through ListSubscriptionsPages to get list of Subscriptions
func Subscriptions(svc snsiface.SNSAPI) ([]*sns.Subscription, error) {
	var results []*sns.Subscription
	err := svc.ListSubscriptionsPages(&sns.ListSubscriptionsInput{},
		func(page *sns.ListSubscriptionsOutput, lastPage bool) bool {
			results = append(results, page.Subscriptions...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Topics ... pages over ListTopics results and returns all Topics parameters
func Topics(svc snsiface.SNSAPI) ([]*SnsTopic, error) {
	topicList, err := listTopics(svc)
	if err != nil {
		return nil, err
	}
	return getTopicAttributes(svc, topicList)
}

// listTopics ... pages through ListTopicsPages to get list of TopicArns
func listTopics(svc snsiface.SNSAPI) ([]*sns.Topic, error) {
	var results []*sns.Topic
	err := svc.ListTopicsPages(&sns.ListTopicsInput{},
		func(page *sns.ListTopicsOutput, lastPage bool) bool {
			results = append(results, page.Topics...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// getTopicAttributes ... loops through list of Topic ARNs to get Topic attributes GetTopicAttributes())
func getTopicAttributes(svc snsiface.SNSAPI, topicList []*sns.Topic) ([]*SnsTopic, error) {
	var topics []*SnsTopic
	for _, t := range topicList {
		input := &sns.GetTopicAttributesInput{TopicArn: t.TopicArn}
		result, err := svc.GetTopicAttributes(input)
		if err != nil {
			return nil, err
		}
		a := result.Attributes
		m := &SnsTopic{
			DisplayName:             a["DisplayName"],
			TopicArn:                a["TopicArn"],
			Owner:                   a["Owner"],
			SubscriptionsPending:    a["SubscriptionsPending"],
			SubscriptionsConfirmed:  a["SubscriptionsConfirmed"],
			SubscriptionsDeleted:    a["SubscriptionsDeleted"],
			DeliveryPolicy:          a["DeliveryPolicy"],
			EffectiveDeliveryPolicy: a["EffectiveDeliveryPolicy"],
		}
		topics = append(topics, m)
	}
	return topics, nil
}
