package managed_instance_aws_integrations

import "github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {
	SetupRoute53(fieldsMap)
	//todo sali add load balncer
}
