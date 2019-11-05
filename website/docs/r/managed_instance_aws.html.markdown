---
layout: "spotinst"
page_title: "Spotinst: managed_instance_aws"
sidebar_current: "docs-do-resource-managed_instance_aws"
description: |-
  Provides a Spotinst AWS group resource.
---

# spotinst\_elastigroup\_aws

Provides a Spotinst AWS managedInstance resource.

## Example Usage

```hcl
# Create an MangedInstance
resource "spotinst_managed_instance_aws" "default-managed-instance" {

  name        = "default-managed-instance"
  description = "created by Terraform"
  region      = "us-west-2"

    //strategy
  lifecycle      = "on_demand"
  orientation    = "balanced"
  draining_timeout = ""
  fallback_to_ondemand  = false
  utilize_reserved_instances = "true"
  optimizationWindows
  revertToSpot


  product     = "Linux/UNIX"
  subnet_ids = ["subnet-79da021e","subnet-7f3fbf06"] 
  image_id              = "ami-a27d8fda"
  iam_instance_profile  = "iam-profile"
  key_name              = "my-key.ssh"
  security_groups       = ["sg-123456"]
  user_data             = "echo hello world"
  enable_monitoring     = false
  ebs_optimized         = false
  placement_tenancy     = "default"

  instance_types_ondemand       = "m3.2xlarge"
  instance_types_spot           = ["m3.xlarge", "m3.2xlarge"]
  instance_types_preferred_spot = ["m3.xlarge"]

  
  cpu_credits           = "unlimited"

  wait_for_capacity         = 5
  wait_for_capacity_timeout = 300
  

  tags = [
  {
     key   = "Env"
     value = "production"
  }, 
  {
     key   = "Name"
     value = "default-production"
  },
  {
     key   = "Project"
     value = "app_v2"
  }
 ]

 
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The managedInstance name.
* `description` - (Optional) The managedInstance description.
* `product` - (Required) Operation system type. Valid values: `"Linux/UNIX"`, `"SUSE Linux"`, `"Windows"`. 
For EC2 Classic instances:  `"Linux/UNIX (Amazon VPC)"`, `"SUSE Linux (Amazon VPC)"`, `"Windows (Amazon VPC)"`.    

lifecycle
* `orientation` - (Optional) Select a prediction strategy. Valid values: `"balanced"`, `"costOriented"`, `"equalAzDistribution"`, `"availabilityOriented"`, `"cheapest"`.    
* `draining_timeout` - (Optional) The time in seconds, the instance is allowed to run while detached from the ELB. This is to allow the instance time to be drained from incoming TCP connections before terminating it, during a scale down operation.
* `fallback_to_ondemand` - (Required) In a case of no Spot instances available, Elastigroup will launch on-demand instances instead.
* `utilize_reserved_instances` - (Optional) In a case of any available reserved instances, Elastigroup will utilize them first before purchasing Spot instances.
optimizationWindows
revertToSpot



* `availability_zones` - (Optional) List of Strings of availability zones. When this parameter is set, `subnet_ids` should be left unused.
Note: `availability_zones` naming syntax follows the convention `availability-zone:subnet:placement-group-name`. For example, to set an AZ in `us-east-1` with subnet `subnet-123456` and placement group `ClusterI03`, you would set:
`availability_zones = ["us-east-1a:subnet-123456:ClusterI03"]`

* `subnet_ids` - (Optional) List of Strings of subnet identifiers.
Note: When this parameter is set, `availability_zones` should be left unused.

* `region` - (Optional) The AWS region your group will be created in.
Note: This parameter is required if you specify subnets (through subnet_ids). This parameter is optional if you specify Availability Zones (through availability_zones).

* `preferred_availability_zones` - The AZs to prioritize when launching Spot instances. If no markets are available in the Preferred AZs, Spot instances are launched in the non-preferred AZs. 
Note: Must be a sublist of `availability_zones` and `orientation` value must not be `"equalAzDistribution"`.

* `security_groups` - (Required) A list of associated security group IDS.
* `image_id` - (Optional) The ID of the AMI used to launch the instance.
* `iam_instance_profile` - (Optional) The ARN or name of an IAM instance profile to associate with launched instances.
* `key_name` - (Optional) The key name that should be used for the instance.
* `enable_monitoring` - (Optional) Indicates whether monitoring is enabled for the instance.
* `user_data` - (Optional) The user data to provide when launching the instance.
* `shutdown_script` - (Optional) The Base64-encoded shutdown script that executes prior to instance termination, for more information please see: [Shutdown Script](https://api.spotinst.com/integration-docs/elastigroup/concepts/compute-concepts/shutdown-scripts/)
* `ebs_optimized` - (Optional) Enable high bandwidth connectivity between instances and AWS’s Elastic Block Store (EBS). For instance types that are EBS-optimized by default this parameter will be ignored.
* `placement_tenancy` - (Optional) Enable dedicated tenancy. Note: There is a flat hourly fee for each region in which dedicated tenancy is used.

* `instance_types_ondemand` - (Required) The type of instance determines your instance's CPU capacity, memory and storage (e.g., m1.small, c1.xlarge).
* `instance_types_spot` - (Required) One or more instance types.
* `instance_types_preferred_spot` - (Optional) Prioritize a subset of spot instance types. Must be a subset of the selected spot instance types.
* `instance_types_weights` - (Optional) List of weights per instance type for weighted groups. Each object in the list should have the following attributes:
    * `weight` - (Required) Weight per instance type (Integer).
    * `instance_type` - (Required) Name of instance type (String).

* `cpu_credits` - (Optional) Controls how T3 instances are launched. Valid values: `standard`, `unlimited`.
* `wait_for_capacity` - (Optional) Minimum number of instances in a 'HEALTHY' status that is required before continuing. This is ignored when updating with blue/green deployment. Cannot exceed `desired_capacity`.
* `wait_for_capacity_timeout` - (Optional) Time (seconds) to wait for instances to report a 'HEALTHY' status. Useful for plans with multiple dependencies that take some time to initialize. Leave undefined or set to `0` to indicate no wait. This is ignored when updating with blue/green deployment. 
* `spot_percentage` - (Optional; Required if not using `ondemand_count`) The percentage of Spot instances that would spin up from the `desired_capacity` number.
* `ondemand_count` - (Optional; Required if not using `spot_percentage`) Number of on demand instances to launch in the group. All other instances will be spot instances. When this parameter is set the `spot_percentage` parameter is being ignored.
* `scaling_strategy` - (Optional) Set termination policy.
    * `terminate_at_end_of_billing_hour` - (Optional) Specify whether to terminate instances at the end of each billing hour.
    * `termination_policy` - (Optional) - Determines whether to terminate the newest instances when performing a scaling action. Valid values: `"default"`, `"newestInstance"`.

* `health_check_type` - (Optional) The service that will perform health checks for the instance. Valid values: `"ELB"`, `"HCS"`, `"TARGET_GROUP"`, `"MLB"`, `"EC2"`, `"MULTAI_TARGET_SET"`, `"MLB_RUNTIME"`, `"K8S_NODE"`, `"NOMAD_NODE"`, `"ECS_CLUSTER_INSTANCE"`.
* `health_check_grace_period` - (Optional) The amount of time, in seconds, after the instance has launched to starts and check its health.
* `health_check_unhealthy_duration_before_replacement` - (Optional) The amount of time, in seconds, that we will wait before replacing an instance that is running and became unhealthy (this is only applicable for instances that were once healthy).

* `tags` - (Optional) A key/value mapping of tags to assign to the resource.
* `elastic_ips` - (Optional) A list of [AWS Elastic IP](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/elastic-ip-addresses-eip.html) allocation IDs to associate to the group instances.
    
* `revert_to_spot` - (Optional) Hold settings for strategy correction – replacing On-Demand for Spot instances. Supported Values: `"never"`, `"always"`, `"timeWindow"`
    * `perform_at` - (Required) In the event of a fallback to On-Demand instances, select the time period to revert back to Spot. Supported Arguments – always (default), timeWindow, never. For timeWindow or never to be valid the group must have availabilityOriented OR persistence defined.
    * `time_windows` - (Optional) Specify a list of time windows for to execute revertToSpot strategy. Time window format: `ddd:hh:mm-ddd:hh:mm`. Example: `Mon:03:00-Wed:02:30`

<a id="load-balancers"></a>
## Load Balancers
    
* `elastic_load_balancers` - (Optional) List of Elastic Load Balancers names (ELB).
* `target_group_arns` - (Optional) List of Target Group ARNs to register the instances to.
* `multai_target_sets` - (Optional) Set of targets to register. 
    * `target_set_id` - (Required) ID of Multai target set.
    * `balancer_id` - (Required) ID of Multai Load Balancer.
    
Usage:

```hcl
  elastic_load_balancers = ["bal5", "bal2"]
  target_group_arns = ["tg-arn"]
  multai_target_sets = [{
    target_set_id = "ts-123",
    balancer_id   = "bal-123"
  },
  {
    target_set_id = "ts-234",
    balancer_id   = "bal-234"
  }]
```


<a id="scheduled-task"></a>
## Scheduled Tasks

Each `scheduled_task` supports the following:

* `task_type`- (Required) The task type to run. Supported task types are: `"scale"`, `"backup_ami"`, `"roll"`, `"scaleUp"`, `"percentageScaleUp"`, `"scaleDown"`, `"percentageScaleDown"`, `"statefulUpdateCapacity"`.  //// todo sali need to udit
* `cron_expression` - (Optional; Required if not using `frequency`) A valid cron expression. The cron is running in UTC time zone and is in [Unix cron format](https://en.wikipedia.org/wiki/Cron).
* `start_time` - (Optional; Format: ISO 8601) Set a start time for one time tasks.
* `frequency` - (Optional; Required if not using `cron_expression`) The recurrence frequency to run this task. Supported values are `"hourly"`, `"daily"`, `"weekly"` and `"continuous"`.
* `is_enabled` - (Optional, Default: `true`) Setting the task to being enabled or disabled.

Usage:

```hcl
  scheduled_task = [{
    task_type             = "backup_ami"
    cron_expression       = ""
    start_time            = "1970-01-01T01:00:00Z"
    frequency             = "hourly"
    is_enabled            = false
  }]
```

<a id="network-interface"></a>
## Network Interfaces

Each of the `network_interface` attributes controls a portion of the AWS
Instance's "Elastic Network Interfaces". It's a good idea to familiarize yourself with [AWS's Elastic Network
Interfaces docs](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-eni.html)
to understand the implications of using these attributes.

* `network_interface_id` - (Optional) The ID of the network interface.
* `device_index` - (Required) The index of the device on the instance for the network interface attachment.
* `description` - (Required) The description of the network interface.
* `private_ip_address` - (Optional) The private IP address of the network interface.
* `delete_on_termination` - (Optional) If set to true, the interface is deleted when the instance is terminated.
* `secondary_private_ip_address_count` - (Optional) The number of secondary private IP addresses.
* `associate_public_ip_address` - (Optional) Indicates whether to assign a public IP address to an instance you launch in a VPC. The public IP address can only be assigned to a network interface for eth0, and can only be assigned to a new network interface, not an existing one.
* `associate_ipv6_address` - (Optional) Indicates whether to assign IPV6 addresses to your instance. Requires a subnet with IPV6 CIDR block ranges.

Usage:

```hcl
  network_interface = [{ 
    network_interface_id               = "" 
    device_index                       = 1
    description                        = "nic description in here"
    private_ip_address                 = "1.1.1.1"
    delete_on_termination              = false
    secondary_private_ip_address_count = 1
    associate_public_ip_address        = true
  }]
```



<a id="health-check"></a>
## Health Check

* `health_check_type` - (Optional) The service that will perform health checks for the instance. Supported values : `"ELB"`, `"HCS"`, `"TARGET_GROUP"`, `"CUSTOM"`, `"K8S_NODE"`, `"MLB"`, `"EC2"`, `"MULTAI_TARGET_SET"`, `"MLB_RUNTIME"`, `"K8S_NODE"`, `"NOMAD_NODE"`, `"ECS_CLUSTER_INSTANCE"`.
* `health_check_grace_period` - (Optional) The amount of time, in seconds, after the instance has launched to starts and check its health
* `health_check_unhealthy_duration_before_replacement` - (Optional) The amount of time, in seconds, that we will wait before replacing an instance that is running and became unhealthy (this is only applicable for instances that were once healthy)

Usage:

```hcl
  health_check_type                                  = "ELB" 
  health_check_grace_period                          = 100
  health_check_unhealthy_duration_before_replacement = 120
```

* `integration_route53` - (Optional) Describes the [Route53](https://aws.amazon.com/documentation/route53/?id=docs_gateway) integration.

    * `domains` - (Required) Collection of one or more domains to register.
        * `hosted_zone_id` - (Required) The id associated with a hosted zone.
        * `spotinst_acct_id` - (Optional) The Spotinst account ID that is linked to the AWS account that holds the Route 53 hosted Zone ID. The default is the user Spotinst account provided as a URL parameter.
        * `record_sets` - (Required) Collection of records containing authoritative DNS information for the specified domain name.
            * `name` - (Required) The record set name.
            * `use_public_ip` - (Optional, Default: `false`) - Designates if the IP address should be exposed to connections outside the VPC.

Usage:

```hcl
    integration_route53 = {
      domains = {
        hosted_zone_id   = "zone-id"
        spotinst_acct_id = "act-123456"
        
        record_sets = {
          name          = "foo.example.com"
          use_public_ip = true
        }
      }
    }
```

## Attributes Reference

The following attributes are exported:

* `id` - The group ID.
