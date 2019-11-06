package managed_instance_aws_integrations

import "github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"

type BalancerType string

const (
	BalancerTypeClassic         BalancerType = "CLASSIC"
	BalancerTypeTargetGroup     BalancerType = "TARGET_GROUP"
	BalancerTypeMultaiTargetSet BalancerType = "MULTAI_TARGET_SET"
)

const (
	// - ROUTE53 -------------------------
	IntegrationRoute53 commons.FieldName = "integration_route53"
	Domains            commons.FieldName = "domains"
	HostedZoneId       commons.FieldName = "hosted_zone_id"
	SpotinstAcctID     commons.FieldName = "spotinst_acct_id"
	RecordSets         commons.FieldName = "record_sets"
	UsePublicIP        commons.FieldName = "use_public_ip"
	Name               commons.FieldName = "name"
	// -----------------------------------

	ElasticLoadBalancers

	//TODO sali ADD load balncer
)
