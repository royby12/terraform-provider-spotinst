package managed_instance_aws_compute

import "github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"

const (
	SubnetIds commons.FieldName = "subnet_ids"
	VpcId     commons.FieldName = "vpc_id"
	ElasticIp commons.FieldName = "elastic_ip"
	PrivateIp commons.FieldName = "private_ip"
)
