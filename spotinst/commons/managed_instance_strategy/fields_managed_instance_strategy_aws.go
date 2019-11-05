package managed_instance_strategy

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/managedinstance/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
	"log"
)

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Setup
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[LifeCycle] = commons.NewGenericField(
		commons.ManagedInstanceAwsStrategy,
		LifeCycle,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value *string = nil
			if managedInstance.Strategy != nil && managedInstance.Strategy.LifeCycle != nil {
				value = managedInstance.Strategy.LifeCycle
			}
			if err := resourceData.Set(string(LifeCycle), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(LifeCycle), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.Get(string(LifeCycle)).(string); ok && v != "" {
				managedInstance.Strategy.SetLifeCycle(spotinst.String(resourceData.Get(string(LifeCycle)).(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.Get(string(LifeCycle)).(string); ok && v != "" {
				managedInstance.Strategy.SetLifeCycle(spotinst.String(resourceData.Get(string(LifeCycle)).(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[Orientation] = commons.NewGenericField(
		commons.ManagedInstanceAwsStrategy,
		Orientation,
		&schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value *string = nil
			if managedInstance.Strategy != nil && managedInstance.Strategy.Orientation != nil {
				value = managedInstance.Strategy.Orientation
			}
			if err := resourceData.Set(string(Orientation), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(Orientation), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.Get(string(Orientation)).(string); ok && v != "" {
				managedInstance.Strategy.SetOrientation(spotinst.String(resourceData.Get(string(Orientation)).(string)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.Get(string(Orientation)).(string); ok && v != "" {
				managedInstance.Strategy.SetOrientation(spotinst.String(resourceData.Get(string(Orientation)).(string)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[DrainingTimeout] = commons.NewGenericField(
		commons.ManagedInstanceAwsStrategy,
		DrainingTimeout,
		&schema.Schema{
			Type:     schema.TypeInt,
			Optional: true,
			//Computed: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value *int = nil
			if managedInstance.Strategy != nil && managedInstance.Strategy.DrainingTimeout != nil {
				value = managedInstance.Strategy.DrainingTimeout
			}
			if err := resourceData.Set(string(DrainingTimeout), spotinst.IntValue(value)); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(DrainingTimeout), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			log.Printf("### tamir in draining_timeout create")
			if v, ok := resourceData.GetOkExists(string(DrainingTimeout)); ok {
				log.Printf("### tamir 1 %v", v.(int))
				log.Printf("### tamir 2 %v", spotinst.Int(v.(int)))
				managedInstance.Strategy.SetDrainingTimeout(spotinst.Int(v.(int)))
			} else {
				v1, o1 := resourceData.GetOkExists(string(DrainingTimeout))
				log.Printf("### tamir 3: %v, %v", v1, o1)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()

			if v, ok := resourceData.GetOkExists(string(DrainingTimeout)); ok {
				managedInstance.Strategy.SetDrainingTimeout(spotinst.Int(v.(int)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[FallbackToOd] = commons.NewGenericField(
		commons.ManagedInstanceAwsStrategy,
		FallbackToOd,
		&schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value *bool = nil
			if managedInstance.Strategy != nil && managedInstance.Strategy.FallbackToOnDemand != nil {
				value = managedInstance.Strategy.FallbackToOnDemand
			}
			if err := resourceData.Set(string(FallbackToOd), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(FallbackToOd), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.GetOk(string(FallbackToOd)); ok {
				managedInstance.Strategy.SetFallbackToOnDemand(spotinst.Bool(v.(bool)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.GetOk(string(FallbackToOd)); ok {
				managedInstance.Strategy.SetFallbackToOnDemand(spotinst.Bool(v.(bool)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[UtilizeReservedInstances] = commons.NewGenericField(
		commons.ManagedInstanceAwsStrategy,
		UtilizeReservedInstances,
		&schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value *bool = nil
			if managedInstance.Strategy != nil && managedInstance.Strategy.UtilizeReservedInstances != nil {
				value = managedInstance.Strategy.UtilizeReservedInstances
			}
			if err := resourceData.Set(string(UtilizeReservedInstances), value); err != nil {
				return fmt.Errorf(string(commons.FailureFieldReadPattern), string(UtilizeReservedInstances), err)
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.GetOk(string(UtilizeReservedInstances)); ok {
				managedInstance.Strategy.SetUtilizeReservedInstances(spotinst.Bool(v.(bool)))
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.GetOk(string(UtilizeReservedInstances)); ok {
				managedInstance.Strategy.SetUtilizeReservedInstances(spotinst.Bool(v.(bool)))
			}
			return nil
		},
		nil,
	)

	fieldsMap[RevertToSpot] = commons.NewGenericField(
		commons.ElastigroupAWS,
		RevertToSpot,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(PerformAt): {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if managedInstance.Strategy != nil && managedInstance.Strategy.RevertToSpot != nil {
				rts := managedInstance.Strategy.RevertToSpot
				result := make(map[string]interface{})
				result[string(PerformAt)] = spotinst.StringValue(rts.PerformAt)
				revertToSpot := []interface{}{result}
				if err := resourceData.Set(string(RevertToSpot), revertToSpot); err != nil {
					return fmt.Errorf("failed to set revertToSpot configuration: %#v", err)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.GetOk(string(RevertToSpot)); ok {
				if revertToSpot, err := expandAWSGroupRevertToSpot(v); err != nil {
					return err
				} else {
					managedInstance.Strategy.SetRevertToSpot(revertToSpot)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var revertToSpot *aws.RevertToSpot = nil
			if v, ok := resourceData.GetOk(string(RevertToSpot)); ok {
				if rts, err := expandAWSGroupRevertToSpot(v); err != nil {
					return err
				} else {
					revertToSpot = rts
				}
			}
			managedInstance.Strategy.SetRevertToSpot(revertToSpot)
			return nil
		},
		nil,
	)
	//todo sali need to add optimizationWindows (array) and revertToSpot with perrformat!
}

func expandAWSGroupRevertToSpot(data interface{}) (*aws.RevertToSpot, error) {
	revertToSpot := &aws.RevertToSpot{}
	list := data.([]interface{})
	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})

		var performAt *string = nil
		if v, ok := m[string(PerformAt)].(string); ok {
			performAt = spotinst.String(v)
		}
		revertToSpot.SetPerformAt(performAt)
	}
	//log.Printf("[DEBUG] Group revert to spot configuration: %s", stringutil.Stringify(revertToSpot))
	return revertToSpot, nil
}
