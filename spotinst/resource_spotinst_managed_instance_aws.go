package spotinst

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/managedinstance/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws_block_devices"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws_instance_types"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws_integrations"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws_launch_configuration"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws_network_interface"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws_scaling_policies"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws_scheduled_task"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws_stateful"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/elastigroup_aws_strategy"
	"log"
	"strings"
	"time"
)

func resourceSpotinstMangedInstanceAws() *schema.Resource {
	setupMangedInstanceResource()

	return &schema.Resource{
		Create: resourceSpotinstManagedInstanceAwsCreate,
		Read:   resourceSpotinstManagedInstanceAwsRead,
		//Update: resourceSpotinstManagedInstanceAwsUpdate,
		//Delete: resourceSpotinstManagedInstanceAwsDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: commons.ManagedInstanceResource.GetSchemaMap(),
	}
}

func setupMangedInstanceResource() {
	fieldsMap := make(map[commons.FieldName]*commons.GenericField)

	elastigroup_aws.Setup(fieldsMap)
	elastigroup_aws_block_devices.Setup(fieldsMap)
	elastigroup_aws_instance_types.Setup(fieldsMap)
	elastigroup_aws_integrations.Setup(fieldsMap)
	elastigroup_aws_launch_configuration.Setup(fieldsMap)
	elastigroup_aws_network_interface.Setup(fieldsMap)
	elastigroup_aws_scaling_policies.Setup(fieldsMap)
	elastigroup_aws_scheduled_task.Setup(fieldsMap)
	elastigroup_aws_stateful.Setup(fieldsMap)
	elastigroup_aws_strategy.Setup(fieldsMap)

	commons.ManagedInstanceResource = commons.NewManagedInstanceResource(fieldsMap)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Delete
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//func resourceSpotinstMangedInstanceAwsDelete(resourceData *schema.ResourceData, meta interface{}) error {
//	id := resourceData.Id()
//	log.Printf(string(commons.ResourceOnDelete),
//		commons.MangedInstanceResource.GetName(), id)
//
//	if err := deleteGroup(resourceData, meta); err != nil {
//		return err
//	}
//
//	log.Printf("===> Elastigroup deleted successfully: %s <===", resourceData.Id())
//	resourceData.SetId("")
//	return nil
//}
//
//func deleteGroup(resourceData *schema.ResourceData, meta interface{}) error {
//	groupId := resourceData.Id()
//	input := &aws.DeleteGroupInput{
//		GroupID: spotinst.String(groupId),
//	}
//
//	if statefulDeallocation, exists := resourceData.GetOkExists(string(elastigroup_aws_stateful.StatefulDeallocation)); exists {
//		list := statefulDeallocation.([]interface{})
//		if list != nil && len(list) > 0 && list[0] != nil {
//			m := list[0].(map[string]interface{})
//
//			var result = &aws.StatefulDeallocation{}
//			if shouldDeleteImage, ok := m[string(elastigroup_aws_stateful.ShouldDeleteImages)].(bool); ok && shouldDeleteImage {
//				result.ShouldDeleteImages = spotinst.Bool(shouldDeleteImage)
//			}
//
//			if shouldDeleteNetworkInterfaces, ok := m[string(elastigroup_aws_stateful.ShouldDeleteNetworkInterfaces)].(bool); ok && shouldDeleteNetworkInterfaces {
//				result.ShouldDeleteNetworkInterfaces = spotinst.Bool(shouldDeleteNetworkInterfaces)
//			}
//
//			if shouldDeleteSnapshots, ok := m[string(elastigroup_aws_stateful.ShouldDeleteSnapshots)].(bool); ok && shouldDeleteSnapshots {
//				result.ShouldDeleteSnapshots = spotinst.Bool(shouldDeleteSnapshots)
//			}
//
//			if shouldDeleteVolumes, ok := m[string(elastigroup_aws_stateful.ShouldDeleteVolumes)].(bool); ok && shouldDeleteVolumes {
//				result.ShouldDeleteVolumes = spotinst.Bool(shouldDeleteVolumes)
//			}
//
//			input.StatefulDeallocation = result
//		}
//	}
//
//	if json, err := commons.ToJson(input); err != nil {
//		return err
//	} else {
//		log.Printf("===> Group delete configuration: %s", json)
//	}
//
//	if _, err := meta.(*Client).mangedInstance.CloudProviderAWS().Delete(context.Background(), input); err != nil {
//		return fmt.Errorf("[ERROR] onDelete() -> Failed to delete group: %s", err)
//	}
//	return nil
//}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Read
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// ErrCodeGroupNotFound for service response error code "GROUP_DOESNT_EXIST".
//const ErrCodeGroupNotFound = "GROUP_DOESNT_EXIST"

func resourceSpotinstManagedInstanceAwsRead(resourceData *schema.ResourceData, meta interface{}) error {
	id := resourceData.Id()
	log.Printf(string(commons.ResourceOnRead),
		commons.ManagedInstanceResource.GetName(), id)

	input := &aws.ReadGroupInput{GroupID: spotinst.String(id)}
	resp, err := meta.(*Client).managedInstance.CloudProviderAWS().Read(context.Background(), input)
	if err != nil {
		// If the group was not found, return nil so that we can show
		// that the group does not exist
		if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
			for _, err := range errs {
				if err.Code == ErrCodeGroupNotFound {
					resourceData.SetId("")
					return nil
				}
			}
		}

		// Some other error, report it.
		return fmt.Errorf("failed to read group: %s", err)
	}

	// If nothing was found, then return no state.
	groupResponse := resp.Group
	if groupResponse == nil {
		resourceData.SetId("")
		return nil
	}

	if err := commons.ElastigroupResource.OnRead(groupResponse, resourceData, meta); err != nil {
		return err
	}
	log.Printf("===> Elastigroup read successfully: %s <===", id)
	return nil
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Create
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func resourceSpotinstManagedInstanceAwsCreate(resourceData *schema.ResourceData, meta interface{}) error {
	log.Printf(string(commons.ResourceOnCreate),
		commons.ManagedInstanceResource.GetName())

	mangedInstance, err := commons.ManagedInstanceResource.OnCreate(resourceData, meta)
	if err != nil {
		return err
	}

	groupId, err := createGroup(resourceData, managedInstance, meta.(*Client))
	if err != nil {
		return err
	}

	resourceData.SetId(spotinst.StringValue(groupId))

	if capacity, ok := resourceData.GetOkExists(string(elastigroup_aws.WaitForCapacity)); ok {
		if *mangedInstance.Capacity.Target < capacity.(int) {

			return fmt.Errorf("[ERROR] Your target healthy capacity must be less than or equal to your desired capcity")
		}
		if timeout, ok := resourceData.GetOkExists(string(elastigroup_aws.WaitForCapacityTimeout)); ok {
			err := awaitReady(groupId, timeout.(int), capacity.(int), meta.(*Client))
			if err != nil {
				return fmt.Errorf("[ERROR] Timed out when creating group: %s", err)
			}
		}
	}

	log.Printf("===> Elastigroup created successfully: %s <===", resourceData.Id())

	return resourceSpotinstElastigroupAwsRead(resourceData, meta)
}

func createGroup(resourceData *schema.ResourceData, group *aws.Group, spotinstClient *Client) (*string, error) {
	if json, err := commons.ToJson(group); err != nil {
		return nil, err
	} else {
		log.Printf("===> Group create configuration: %s", json)
	}

	if v, ok := resourceData.Get(string(elastigroup_aws_launch_configuration.IamInstanceProfile)).(string); ok && v != "" {
		time.Sleep(5 * time.Second)
	}
	input := &aws.CreateGroupInput{Group: group}

	var resp *aws.CreateGroupOutput = nil
	err := resource.Retry(time.Minute, func() *resource.RetryError {
		r, err := spotinstClient.elastigroup.CloudProviderAWS().Create(context.Background(), input)
		if err != nil {
			// Checks whether we should retry the group creation.
			if errs, ok := err.(client.Errors); ok && len(errs) > 0 {
				for _, err := range errs {
					if err.Code == "InvalidParameterValue" &&
						strings.Contains(err.Message, "Invalid IAM Instance Profile") {
						return resource.RetryableError(err)
					}
				}
			}

			// Some other error, report it.
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("[ERROR] failed to create group: %s", err)
	}
	return resp.Group.ID, nil
}
