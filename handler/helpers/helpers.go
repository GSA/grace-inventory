package helpers

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/configservice/configserviceiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/aws/aws-sdk-go/service/glacier/glacieriface"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
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
func Alarms(svc cloudwatchiface.CloudWatchAPI) ([]*cloudwatch.MetricAlarm, error) {
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
func ConfigRules(svc configserviceiface.ConfigServiceAPI) ([]*configservice.ConfigRule, error) {
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
func LoadBalancers(svc elbv2iface.ELBV2API) ([]*elbv2.LoadBalancer, error) {
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
func Keys(svc kmsiface.KMSAPI) ([]*KmsKey, error) {
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
func listKeys(svc kmsiface.KMSAPI) ([]*kms.KeyListEntry, error) {
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
func getKeyDescriptions(svc kmsiface.KMSAPI, keyList []*kms.KeyListEntry) ([]*kms.KeyMetadata, error) {
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
func listKeyAliases(svc kmsiface.KMSAPI) ([]*kms.AliasListEntry, error) {
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

// Parameters ... pages through DescribeParametersPages to get SSM Parameters
func Parameters(svc ssmiface.SSMAPI) ([]*ssm.ParameterMetadata, error) {
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
