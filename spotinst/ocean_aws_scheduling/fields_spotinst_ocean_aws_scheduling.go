package ocean_aws_scheduling

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) {

	fieldsMap[ScheduledTask] = commons.NewGenericField(
		commons.OceanAWSScheduling,
		ScheduledTask,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(IsEnabled): {
						Type:     schema.TypeBool,
						Required: true,
					},

					string(TaskType): {
						Type:     schema.TypeString,
						Required: true,
					},

					string(CronExpression): {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.AWSClusterWrapper)
			elastigroup := egWrapper.GetCluster()
			if v, ok := resourceData.GetOk(string(ScheduledTask)); ok {
				if tasks, err := expandScheduledTasks(v); err != nil {
					return err
				} else {
					elastigroup.Scheduling.SetTasks(tasks)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.AWSClusterWrapper)
			elastigroup := egWrapper.GetCluster()
			var value []*aws.Tasks = nil
			if v, ok := resourceData.GetOk(string(ScheduledTask)); ok {
				if interfaces, err := expandScheduledTasks(v); err != nil {
					return err
				} else {
					value = interfaces
				}
			}
			elastigroup.Scheduling.SetTasks(value)
			return nil
		},
		nil,
	)

	fieldsMap[ShutdownHours] = commons.NewGenericField(
		commons.OceanAWSScheduling,
		ShutdownHours,
		&schema.Schema{
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(tasksIsEnabled): {
						Type:     schema.TypeBool,
						Optional: true,
					},

					string(TimeWindows): {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.AWSClusterWrapper)
			elastigroup := egWrapper.GetCluster()
			if v, ok := resourceData.GetOk(string(ScheduledTask)); ok {
				if tasks, err := expandScheduledTasks(v); err != nil {
					return err
				} else {
					elastigroup.Scheduling.SetTasks(tasks)
				}
			}
			return nil
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			egWrapper := resourceObject.(*commons.AWSClusterWrapper)
			elastigroup := egWrapper.GetCluster()
			var value []*aws.Tasks = nil
			if v, ok := resourceData.GetOk(string(ScheduledTask)); ok {
				if interfaces, err := expandScheduledTasks(v); err != nil {
					return err
				} else {
					value = interfaces
				}
			}
			elastigroup.Scheduling.SetTasks(value)
			return nil
		},
		nil,
	)

}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func flattenAzureGroupScheduledTasks(tasks []*aws.Tasks) []interface{} {
	result := make([]interface{}, 0, len(tasks))
	for _, t := range tasks {
		m := make(map[string]interface{})
		m[string(IsEnabled)] = spotinst.BoolValue(t.IsEnabled)
		m[string(TaskType)] = spotinst.StringValue(t.TaskType)
		m[string(CronExpression)] = spotinst.StringValue(t.CronExpression)
		result = append(result, m)
	}
	return result
}

func expandScheduledTasks(data interface{}) ([]*aws.Tasks, error) {
	list := data.(*schema.Set).List()
	tasks := make([]*aws.Tasks, 0, len(list))
	for _, item := range list {
		m := item.(map[string]interface{})
		task := &aws.Tasks{}

		if v, ok := m[string(IsEnabled)].(bool); ok {
			task.SetIsEnabled(spotinst.Bool(v))
		}

		if v, ok := m[string(TaskType)].(string); ok && v != "" {
			task.SetTaskType(spotinst.String(v))
		}

		if v, ok := m[string(CronExpression)].(string); ok && v != "" {
			task.SetCronExpression(spotinst.String(v))
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Setup
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//func Setup(fieldsMap map[commons.FieldName]*commons.GenericField) { //
//
//	fieldsMap[ScheduledTask] = commons.NewGenericField(
//		commons.OceanAWSScheduling,
//		ScheduledTask,
//		&schema.Schema{
//			Type:     schema.TypeList,
//			Optional: true,
//			MaxItems: 1,
//			Elem: &schema.Resource{
//				Schema: map[string]*schema.Schema{
//					string(ShutdownHours): {
//						Type:     schema.TypeList,
//						Optional: true,
//						MaxItems: 1,
//						Elem: &schema.Resource{
//							Schema: map[string]*schema.Schema{
//								string(IsEnabled): {
//									Type:     schema.TypeBool,
//									Optional: true,
//								},
//								string(TimeWindows): {
//									Type:     schema.TypeList,
//									Required: true,
//									Elem:     &schema.Schema{Type: schema.TypeString},
//								},
//							},
//						},
//					},
//					string(tasks): {
//						Type:     schema.TypeList,
//						Optional: true,
//						Elem: &schema.Resource{
//							Schema: map[string]*schema.Schema{
//								string(tasksIsEnabled): {
//									Type:     schema.TypeBool,
//									Optional: true,
//								},
//								string(CronExpression): {
//									Type:     schema.TypeString,
//									Required: true,
//								},
//								string(TaskType): {
//									Type:     schema.TypeString,
//									Required: true,
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//
//		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
//			return nil
//		},
//		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
//			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
//			cluster := clusterWrapper.GetCluster()
//			if v, ok := resourceData.GetOk(string(ScheduledTask)); ok {
//				if scheduling, err := expandScheduledTask(v); err != nil {
//					return err
//				} else {
//					cluster.SetScheduling(scheduling)
//				}
//			}
//			return nil
//		},
//		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
//			clusterWrapper := resourceObject.(*commons.AWSClusterWrapper)
//			cluster := clusterWrapper.GetCluster()
//			var value *aws.Scheduling = nil
//
//			if v, ok := resourceData.GetOk(string(ScheduledTask)); ok {
//				if scheduling, err := expandScheduledTask(v); err != nil {
//					return err
//				} else {
//					value = scheduling
//				}
//			}
//			cluster.SetScheduling(value)
//			return nil
//		},
//
//		nil,
//	)
//}
//
////-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
////            Utils
////-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//
//func expandScheduledTask(data interface{}) (*aws.Scheduling, error) {
//	scheduling := &aws.Scheduling{}
//	list := data.([]interface{})
//	if list != nil && list[0] != nil {
//		m := list[0].(map[string]interface{})
//
//		var shutdownHours *aws.ShutdownHours = nil
//		if v, ok := m[string(ShutdownHours)].(string); ok && v != "" {
//			runner, err := expandAShutdownHours(v)
//			if err != nil {
//				return nil, err
//			}
//			if runner != nil {
//				shutdownHours = runner
//			}
//		}
//		scheduling.SetShutdownHours(shutdownHours)
//
//	}
//
//	return scheduling, nil
//}
//
//func expandAShutdownHours(data interface{}) (*aws.ShutdownHours, error) {
//	if list := data.([]interface{}); len(list) > 0 && list[0] != nil {
//		shutdownHours := &aws.ShutdownHours{}
//		m := list[0].(map[string]interface{})
//
//		if v, ok := m[string(IsEnabled)].(bool); ok {
//			shutdownHours.SetIsEnabled(spotinst.Bool(v))
//		}
//
//		return shutdownHours, nil
//	}
//
//	return nil, nil
//}
