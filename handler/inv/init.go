package inv

import (
	"os"
	"regexp"

	"github.com/GSA/grace-inventory-lambda/handler/spreadsheet"
)

// Sheet name constants
const (
	SheetRoles          = "Roles"
	SheetAccounts       = "Accounts"
	SheetGroups         = "Groups"
	SheetPolicies       = "Policies"
	SheetUsers          = "Users"
	SheetBuckets        = "Buckets"
	SheetInstances      = "Instances"
	SheetImages         = "Images"
	SheetVolumes        = "Volumes"
	SheetSnapshots      = "Snapshots"
	SheetVpcs           = "VPCs"
	SheetSubnets        = "Subnets"
	SheetSecurityGroups = "SecurityGroups"
	SheetAddresses      = "Addresses"
	SheetKeyPairs       = "KeyPairs"
	SheetStacks         = "Stacks"
	SheetAlarms         = "Alarms"
	SheetConfigRules    = "ConfigRules"
	SheetLoadBalancers  = "LoadBlancers"
	SheetVaults         = "Vaults"
	SheetKeys           = "Keys"
	SheetDBInstances    = "DBInstances"
	SheetDBSnapshots    = "DBSnapshots"
	SheetSecrets        = "Secrets"
	SheetSubscriptions  = "Subscriptions"
	SheetTopics         = "Topics"
	SheetParameters     = "Parameters"
)

func init() {
	accountsInfo := os.Getenv("accounts_info")
	r := regexp.MustCompile(`^\d{12}`)
	if accountsInfo == "self" || r.MatchString(accountsInfo) {
		spreadsheet.RegisterSheet(SheetAccounts, func() *spreadsheet.Sheet {
			return &spreadsheet.Sheet{Name: "Accounts", Columns: []*spreadsheet.Column{
				{FriendlyName: "Alias", FieldName: "Name"},
				{FriendlyName: "Id", FieldName: "Id"},
			}}
		})
	} else {
		spreadsheet.RegisterSheet(SheetAccounts, func() *spreadsheet.Sheet {
			return &spreadsheet.Sheet{Name: "Accounts", Columns: []*spreadsheet.Column{
				{FriendlyName: "Name", FieldName: "Name"},
				{FriendlyName: "Id", FieldName: "Id"},
				{FriendlyName: "Status", FieldName: "Status"},
				{FriendlyName: "Email", FieldName: "Email"},
				{FriendlyName: "JoinedMethod", FieldName: "JoinedMethod"},
				{FriendlyName: "JoinedTimestamp", FieldName: "JoinedTimestamp"},
				{FriendlyName: "Arn", FieldName: "Arn"},
			}}
		})
	}
	spreadsheet.RegisterSheet(SheetRoles, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "IAM Roles", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "RoleName", FieldName: "RoleName"},
			{FriendlyName: "RoleId", FieldName: "RoleId"},
			{FriendlyName: "Description", FieldName: "Description"},
			{FriendlyName: "CreateDate", FieldName: "CreateDate"},
		}}
	})
	spreadsheet.RegisterSheet(SheetGroups, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "IAM Groups", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "GroupName", FieldName: "GroupName"},
			{FriendlyName: "GroupId", FieldName: "GroupId"},
			{FriendlyName: "CreateDate", FieldName: "CreateDate"},
		}}
	})
	spreadsheet.RegisterSheet(SheetPolicies, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "IAM Policies", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "PolicyName", FieldName: "PolicyName"},
			{FriendlyName: "PolicyId", FieldName: "PolicyId"},
			{FriendlyName: "CreateDate", FieldName: "CreateDate"},
		}}
	})
	spreadsheet.RegisterSheet(SheetUsers, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "IAM Users", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "UserName", FieldName: "UserName"},
			{FriendlyName: "UserId", FieldName: "UserId"},
			{FriendlyName: "CreateDate", FieldName: "CreateDate"},
		}}
	})
	spreadsheet.RegisterSheet(SheetBuckets, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "S3 Buckets", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Name", FieldName: "Name"},
			{FriendlyName: "CreateDate", FieldName: "CreationDate"},
		}}
	})
	spreadsheet.RegisterSheet(SheetInstances, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "EC2 Instances", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "Tags"},
			{FriendlyName: "InstanceId", FieldName: "InstanceId"},
			{FriendlyName: "PrivateIpAddress", FieldName: "PrivateIpAddress"},
			{FriendlyName: "PublicIpAddress", FieldName: "PublicIpAddress"},
			{FriendlyName: "State", FieldName: "State"},
			{FriendlyName: "ImageId", FieldName: "ImageId"},
			{FriendlyName: "LaunchTime", FieldName: "LaunchTime"},
		}}
	})
	spreadsheet.RegisterSheet(SheetImages, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Images", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "Tags"},
			{FriendlyName: "AMI Name", FieldName: "Name"},
			{FriendlyName: "ImageId", FieldName: "ImageId"},
			{FriendlyName: "Description", FieldName: "Description"},
			{FriendlyName: "State", FieldName: "State"},
			{FriendlyName: "CreationDate", FieldName: "CreationDate"},
		}}
	})
	spreadsheet.RegisterSheet(SheetVolumes, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Volumes", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "VolumeId", FieldName: "VolumeId"},
			{FriendlyName: "State", FieldName: "State"},
			{FriendlyName: "Size", FieldName: "Size"},
			{FriendlyName: "VolumeType", FieldName: "VolumeType"},
			{FriendlyName: "Encrypted", FieldName: "Encrypted"},
			{FriendlyName: "CreateTime", FieldName: "CreateTime"},
		}}
	})
	spreadsheet.RegisterSheet(SheetSnapshots, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Snapshots", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "Tags"},
			{FriendlyName: "SnapshotId", FieldName: "SnapshotId"},
			{FriendlyName: "Description", FieldName: "Description"},
			{FriendlyName: "State", FieldName: "State"},
			{FriendlyName: "VolumeId", FieldName: "VolumeId"},
			{FriendlyName: "VolumeSize", FieldName: "VolumeSize"},
			{FriendlyName: "StartTime", FieldName: "StartTime"},
		}}
	})
	spreadsheet.RegisterSheet(SheetVpcs, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "VPCs", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "Tags"},
			{FriendlyName: "VpcId", FieldName: "VpcId"},
			{FriendlyName: "State", FieldName: "State"},
			{FriendlyName: "CidrBlock", FieldName: "CidrBlock"},
			{FriendlyName: "DhcpOptionsId", FieldName: "DhcpOptionsId"},
		}}
	})
	spreadsheet.RegisterSheet(SheetSubnets, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Subnets", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "Tags"},
			{FriendlyName: "SubnetId", FieldName: "SubnetId"},
			{FriendlyName: "VpcId", FieldName: "VpcId"},
			{FriendlyName: "State", FieldName: "State"},
			{FriendlyName: "CidrBlock", FieldName: "CidrBlock"},
			{FriendlyName: "AvailabilityZone", FieldName: "AvailabilityZone"},
		}}
	})
	spreadsheet.RegisterSheet(SheetSecurityGroups, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "SecurityGroups", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "GroupName", FieldName: "GroupName"},
			{FriendlyName: "GroupId", FieldName: "GroupId"},
			{FriendlyName: "Description", FieldName: "Description"},
			{FriendlyName: "VpcId", FieldName: "VpcId"},
		}}
	})
	spreadsheet.RegisterSheet(SheetAddresses, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "EC2 IP Addresses", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "AllocationId", FieldName: "AllocationId"},
			{FriendlyName: "AssociationId", FieldName: "AssociationId"},
			{FriendlyName: "Domain", FieldName: "Domain"},
			{FriendlyName: "InstanceId", FieldName: "InstanceId"},
			{FriendlyName: "NetworkInterfaceId", FieldName: "NetworkInterfaceId"},
			{FriendlyName: "NetworkInterfaceOwnerId", FieldName: "NetworkInterfaceOwnerId"},
			{FriendlyName: "PrivateIpAddress", FieldName: "PrivateIpAddress"},
			{FriendlyName: "PublicIp", FieldName: "PublicIp"},
			{FriendlyName: "PublicIpv4Pool", FieldName: "PublicIpv4Pool"},
		}}
	})
	spreadsheet.RegisterSheet(SheetKeyPairs, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Key Pairs", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "KeyName", FieldName: "KeyName"},
			{FriendlyName: "KeyFingerprint", FieldName: "KeyFingerprint"},
		}}
	})
	spreadsheet.RegisterSheet(SheetStacks, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "CloudFormation Stacks", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "StackName", FieldName: "StackName"},
			{FriendlyName: "Description", FieldName: "Description"},
			{FriendlyName: "RootId", FieldName: "RootId"},
			{FriendlyName: "StackId", FieldName: "StackId"},
			{FriendlyName: "ParentId", FieldName: "ParentId"},
			{FriendlyName: "RoleARN", FieldName: "RoleARN"},
			{FriendlyName: "CreationTime", FieldName: "CreationTime"},
			{FriendlyName: "LastUpdatedTime", FieldName: "LastUpdatedTime"},
			{FriendlyName: "DeletionTime", FieldName: "DeletionTime"},
			{FriendlyName: "ChangeSetId", FieldName: "ChangeSetId"},
			{FriendlyName: "StackStatus", FieldName: "StackStatus"},
			{FriendlyName: "StackStatusReason", FieldName: "StackStatusReason"},
			{FriendlyName: "TimeoutInMinutes", FieldName: "TimeoutInMinutes"},
		}}
	})
	spreadsheet.RegisterSheet(SheetAlarms, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Alarms", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "AlarmName"},
			{FriendlyName: "Description", FieldName: "AlarmDescription"},
			{FriendlyName: "AlarmArn", FieldName: "AlarmArn"},
			{FriendlyName: "ActionsEnabled", FieldName: "ActionsEnabled"},
			{FriendlyName: "Updated", FieldName: "AlarmConfigurationUpdatedTimestamp"},
			{FriendlyName: "ComparisonOperator", FieldName: "ComparisonOperator"},
			{FriendlyName: "DatapointsToAlarm", FieldName: "DatapointsToAlarm"},
			{FriendlyName: "EvaluateLowSampleCountPercentile", FieldName: "EvaluateLowSampleCountPercentile"},
			{FriendlyName: "EvaluationPeriods", FieldName: "EvaluationPeriods"},
			{FriendlyName: "ExtendedStatistic", FieldName: "ExtendedStatistic"},
			{FriendlyName: "MetricName", FieldName: "MetricName"},
			{FriendlyName: "Namespace", FieldName: "Namespace"},
			{FriendlyName: "Period", FieldName: "Period"},
			{FriendlyName: "StateReason", FieldName: "StateReason"},
			{FriendlyName: "StateReasonData", FieldName: "StateReasonData"},
			{FriendlyName: "StateUpdatedTimestamp", FieldName: "StateUpdatedTimestamp"},
			{FriendlyName: "StateValue", FieldName: "StateValue"},
			{FriendlyName: "Statistic", FieldName: "Statistic"},
			{FriendlyName: "Threshold", FieldName: "Threshold"},
			{FriendlyName: "TreatMissingData", FieldName: "TreatMissingData"},
			{FriendlyName: "Unit", FieldName: "Unit"},
		}}
	})
	spreadsheet.RegisterSheet(SheetConfigRules, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Config Rules", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "ConfigRuleName"},
			{FriendlyName: "Description", FieldName: "Description"},
			{FriendlyName: "ConfigRuleId", FieldName: "ConfigRuleId"},
			{FriendlyName: "ConfigRuleArn", FieldName: "ConfigRuleArn"},
			{FriendlyName: "ConfigRuleState", FieldName: "ConfigRuleState"},
			{FriendlyName: "CreatedBy", FieldName: "CreatedBy"},
			{FriendlyName: "InputParameters", FieldName: "InputParameters"},
			{FriendlyName: "MaximumExecutionFrequency", FieldName: "MaximumExecutionFrequency"},
		}}
	})
	spreadsheet.RegisterSheet(SheetLoadBalancers, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Load Balancers", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "LoadBalancerName"},
			{FriendlyName: "DNSName", FieldName: "DNSName"},
			{FriendlyName: "CanonicalHostedZoneId", FieldName: "CanonicalHostedZoneId"},
			{FriendlyName: "CreatedTime", FieldName: "CreatedTime"},
			{FriendlyName: "IpAddressType", FieldName: "IpAddressType"},
			{FriendlyName: "LoadBalancerArn", FieldName: "LoadBalancerArn"},
			{FriendlyName: "Scheme", FieldName: "Scheme"},
			{FriendlyName: "State", FieldName: "State"},
			{FriendlyName: "Type", FieldName: "Type"},
			{FriendlyName: "VpcId", FieldName: "VpcId"},
		}}
	})
	spreadsheet.RegisterSheet(SheetVaults, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Glacier Vaults", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "VaultName"},
			{FriendlyName: "VaultARN", FieldName: "VaultARN"},
			{FriendlyName: "SizeInBytes", FieldName: "SizeInBytes"},
			{FriendlyName: "NumberOfArchives", FieldName: "NumberOfArchives"},
			{FriendlyName: "CreationDate", FieldName: "CreationDate"},
			{FriendlyName: "LastInventoryDate", FieldName: "LastInventoryDate"},
		}}
	})
	spreadsheet.RegisterSheet(SheetKeys, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "KMS Keys", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "AliasName", FieldName: "AliasName"},
			{FriendlyName: "Arn", FieldName: "Arn"},
			{FriendlyName: "KeyId", FieldName: "KeyID"},
			{FriendlyName: "CloudHsmClusterId", FieldName: "CloudHsmClusterID"},
			{FriendlyName: "CreationDate", FieldName: "CreationDate"},
			{FriendlyName: "CustomKeyStoreId", FieldName: "CustomKeyStoreID"},
			{FriendlyName: "DeletionDate", FieldName: "DeletionDate"},
			{FriendlyName: "Description", FieldName: "Description"},
			{FriendlyName: "Enabled", FieldName: "Enabled"},
			{FriendlyName: "ExpirationModel", FieldName: "ExpirationModel"},
			{FriendlyName: "KeyManager", FieldName: "KeyManager"},
			{FriendlyName: "KeyState", FieldName: "KeyState"},
			{FriendlyName: "KeyUsage", FieldName: "KeyUsage"},
			{FriendlyName: "Origin", FieldName: "Origin"},
			{FriendlyName: "ValidTo", FieldName: "ValidTo"},
		}}
	})
	spreadsheet.RegisterSheet(SheetDBInstances, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "RDS DB Instances", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "AvailabilityZone", FieldName: "AvailabilityZone"},
			{FriendlyName: "DBClusterIdentifier", FieldName: "DBClusterIdentifier"},
			{FriendlyName: "DBInstanceIdentifier", FieldName: "DBInstanceIdentifier"},
			{FriendlyName: "DBName", FieldName: "DBName"},
			{FriendlyName: "Engine", FieldName: "Engine"},
			{FriendlyName: "EngineVersion", FieldName: "EngineVersion"},
			{FriendlyName: "Endpoint", FieldName: "Endpoint"},
			{FriendlyName: "DBInstanceArn", FieldName: "DBInstanceArn"},
			{FriendlyName: "DBInstanceClass", FieldName: "DBInstanceClass"},
			{FriendlyName: "DBInstanceStatus", FieldName: "DBInstanceStatus"},
			{FriendlyName: "MultiAZ", FieldName: "MultiAZ"},
			{FriendlyName: "PubliclyAccessible", FieldName: "PubliclyAccessible"},
			{FriendlyName: "StorageEncrypted", FieldName: "StorageEncrypted"},
			{FriendlyName: "InstanceCreateTime", FieldName: "InstanceCreateTime"},
		}}
	})
	spreadsheet.RegisterSheet(SheetDBSnapshots, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "RDS DB Snapshots", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "AllocatedStorage", FieldName: "AllocatedStorage"},
			{FriendlyName: "AvailabilityZone", FieldName: "AvailabilityZone"},
			{FriendlyName: "DBInstanceIdentifier", FieldName: "DBInstanceIdentifier"},
			{FriendlyName: "DBSnapshotArn", FieldName: "DBSnapshotArn"},
			{FriendlyName: "DBSnapshotIdentifier", FieldName: "DBSnapshotIdentifier"},
			{FriendlyName: "DbiResourceId", FieldName: "DbiResourceId"},
			{FriendlyName: "Encrypted", FieldName: "Encrypted"},
			{FriendlyName: "Engine", FieldName: "Engine"},
			{FriendlyName: "EngineVersion", FieldName: "EngineVersion"},
			{FriendlyName: "IAMDatabaseAuthenticationEnabled", FieldName: "IAMDatabaseAuthenticationEnabled"},
			{FriendlyName: "InstanceCreateTime", FieldName: "InstanceCreateTime"},
			{FriendlyName: "SnapshotCreateTime", FieldName: "SnapshotCreateTime"},
			{FriendlyName: "SourceDBSnapshotIdentifier", FieldName: "SourceDBSnapshotIdentifier"},
			{FriendlyName: "SourceRegion", FieldName: "SourceRegion"},
			{FriendlyName: "Status", FieldName: "Status"},
			{FriendlyName: "StorageType", FieldName: "StorageType"},
			{FriendlyName: "VpcId", FieldName: "VpcId"},
		}}
	})
	spreadsheet.RegisterSheet(SheetSecrets, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "Secrets", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "Name"},
			{FriendlyName: "Description", FieldName: "Description"},
			{FriendlyName: "ARN", FieldName: "ARN"},
			{FriendlyName: "DeletedDate", FieldName: "DeletedDate"},
			{FriendlyName: "KmsKeyId", FieldName: "KmsKeyId"},
			{FriendlyName: "LastAccessedDate", FieldName: "LastAccessedDate"},
			{FriendlyName: "LastChangedDate", FieldName: "LastChangedDate"},
			{FriendlyName: "LastRotatedDate", FieldName: "LastRotatedDate"},
			{FriendlyName: "RotationEnabled", FieldName: "RotationEnabled"},
		}}
	})
	spreadsheet.RegisterSheet(SheetSubscriptions, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "SNS Subscriptions", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Endpoint", FieldName: "Endpoint"},
			{FriendlyName: "Owner", FieldName: "Owner"},
			{FriendlyName: "Protocol", FieldName: "Protocol"},
			{FriendlyName: "SubscriptionArn", FieldName: "SubscriptionArn"},
			{FriendlyName: "TopicArn", FieldName: "TopicArn"},
		}}
	})
	spreadsheet.RegisterSheet(SheetTopics, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "SNS Topics", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "DisplayName"},
			{FriendlyName: "TopicArn", FieldName: "TopicArn"},
			{FriendlyName: "Owner", FieldName: "Owner"},
			{FriendlyName: "SubscriptionsPending", FieldName: "SubscriptionsPending"},
			{FriendlyName: "SubscriptionsConfirmed", FieldName: "SubscriptionsConfirmed"},
			{FriendlyName: "SubscriptionsDeleted", FieldName: "SubscriptionsDeleted"},
			{FriendlyName: "DeliveryPolicy", FieldName: "DeliveryPolicy"},
			{FriendlyName: "EffectiveDeliveryPolicy", FieldName: "EffectiveDeliveryPolicy"},
		}}
	})
	spreadsheet.RegisterSheet(SheetParameters, func() *spreadsheet.Sheet {
		return &spreadsheet.Sheet{Name: "SSM Parameters", Columns: []*spreadsheet.Column{
			{FriendlyName: "Account", FieldName: ""},
			{FriendlyName: "Region", FieldName: ""},
			{FriendlyName: "Name", FieldName: "Name"},
			{FriendlyName: "Description", FieldName: "Description"},
			{FriendlyName: "KeyId", FieldName: "KeyId"},
			{FriendlyName: "AllowedPattern", FieldName: "AllowedPattern"},
			{FriendlyName: "Tier", FieldName: "Tier"},
			{FriendlyName: "Type", FieldName: "Type"},
			{FriendlyName: "Version", FieldName: "Version"},
			{FriendlyName: "LastModifiedDate", FieldName: "LastModifiedDate"},
			{FriendlyName: "LastModifiedUser", FieldName: "LastModifiedUser"},
		}}
	})
}
