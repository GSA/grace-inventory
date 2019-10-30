package helpers

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glacier/glacieriface"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var self = []*string{aws.String("self")}

// KmsKey ... extends kms.KeyMetadata to add AliasName
type KmsKey struct {

	// The Amazon Resource Name (ARN) of the CMK. For examples, see AWS Key Management
	// Service (AWS KMS) (https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html#arn-syntax-kms)
	// in the Example ARNs section of the AWS General Reference.
	Arn *string `min:"20" type:"string"`

	// The cluster ID of the AWS CloudHSM cluster that contains the key material
	// for the CMK. When you create a CMK in a custom key store (https://docs.aws.amazon.com/kms/latest/developerguide/custom-key-store-overview.html),
	// AWS KMS creates the key material for the CMK in the associated AWS CloudHSM
	// cluster. This value is present only when the CMK is created in a custom key
	// store.
	CloudHsmClusterID *string `min:"19" type:"string"`

	// The date and time when the CMK was created.
	CreationDate *time.Time `type:"timestamp"`

	// A unique identifier for the custom key store (https://docs.aws.amazon.com/kms/latest/developerguide/custom-key-store-overview.html)
	// that contains the CMK. This value is present only when the CMK is created
	// in a custom key store.
	CustomKeyStoreID *string `min:"1" type:"string"`

	// The date and time after which AWS KMS deletes the CMK. This value is present
	// only when KeyState is PendingDeletion.
	DeletionDate *time.Time `type:"timestamp"`

	// The description of the CMK.
	Description *string `type:"string"`

	// Specifies whether the CMK is enabled. When KeyState is Enabled this value
	// is true, otherwise it is false.
	Enabled *bool `type:"boolean"`

	// Specifies whether the CMK's key material expires. This value is present only
	// when Origin is EXTERNAL, otherwise this value is omitted.
	ExpirationModel *string `type:"string" enum:"ExpirationModelType"`

	// The globally unique identifier for the CMK.
	//
	// KeyId is a required field
	KeyID *string `min:"1" type:"string" required:"true"`

	// The manager of the CMK. CMKs in your AWS account are either customer managed
	// or AWS managed. For more information about the difference, see Customer Master
	// Keys (https://docs.aws.amazon.com/kms/latest/developerguide/concepts.html#master_keys)
	// in the AWS Key Management Service Developer Guide.
	KeyManager *string `type:"string" enum:"KeyManagerType"`

	// The state of the CMK.
	//
	// For more information about how key state affects the use of a CMK, see How
	// Key State Affects the Use of a Customer Master Key (https://docs.aws.amazon.com/kms/latest/developerguide/key-state.html)
	// in the AWS Key Management Service Developer Guide.
	KeyState *string `type:"string" enum:"KeyState"`

	// The cryptographic operations for which you can use the CMK. The only valid
	// value is ENCRYPT_DECRYPT, which means you can use the CMK to encrypt and
	// decrypt data.
	KeyUsage *string `type:"string" enum:"KeyUsageType"`

	// The source of the CMK's key material. When this value is AWS_KMS, AWS KMS
	// created the key material. When this value is EXTERNAL, the key material was
	// imported from your existing key management infrastructure or the CMK lacks
	// key material. When this value is AWS_CLOUDHSM, the key material was created
	// in the AWS CloudHSM cluster associated with a custom key store.
	Origin *string `type:"string" enum:"OriginType"`

	// The time at which the imported key material expires. When the key material
	// expires, AWS KMS deletes the key material and the CMK becomes unusable. This
	// value is present only for CMKs whose Origin is EXTERNAL and whose ExpirationModel
	// is KEY_MATERIAL_EXPIRES, otherwise this value is omitted.
	ValidTo *time.Time `type:"timestamp"`
	// String that contains the alias. This value begins with alias/.
	AliasName *string `min:"1" type:"string"`
}

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

// Buckets ... performs ListBuckets and returns all S3 buckets
func Buckets(svc s3iface.S3API) ([]*s3.Bucket, error) {
	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	return result.Buckets, err
}

// Stacks ... pages through DescribeStacksPages and returns all CloudFormation Stacks
func Stacks(svc cloudformationiface.CloudFormationAPI) ([]*cloudformation.Stack, error) {
	var results []*cloudformation.Stack
	err := svc.DescribeStacksPages(&cloudformation.DescribeStacksInput{},
		func(page *cloudformation.DescribeStacksOutput, lastPage bool) bool {
			results = append(results, page.Stacks...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Alarms ... pages through DescribeAlarmsPages and returns all CloudWatch Metric Alarms
func Alarms(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*cloudwatch.MetricAlarm, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := cloudwatch.New(cfg, &aws.Config{Credentials: cred})
	var results []*cloudwatch.MetricAlarm
	err := svc.DescribeAlarmsPages(&cloudwatch.DescribeAlarmsInput{},
		func(page *cloudwatch.DescribeAlarmsOutput, lastPage bool) bool {
			results = append(results, page.MetricAlarms...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// ConfigRules ... performs DescribeConfigRules and returns all Config Service ConfigRules
func ConfigRules(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*configservice.ConfigRule, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := configservice.New(cfg, &aws.Config{Credentials: cred})
	input := &configservice.DescribeConfigRulesInput{}
	result, err := svc.DescribeConfigRules(input)
	if err != nil {
		return nil, err
	}
	rules := result.ConfigRules
	for result.NextToken != nil {
		input.NextToken = result.NextToken
		result, err = svc.DescribeConfigRules(input)
		if err != nil {
			return nil, err
		}
		rules = append(rules, result.ConfigRules...)
	}
	return rules, nil
}

// LoadBalancers ... pages through DescribeLoadBalancersPages and returns all ELB v2 LoadBalancers
func LoadBalancers(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*elbv2.LoadBalancer, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := elbv2.New(cfg, &aws.Config{Credentials: cred})
	var results []*elbv2.LoadBalancer
	err := svc.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{},
		func(page *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
			results = append(results, page.LoadBalancers...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// GlacierSvc is used to call glacier functions
type GlacierSvc struct {
	Client glacieriface.GlacierAPI
}

// Vaults ... pages through ListVaultsPages and returns all Glacier Vaults
func (svc *GlacierSvc) Vaults() ([]*glacier.DescribeVaultOutput, error) {
	var results []*glacier.DescribeVaultOutput
	err := svc.Client.ListVaultsPages(&glacier.ListVaultsInput{},
		func(page *glacier.ListVaultsOutput, lastPage bool) bool {
			results = append(results, page.VaultList...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Keys ... pages over ListKeys results and returns all KMS Keys w/ AliasName
func Keys(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*KmsKey, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := kms.New(cfg, &aws.Config{Credentials: cred})
	keyList, err := listKeys(svc)
	if err != nil {
		return nil, err
	}
	keyDescriptions, err := getKeyDescriptions(svc, keyList)
	if err != nil {
		return nil, err
	}
	keyAliases, err := listKeyAliases(svc)
	if err != nil {
		return nil, err
	}
	return mergeKeyAliases(keyDescriptions, keyAliases)
}

// listKeys ... pages through ListKeysPages to get list of KeyIDs
func listKeys(svc *kms.KMS) ([]*kms.KeyListEntry, error) {
	var results []*kms.KeyListEntry
	err := svc.ListKeysPages(&kms.ListKeysInput{},
		func(page *kms.ListKeysOutput, lastPage bool) bool {
			results = append(results, page.Keys...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// getKeyDescriptions ... loops through list of KeyIds to get KeyMetadata (kms.DescribeKey)
func getKeyDescriptions(svc *kms.KMS, keyList []*kms.KeyListEntry) ([]*kms.KeyMetadata, error) {
	var keys []*kms.KeyMetadata
	for _, key := range keyList {
		input := &kms.DescribeKeyInput{KeyId: key.KeyId}
		result, err := svc.DescribeKey(input)
		if err != nil {
			return nil, err
		}
		keys = append(keys, result.KeyMetadata)
	}
	return keys, nil
}

// listKeyAliases ... pages over ListAliasesPages and returns list of Aliases
func listKeyAliases(svc *kms.KMS) ([]*kms.AliasListEntry, error) {
	var results []*kms.AliasListEntry
	err := svc.ListAliasesPages(&kms.ListAliasesInput{},
		func(page *kms.ListAliasesOutput, lastPage bool) bool {
			results = append(results, page.Aliases...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Merges AliasListEntries with KeyMetadata to add AliasName to KeyMetadata
func mergeKeyAliases(keyDescriptions []*kms.KeyMetadata, keyAliases []*kms.AliasListEntry) ([]*KmsKey, error) {
	var keys []*KmsKey
	for _, metadata := range keyDescriptions {
		key := KmsKey{
			AliasName:         getAliasName(aws.StringValue(metadata.KeyId), keyAliases),
			Arn:               metadata.Arn,
			KeyID:             metadata.KeyId,
			CloudHsmClusterID: metadata.CloudHsmClusterId,
			CreationDate:      metadata.CreationDate,
			CustomKeyStoreID:  metadata.CustomKeyStoreId,
			DeletionDate:      metadata.DeletionDate,
			Description:       metadata.Description,
			Enabled:           metadata.Enabled,
			ExpirationModel:   metadata.ExpirationModel,
			KeyManager:        metadata.KeyManager,
			KeyState:          metadata.KeyState,
			KeyUsage:          metadata.KeyUsage,
			Origin:            metadata.Origin,
			ValidTo:           metadata.ValidTo,
		}
		keys = append(keys, &key)
	}
	return keys, nil
}

func getAliasName(str string, keyAliases []*kms.AliasListEntry) *string {
	for _, v := range keyAliases {
		if str == aws.StringValue(v.TargetKeyId) {
			return v.AliasName
		}
	}
	return aws.String("")
}

// DBInstances ... pages through DescribeDBInstancesPages to get list of DBInstances
func DBInstances(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*rds.DBInstance, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := rds.New(cfg, &aws.Config{Credentials: cred})
	var results []*rds.DBInstance
	err := svc.DescribeDBInstancesPages(&rds.DescribeDBInstancesInput{},
		func(page *rds.DescribeDBInstancesOutput, lastPage bool) bool {
			results = append(results, page.DBInstances...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// DBSnapshots ... pages through DescribeDBSnapshotsPages to get list of DBSnapshots
func DBSnapshots(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*rds.DBSnapshot, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := rds.New(cfg, &aws.Config{Credentials: cred})
	var results []*rds.DBSnapshot
	err := svc.DescribeDBSnapshotsPages(&rds.DescribeDBSnapshotsInput{},
		func(page *rds.DescribeDBSnapshotsOutput, lastPage bool) bool {
			results = append(results, page.DBSnapshots...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Secrets ... pages through ListSecretsPages to get list of Secrets
func Secrets(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*secretsmanager.SecretListEntry, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := secretsmanager.New(cfg, &aws.Config{Credentials: cred})
	var results []*secretsmanager.SecretListEntry
	err := svc.ListSecretsPages(&secretsmanager.ListSecretsInput{},
		func(page *secretsmanager.ListSecretsOutput, lastPage bool) bool {
			results = append(results, page.SecretList...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Subscriptions ... pages through ListSubscriptionsPages to get list of Subscriptions
func Subscriptions(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*sns.Subscription, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := sns.New(cfg, &aws.Config{Credentials: cred})
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
func Topics(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*SnsTopic, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := sns.New(cfg, &aws.Config{Credentials: cred})
	topicList, err := listTopics(svc)
	if err != nil {
		return nil, err
	}
	return getTopicAttributes(svc, topicList)
}

// listTopics ... pages through ListTopicsPages to get list of TopicArns
func listTopics(svc *sns.SNS) ([]*sns.Topic, error) {
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
func getTopicAttributes(svc *sns.SNS, topicList []*sns.Topic) ([]*SnsTopic, error) {
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

// Parameters ... pages through DescribeParametersPages to get SSM Parameters
func Parameters(cfg client.ConfigProvider, cred *credentials.Credentials) ([]*ssm.ParameterMetadata, error) {
	if cfg == nil {
		return nil, errors.New("nil ConfigProvider")
	}
	svc := ssm.New(cfg, &aws.Config{Credentials: cred})
	var results []*ssm.ParameterMetadata
	err := svc.DescribeParametersPages(&ssm.DescribeParametersInput{},
		func(page *ssm.DescribeParametersOutput, lastPage bool) bool {
			results = append(results, page.Parameters...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}
