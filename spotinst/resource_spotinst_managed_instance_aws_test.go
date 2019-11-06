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
	"testing"
)

func init() {
	resource.AddTestSweepers("spotinst_managed_instance_aws", &resource.Sweeper{
		Name: "spotinst_managed_instance_aws",
		F:    testSweepManagedInstance,
	})
}

func testSweepManagedInstance(region string) error {
	client, err := getProviderClient("aws")
	if err != nil {
		return fmt.Errorf("error getting client: %v", err)
	}

	conn := client.(*Client).managedInstance.CloudProviderAWS()

	input := &aws.ListManagedInstancesInput{}
	if resp, err := conn.List(context.Background(), input); err != nil {
		return fmt.Errorf("error getting list of groups to sweep")
	} else {
		if len(resp.ManagedInstances) == 0 {
			log.Printf("[INFO] No groups to sweep")
		}
		for _, managedInstance := range resp.ManagedInstances {
			if strings.Contains(spotinst.StringValue(managedInstance.Name), "test-acc-") {
				if _, err := conn.Delete(context.Background(), &aws.DeleteManagedInstanceInput{ManagedInstanceID: managedInstance.ID}); err != nil {
					return fmt.Errorf("unable to delete managedInstance %v in sweep", spotinst.StringValue(managedInstance.ID))
				} else {
					log.Printf("Sweeper deleted %v\n", spotinst.StringValue(managedInstance.ID))
				}
			}
		}
	}
	return nil
}

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
	provider string
	name     string
	region   string

	variables            string
	fieldsToAppend       string
	updateBaselineFields bool
}

func createManagedInstanceTerraform(ccm *ManagedInstanceConfigMetadata) string {
	if ccm == nil {
		return ""
	}

	if ccm.provider == "" {
		ccm.provider = "aws"
	}

	template :=
		`provider "aws" {
	 token   = "fake"
	 account = "fake"
	}
	`

	if ccm.updateBaselineFields {
		format := testBaselineManagedInstanceConfig_Update
		template += fmt.Sprintf(format,
			ccm.name,
			ccm.provider,
			ccm.name,
			ccm.fieldsToAppend,
		)
	} else {
		format := testBaselineManagedInstanceConfig_Create
		template += fmt.Sprintf(format,
			ccm.name,
			ccm.provider,
			ccm.name,
			ccm.fieldsToAppend,
		)
	}

	if ccm.variables != "" {
		template = ccm.variables + "\n" + template
	}

	log.Printf("Terraform [%v] template:\n%v", ccm.name, template)
	return template
}

// region managedInstance: Baseline
func TestAccSpotinstManagedInstance_Baseline(t *testing.T) {
	name := "test-acc-cluster-managed-instance"
	resourceName := createManagedInstanceAWSResourceName(name)

	var cluster aws.ManagedInstance
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testManagedInstanceAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createManagedInstanceTerraform(&ManagedInstanceConfigMetadata{
					name: name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagedInstanceAWSExists(&cluster, resourceName),
					testCheckManagedInstanceAWSAttributes(&cluster, name),
					resource.TestCheckResourceAttr(resourceName, "persist_private_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "persist_block_devices", "true"),
					resource.TestCheckResourceAttr(resourceName, "persist_root_device", "true"),
					resource.TestCheckResourceAttr(resourceName, "block_devices_mode", "reattach"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", "vpc-9dee6bfa"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"), //need to add
					resource.TestCheckResourceAttr(resourceName, "types.#", "1"),      //need to add
					resource.TestCheckResourceAttr(resourceName, "types.0", "t1.micro"),
					resource.TestCheckResourceAttr(resourceName, "image_id", "ami-082b5a644766e0e6f"),
					resource.TestCheckResourceAttr(resourceName, "product", "Linux/UNIX"),

					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.0", "subnet-7f3fbf06"),

					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", "sg-1a29b065"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.1", "sg-5750fb2f"),
				),
			},
			{
				Config: createManagedInstanceTerraform(&ManagedInstanceConfigMetadata{
					//clusterName:          clusterName,
					name:                 name,
					updateBaselineFields: true}),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagedInstanceAWSExists(&cluster, resourceName),
					testCheckManagedInstanceAWSAttributes(&cluster, name),
					resource.TestCheckResourceAttr(resourceName, "persist_private_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "persist_block_devices", "true"),
					resource.TestCheckResourceAttr(resourceName, "persist_root_device", "true"),
					resource.TestCheckResourceAttr(resourceName, "block_devices_mode", "reattach"),
					resource.TestCheckResourceAttr(resourceName, "vpc_id", "vpc-9dee6bfa"),
					resource.TestCheckResourceAttr(resourceName, "types.#", "2"), //need to add
					resource.TestCheckResourceAttr(resourceName, "types.0", "t1.micro"),
					resource.TestCheckResourceAttr(resourceName, "types.1", "t3.medium"),
					resource.TestCheckResourceAttr(resourceName, "image_id", "ami-082b5a644766e0e6f"),
					resource.TestCheckResourceAttr(resourceName, "product", "Linux/UNIX"),

					resource.TestCheckResourceAttr(resourceName, "subnet_ids.#", "1"), //need to add
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.0", "subnet-7f3fbf06"),

					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", "sg-1a29b065"),
				),
			},
		},
	})
}

const testBaselineManagedInstanceConfig_Create = `
resource "` + string(commons.ManagedInstanceAwsResourceName) + `" "%v" {
  provider = "%v"  

  name = "%v"
  region = "us-west-2"
  product = "Linux/UNIX"


  persist_private_ip = "false"
  persist_block_devices = "true"
  persist_root_device = "true"
  block_devices_mode = "reattach"

  subnet_ids = ["subnet-7f3fbf06"]  
  vpc_id = "vpc-9dee6bfa"

  types = ["t1.micro"]
  preferred_type = "t1.micro"

  image_id = "ami-082b5a644766e0e6f"

  security_group_ids = ["sg-1a29b065","sg-5750fb2f"]


 %v
}
`

const testBaselineManagedInstanceConfig_Update = `
resource "` + string(commons.ManagedInstanceAwsResourceName) + `" "%v" {
  provider = "%v"

  name = "%v"
  region = "us-west-2"
  product = "Linux/UNIX"
  //strategy

  persist_private_ip = "false"
  persist_block_devices = "true"
  persist_root_device = "true"
  block_devices_mode = "reattach"

  subnet_ids = ["subnet-7f3fbf06"]   //need to add
  vpc_id = "vpc-9dee6bfa"

  types = [
    "t1.micro",
    "t3.medium",]
  preferred_type = "t1.micro"

  image_id = "ami-082b5a644766e0e6f"

  security_group_ids = ["sg-1a29b065"]

 %v
}
`

// endregion

// region oceanECS: Strategy
func TestAccSpotinstManagedInstance_All(t *testing.T) {
	name := "test-acc-cluster-managed-instance-all"
	resourceName := createManagedInstanceAWSResourceName(name)

	var cluster aws.ManagedInstance
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "aws") },
		Providers:    TestAccProviders,
		CheckDestroy: testManagedInstanceAWSDestroy,

		Steps: []resource.TestStep{
			{
				Config: createManagedInstanceTerraform(&ManagedInstanceConfigMetadata{
					name:           name,
					fieldsToAppend: managedInstanceAll_Create,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagedInstanceAWSExists(&cluster, resourceName),
					testCheckManagedInstanceAWSAttributes(&cluster, name),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttr(resourceName, "draining_timeout", "120"),
					resource.TestCheckResourceAttr(resourceName, "life_cycle", "on_demand"),
					resource.TestCheckResourceAttr(resourceName, "orientation", "balanced"),
					resource.TestCheckResourceAttr(resourceName, "fall_back_to_od", "false"),
					resource.TestCheckResourceAttr(resourceName, "utilize_reserved_instances", "false"),
					//resource.TestCheckResourceAttr(resourceName, "perform_at", "never"),
					resource.TestCheckResourceAttr(resourceName, "health_check_type", "EC2"),
					//resource.TestCheckResourceAttr(resourceName, "auto_healing", "true"),
					resource.TestCheckResourceAttr(resourceName, "grace_period", "180"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_duration", "60"),
					resource.TestCheckResourceAttr(resourceName, "ebs_optimized", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_monitoring", "false"),
					resource.TestCheckResourceAttr(resourceName, "placement_tenancy", "default"),
					resource.TestCheckResourceAttr(resourceName, "key_pair", "TamirKeyPair"),
					//resource.TestCheckResourceAttr(resourceName, "user_data", "false"),
				),
			},
			{
				ResourceName: resourceName,
				Config: createManagedInstanceTerraform(&ManagedInstanceConfigMetadata{
					name:           name,
					fieldsToAppend: managedInstanceAll_Update,
				}),
				Check: resource.ComposeTestCheckFunc(
					testCheckManagedInstanceAWSExists(&cluster, resourceName),
					testCheckManagedInstanceAWSAttributes(&cluster, name),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
					resource.TestCheckResourceAttr(resourceName, "draining_timeout", "240"),
					resource.TestCheckResourceAttr(resourceName, "life_cycle", "spot"),
					resource.TestCheckResourceAttr(resourceName, "orientation", "cheapest"),
					resource.TestCheckResourceAttr(resourceName, "fall_back_to_od", "true"),
					resource.TestCheckResourceAttr(resourceName, "utilize_reserved_instances", "true"),
					//resource.TestCheckResourceAttr(resourceName, "perform_at", "always"),
					resource.TestCheckResourceAttr(resourceName, "health_check_type", "MULTAI_TARGET_SET"),
					//resource.TestCheckResourceAttr(resourceName, "auto_healing", "false"),
					resource.TestCheckResourceAttr(resourceName, "grace_period", "100"),
					resource.TestCheckResourceAttr(resourceName, "unhealthy_duration", "120"),
					resource.TestCheckResourceAttr(resourceName, "ebs_optimized", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_monitoring", "true"),
					resource.TestCheckResourceAttr(resourceName, "placement_tenancy", "dedicated"),
					resource.TestCheckResourceAttr(resourceName, "key_pair", "TamirKeyPair"), //need to change the key pair
					//resource.TestCheckResourceAttr(resourceName, "user_data", "true"),
				),
			},
		},
	})
}

const managedInstanceAll_Create = `
  description = "description"
 //strategy
  life_cycle = "on_demand"
  orientation = "balanced"
  draining_timeout = 120
  fall_back_to_od = "false"
  utilize_reserved_instances = "false"
 //  optimization_windows = ""   // need to test
 //revert_to_spot = {   need to test here
 //perform_at = "never"
 //}
 //healthCheck
 health_check_type = "EC2"
 //auto_healing = "true"
 grace_period = "180"
 unhealthy_duration = "60"

 //compute
//  elastic_ip = "1.1.1.1"
//  private_ip = "1.1.1.1"

 //  launchSpecification
 ebs_optimized = "false"
 enable_monitoring = "false"
 placement_tenancy = "default"
 //  iam_instance_profile = "arn:aws:iam::842422002533:instance-profile/ecsInstanceRole"
 key_pair = "TamirKeyPair"
 //tags = [
 //  {
 //    key   = "Env"
 //    value = "production"
 //  },
 //  {
 //    key   = "Name"
 //    value = "default-production"
 //  }
 //]
 user_data  = "echo hello world"
 //shutdown_script = "echo bye world"
//  creditSpecification
//    cpu_credits = "standard"

 //network_interface = {   ///cant do it the oder way
 //  device_index = 0
 //  associate_public_ip_address = "true"
 //  associate_ipv6_address = "false"
 //  }

 ////  scheduling
 //scheduled_task = {
 //  task_type             = "pause"
 //  cron_expression       = "" // need to test
 //  //start_time            = "2020-01-01T01:00:00Z"  ///cant read it?
 //  frequency             = "hourly"
 //  is_enabled            = "true"
 //}


 //  integrations
 //loadBalancersConfig
 //route53  //need to test

// --------------------------------
`

const managedInstanceAll_Update = `
// --------------------
  description = "description2"
 //strategy
  life_cycle = "spot"
  orientation = "cheapest"
  draining_timeout = 240
  fall_back_to_od = "true"
  utilize_reserved_instances = "true"
 //  optimization_windows = ""   // need to test
 //revert_to_spot = { need to test here
 //perform_at = "always"
 //}
 //healthCheck
 health_check_type = "MULTAI_TARGET_SET"
 //auto_healing = "false"
 grace_period = "100"
 unhealthy_duration = "120"

 ebs_optimized = "true"
 enable_monitoring = "true"
 placement_tenancy = "dedicated"
 key_pair = "TamirKeyPair"  //need to change the key pair 

// --------------------------------
`

// endregion
