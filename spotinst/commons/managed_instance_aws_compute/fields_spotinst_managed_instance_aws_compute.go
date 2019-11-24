package managed_instance_aws_compute

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[SubnetIds] = commons.NewGenericField(
		commons.ManagedInstanceAwsCompute,
		SubnetIds,
		&schema.Schema{
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value []string = nil
			if managedInstance.Compute != nil && managedInstance.Compute.SubnetIDs != nil {
				value = managedInstance.Compute.SubnetIDs
			}
			if err := resourceData.Set(string(SubnetIds), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(SubnetIds), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if value, ok := resourceData.GetOk(string(SubnetIds)); ok && value != nil {
				if subnetIds, err := expandSubnetIDs(value); err != nil {
					return err
				} else {
					managedInstance.Compute.SetSubnetIDs(subnetIds)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if value, ok := resourceData.GetOk(string(SubnetIds)); ok && value != nil {
				if subnetIds, err := expandSubnetIDs(value); err != nil {
					return err
				} else {
					managedInstance.Compute.SetSubnetIDs(subnetIds)
				}
			}
			return nil
		},
		nil,
	)

	fieldsMap[VpcId] = commons.NewGenericField(
		commons.ManagedInstanceAwsCompute,
		VpcId,
		&schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value *string = nil
			if managedInstance.Compute != nil && managedInstance.Compute.VpcId != nil {
				value = managedInstance.Compute.VpcId
			}
			if err := resourceData.Set(string(VpcId), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(VpcId), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			managedInstance.Compute.SetVpcId(spotinst.String(resourceData.Get(string(VpcId)).(string)))
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			managedInstance.Compute.SetVpcId(spotinst.String(resourceData.Get(string(VpcId)).(string)))
			return nil
		},
		nil,
	)

	fieldsMap[ElasticIp] = commons.NewGenericField(
		commons.ManagedInstanceAwsCompute,
		ElasticIp,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value *string = nil
			if managedInstance.Compute != nil && managedInstance.Compute.ElasticIP != nil {
				value = managedInstance.Compute.ElasticIP
			}
			if err := resourceData.Set(string(ElasticIp), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(ElasticIp), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if value, ok := resourceData.GetOk(string(ElasticIp)); ok && value != nil {
				managedInstance.Compute.SetElasticIP(spotinst.String(resourceData.Get(string(ElasticIp)).(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if value, ok := resourceData.GetOk(string(ElasticIp)); ok && value != nil {
				managedInstance.Compute.SetElasticIP(spotinst.String(resourceData.Get(string(ElasticIp)).(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[PrivateIp] = commons.NewGenericField(
		commons.ManagedInstanceAwsCompute,
		PrivateIp,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value *string = nil
			if managedInstance.Compute != nil && managedInstance.Compute.PrivateIP != nil {
				value = managedInstance.Compute.PrivateIP
			}
			if err := resourceData.Set(string(PrivateIp), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(PrivateIp), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if value, ok := resourceData.GetOk(string(PrivateIp)); ok && value != nil {
				managedInstance.Compute.SetPrivateIP(spotinst.String(resourceData.Get(string(PrivateIp)).(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if value, ok := resourceData.GetOk(string(PrivateIp)); ok && value != nil {
				managedInstance.Compute.SetPrivateIP(spotinst.String(resourceData.Get(string(PrivateIp)).(string)))
			}
			return nil
		},
		nil,
	)
}

func expandSubnetIDs(data interface{}) ([]string, error) {
	list := data.([]interface{})
	result := make([]string, 0, len(list))

	for _, v := range list {
		if subnetID, ok := v.(string); ok && subnetID != "" {
			result = append(result, subnetID)
		}
	}

	return result, nil
}
