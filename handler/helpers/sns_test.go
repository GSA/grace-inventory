package helpers

import (
	"reflect"
	"testing"

	"github.com/GSA/grace-inventory-lambda/handler/inv"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type mockSnsClient struct {
	snsiface.SNSAPI
}

func (m mockSnsClient) ListSubscriptionsPages(in *sns.ListSubscriptionsInput, fn func(*sns.ListSubscriptionsOutput, bool) bool) error {
	fn(&sns.ListSubscriptionsOutput{Subscriptions: []*sns.Subscription{{}}}, true)
	return nil
}

func (m mockSnsClient) ListTopicsPages(in *sns.ListTopicsInput, fn func(*sns.ListTopicsOutput, bool) bool) error {
	fn(&sns.ListTopicsOutput{Topics: []*sns.Topic{{}}}, true)
	return nil
}

func (m mockSnsClient) GetTopicAttributes(in *sns.GetTopicAttributesInput) (*sns.GetTopicAttributesOutput, error) {
	return &sns.GetTopicAttributesOutput{Attributes: map[string]*string{}}, nil
}

// func Subscriptions(svc snsiface.SNSAPI) ([]*sns.Subscription, error) {
func TestSubscriptions(t *testing.T) {
	expected := []*sns.Subscription{{}}
	svc := mockSnsClient{}
	got, err := Subscriptions(svc)
	if err != nil {
		t.Fatalf("Subscriptions() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Subscriptions() failed.\nExpected %#v (%T)\nGot: %#v (%T)\n", expected, expected, got, got)
	}
	_, err = inv.TypeToSheet(expected)
	if err != nil {
		t.Fatalf("inv.TypeToSheet failed: %v", err)
	}
}

// func Topics(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*SnsTopic, error) {
func TestTopics(t *testing.T) {
	expected := []*SnsTopic{{}}
	svc := mockSnsClient{}
	got, err := Topics(svc)
	if err != nil {
		t.Fatalf("Topics() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Topics() failed.\nExpected %#v (%T)\nGot: %#v (%T)\n", expected, expected, got, got)
	}
	_, err = inv.TypeToSheet(expected)
	if err != nil {
		t.Fatalf("inv.TypeToSheet failed: %v", err)
	}
}
