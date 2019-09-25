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
	MangedInstanceAwsResourceName ResourceName = "spotinst_managed_instance_aws"
)

var MangedInstanceResource *MangedInstanceTerraformResource

type MangedInstanceTerraformResource struct {
	GenericResource // embedding

}

type MangedInstanceAWSBeanstalkWrapper struct {
	elastigroup *aws.Group

	// Load balancer states
	StatusElbUpdated bool
	StatusTgUpdated  bool
	StatusMlbUpdated bool

	// Block devices states
	StatusEphemeralBlockDeviceUpdated bool
	StatusEbsBlockDeviceUpdated       bool
}

func NewMangedInstanceResource(fieldsMap map[FieldName]*GenericField) *MangedInstanceTerraformResource {
	return &MangedInstanceTerraformResource{
		GenericResource: GenericResource{
			resourceName: MangedInstanceAwsResourceName,
			fields:       NewGenericFields(fieldsMap),
		},
	}
}

func (res *MangedInstanceTerraformResource) OnRead(
	elastigroup *aws.Group, /// need to change it to the managed instance on the sdk
	resourceData *schema.ResourceData,
	meta interface{}) error {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return fmt.Errorf("resource fields are nil or empty, cannot read")
	}

	egWrapper := NewMangedInstanceWrapper()
	egWrapper.SetMangedInstance(elastigroup)

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
func (res *MangedInstanceTerraformResource) OnCreate(
	resourceData *schema.ResourceData,
	meta interface{}) (*aws.Group, error) {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return nil, fmt.Errorf("resource fields are nil or empty, cannot create")
	}

	egWrapper := NewMangedInstanceWrapper()

	for _, field := range res.fields.fieldsMap {
		if field.onCreate == nil {
			continue
		}
		log.Printf(string(ResourceFieldOnCreate), field.resourceAffinity, field.fieldNameStr)
		if err := field.onCreate(egWrapper, resourceData, meta); err != nil {
			return nil, err
		}
	}
	return egWrapper.GetMangedInstance(), nil
}

func (res *MangedInstanceTerraformResource) OnUpdate(
	resourceData *schema.ResourceData,
	meta interface{}) (bool, *aws.Group, error) {

	if res.fields == nil || res.fields.fieldsMap == nil || len(res.fields.fieldsMap) == 0 {
		return false, nil, fmt.Errorf("resource fields are nil or empty, cannot update")
	}

	egWrapper := NewMangedInstanceWrapper()
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

	return hasChanged, egWrapper.GetMangedInstance(), nil
}

func NewMangedInstanceWrapper() *MangedInstanceAWSBeanstalkWrapper { ////need to look into it
	return &MangedInstanceAWSBeanstalkWrapper{
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

func (egWrapper *MangedInstanceAWSBeanstalkWrapper) GetMangedInstance() *aws.Group {
	return egWrapper.elastigroup
}

func (egWrapper *MangedInstanceAWSBeanstalkWrapper) SetMangedInstance(elastigroup *aws.Group) { ////i think i need to change it change it as the sdk
	egWrapper.elastigroup = elastigroup
}
