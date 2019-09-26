package commons

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/elastigroup/providers/aws"
	"log"
)

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Variables
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
const (
	ManagedInstanceAwsResourceName ResourceName = "spotinst_managed_instance_aws"
)

var ManagedInstanceResource *ManagedInstanceTerraformResource

type ManagedInstanceTerraformResource struct {
	GenericResource // embedding

}

type MangedInstanceAWSWrapper struct {
	elastigroup *aws.Group

	// Load balancer states
	StatusElbUpdated bool
	StatusTgUpdated  bool
	StatusMlbUpdated bool

	// Block devices states
	StatusEphemeralBlockDeviceUpdated bool
	StatusEbsBlockDeviceUpdated       bool
}

func NewManagedInstanceResource(fieldsMap map[FieldName]*GenericField) *ManagedInstanceTerraformResource {
	return &ManagedInstanceTerraformResource{
		GenericResource: GenericResource{
			resourceName: ManagedInstanceAwsResourceName,
			fields:       NewGenericFields(fieldsMap),
		},
	}
}

func (res *ManagedInstanceTerraformResource) OnRead(
	elastigroup *aws.Group, /// need to change it to the managed instance on the sdk
	resourceData *schema.ResourceData,
	meta interface{}) error {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return fmt.Errorf("resource fields are nil or empty, cannot read")
	}

	egWrapper := NewManagedInstanceWrapper()
	egWrapper.SetManagedInstance(elastigroup)

	for _, field := range res.fields.fieldsMap {
		if field.onRead == nil {
			continue
		}
		log.Printf(string(ResourceFieldOnRead), field.resourceAffinity, field.fieldNameStr)
		if err := field.onRead(egWrapper, resourceData, meta); err != nil {
			return err
		}
	}
	return nil
}
func (res *ManagedInstanceTerraformResource) OnCreate(
	resourceData *schema.ResourceData,
	meta interface{}) (*aws.Group, error) {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return nil, fmt.Errorf("resource fields are nil or empty, cannot create")
	}

	egWrapper := NewManagedInstanceWrapper()

	for _, field := range res.fields.fieldsMap {
		if field.onCreate == nil {
			continue
		}
		log.Printf(string(ResourceFieldOnCreate), field.resourceAffinity, field.fieldNameStr)
		if err := field.onCreate(egWrapper, resourceData, meta); err != nil {
			return nil, err
		}
	}
	return egWrapper.GetManagedInstance(), nil
}

func (res *ManagedInstanceTerraformResource) OnUpdate(
	resourceData *schema.ResourceData,
	meta interface{}) (bool, *aws.Group, error) {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return false, nil, fmt.Errorf("resource fields are nil or empty, cannot update")
	}

	egWrapper := NewManagedInstanceWrapper()
	hasChanged := false
	for _, field := range res.fields.fieldsMap {
		if field.onUpdate == nil {
			continue
		}
		if field.hasFieldChange(resourceData, meta) {
			log.Printf(string(ResourceFieldOnUpdate), field.resourceAffinity, field.fieldNameStr)
			if err := field.onUpdate(egWrapper, resourceData, meta); err != nil {
				return false, nil, err
			}
			hasChanged = true
		}
	}

	return hasChanged, egWrapper.GetManagedInstance(), nil
}

func NewManagedInstanceWrapper() *MangedInstanceAWSWrapper { ////need to look into it
	return &MangedInstanceAWSWrapper{
		elastigroup: &aws.Group{
			Scaling:     &aws.Scaling{},
			Scheduling:  &aws.Scheduling{},
			Integration: &aws.Integration{},
			Compute: &aws.Compute{
				LaunchSpecification: &aws.LaunchSpecification{
					LoadBalancersConfig: &aws.LoadBalancersConfig{},
				},
				InstanceTypes: &aws.InstanceTypes{},
			},
			Capacity: &aws.Capacity{},
			Strategy: &aws.Strategy{
				Persistence: &aws.Persistence{},
			},
		},
	}
}

func (egWrapper *MangedInstanceAWSWrapper) GetManagedInstance() *aws.Group {
	return egWrapper.elastigroup
}

func (egWrapper *MangedInstanceAWSWrapper) SetManagedInstance(elastigroup *aws.Group) { ////i think i need to change it change it as the sdk
	egWrapper.elastigroup = elastigroup
}
