package spotinst

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spotinst/spotinst-sdk-go/service/managedinstance/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
	"log"
	"strings"
)

//func init() {
//	resource.AddTestSweepers("spotinst_managed_instance", &resource.Sweeper{
//		Name: "spotinst_managed_instance",
//		F:    testSweepManagedInstance,
//	})
//}

//func testSweepManagedInstance(region string) error {
//	client, err := getProviderClient("aws")
//	if err != nil {
//		return fmt.Errorf("error getting client: %v", err)
//	}
//
//	conn := client.(*Client).managedInstance.CloudProviderAWS()
//
//	input := &aws.ListMangedInstancesInput{}
//	if resp, err := conn.List(context.Background(), input); err != nil {
//		return fmt.Errorf("error getting list of groups to sweep")
//	} else {
//		if len(resp.MangedInstances) == 0 {
//			log.Printf("[INFO] No groups to sweep")
//		}
//		for _, managedInstance := range resp.MangedInstances {
//			if strings.Contains(spotinst.StringValue(managedInstance.Name), "test-acc-") {
//				if _, err := conn.Delete(context.Background(), &aws.DeleteManagedInstanceInput{ManagedInstanceID: managedInstance.ID}); err != nil {
//					return fmt.Errorf("unable to delete managedInstance %v in sweep", spotinst.StringValue(managedInstance.ID))
//				} else {
//					log.Printf("Sweeper deleted %v\n", spotinst.StringValue(managedInstance.ID))
//				}
//			}
//		}
//	}
//	return nil
//}

func createManagedInstanceAWSResourceName(name string) string {
	return fmt.Sprintf("%v.%v", string(commons.ManagedInstanceAwsResourceName), name)
}

func testManagedInstanceAWSDestroy(s *terraform.State) error {
	client := testAccProviderAWS.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != string(commons.ManagedInstanceAwsResourceName) {
			continue
		}
		input := &aws.ReadManagedInstanceInput{ManagedInstanceID: spotinst.String(rs.Primary.ID)}
		resp, err := client.managedInstance.CloudProviderAWS().Read(context.Background(), input)
		if err == nil && resp != nil && resp.ManagedInstance != nil {
			return fmt.Errorf("managedInstance still exists")
		}
	}
	return nil
}

func testCheckManagedInstanceAWSAttributes(managedInstance *aws.ManagedInstance, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if spotinst.StringValue(managedInstance.Name) != expectedName {
			return fmt.Errorf("bad content: %v", managedInstance.Name)
		}
		return nil
	}
}

func testCheckManagedInstanceAWSExists(managedInstance *aws.ManagedInstance, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}
		client := testAccProviderAWS.Meta().(*Client)
		input := &aws.ReadManagedInstanceInput{ManagedInstanceID: spotinst.String(rs.Primary.ID)}
		resp, err := client.managedInstance.CloudProviderAWS().Read(context.Background(), input)
		if err != nil {
			return err
		}
		if spotinst.StringValue(resp.ManagedInstance.Name) != rs.Primary.Attributes["name"] {
			return fmt.Errorf("ManagedInstance not found: %+v,\n %+v\n", resp.ManagedInstance, rs.Primary.Attributes)
		}
		*managedInstance = *resp.ManagedInstance
		return nil
	}
}

type ManagedInstanceConfigMetadata struct {
	managedInstanceID    string //maybe need id
	provider             string
	fieldsToAppend       string
	updateBaselineFields bool
}

func createManagedInstanceAWSTerraform(gcm *ManagedInstanceConfigMetadata) string {
	if gcm == nil {
		return ""
	}

	if gcm.provider == "" {
		gcm.provider = "aws"
	}

	template :=
		`provider "aws" {
	token   = "fake"
	account = "fake"
	}
	`
	if gcm.updateBaselineFields {
		format := testBaselineMnagedInstanceAWSUpdate
		template += fmt.Sprintf(format,
			gcm.managedInstanceID,
			gcm.provider,
			gcm.managedInstanceID,
			gcm.fieldsToAppend,
		)
	} else {
		//	format :=
		//	template += fmt.Sprintf(format,
		//		gcm.managedInstanceID,
		//		gcm.provider,
		//		gcm.managedInstanceID,
		//	)
	}

	log.Printf("Terraform [%v] template:\n%v", gcm.managedInstanceID, template)
	return template
}

const testBaselineMnagedInstanceAWSUpdate = `
//resource "` + string(commons.ManagedInstanceAwsResourceName) + `" "%v" {
//  provider = "%v"
//
//  ocean_id = "%v"
//  image_id = "ami-79826301"
//  security_groups = ["sg-0041bd3fd6aa2ee3c", "sg-0195f2ac3a6014a15"]
//  user_data = "hello world updated"
//  iam_instance_profile = "updated"
//  
//  labels = {
//    key = "label key updated"
//    value = "label value updated"
//  }
//
//  taints = [{
//    key = "taint key updated"
//    value = "taint value updated"
//    effect = "NoExecute"
//  }]
//
//%v
//}
`

// endregion
