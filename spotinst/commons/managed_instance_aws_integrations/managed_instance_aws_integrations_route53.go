package managed_instance_aws_integrations

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spotinst/spotinst-sdk-go/service/managedinstance/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/terraform-providers/terraform-provider-spotinst/spotinst/commons"
)

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Setup
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func SetupRoute53(fieldsMap map[commons.FieldName]*commons.GenericField) {
	fieldsMap[IntegrationRoute53] = commons.NewGenericField(
		commons.ManagedInstanceAwsIntegrations,
		IntegrationRoute53,
		&schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					string(Domains): {
						Type:     schema.TypeSet,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								string(HostedZoneId): {
									Type:     schema.TypeString,
									Required: true,
								},

								string(SpotinstAcctID): {
									Type:     schema.TypeString,
									Optional: true,
								},

								string(RecordSets): {
									Type:     schema.TypeSet,
									Required: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											string(UsePublicIP): {
												Type:     schema.TypeBool,
												Optional: true,
											},

											string(Name): {
												Type:     schema.TypeString,
												Required: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			return nil
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			if v, ok := resourceData.GetOk(string(IntegrationRoute53)); ok {
				if integration, err := expandAWSGroupRoute53Integration(v); err != nil {
					return err
				} else {
					managedInstance.Integration.SetRoute53(integration)
				}
			}
			return nil
		},

		func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
			miWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
			managedInstance := miWrapper.GetManagedInstance()
			var value *aws.Route53Integration = nil

			if v, ok := resourceData.GetOk(string(IntegrationRoute53)); ok {
				if integration, err := expandAWSGroupRoute53Integration(v); err != nil {
					return err
				} else {
					value = integration
				}
			}
			managedInstance.Integration.SetRoute53(value)
			return nil
		},

		nil,
	)

	//fieldsMap[ElasticLoadBalancers] = commons.NewGenericField(
	//	commons.ElastigroupAWS,
	//	ElasticLoadBalancers,
	//	&schema.Schema{
	//		Type:     schema.TypeList,
	//		Elem:     &schema.Schema{Type: schema.TypeString},
	//		Optional: true,
	//	},
	//	func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
	//		egWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
	//		elastigroup := egWrapper.GetManagedInstance()
	//		var balNames []string = nil
	//		if elastigroup.Integration != nil && elastigroup.Integration.LoadBalancersConfig != nil &&
	//			elastigroup.Integration.LoadBalancersConfig.LoadBalancers != nil {
	//
	//			balancers := elastigroup.Integration.LoadBalancersConfig.LoadBalancers
	//			for _, balancer := range balancers {
	//				balType := spotinst.StringValue(balancer.Type)
	//				if strings.ToUpper(balType) == string(BalancerTypeClassic) {
	//					balName := spotinst.StringValue(balancer.Name)
	//					balNames = append(balNames, balName)
	//				}
	//			}
	//		}
	//		resourceData.Set(string(ElasticLoadBalancers), balNames)
	//		return nil
	//	},
	//	func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
	//		egWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
	//		elastigroup := egWrapper.GetManagedInstance()
	//		if balNames, ok := resourceData.GetOk(string(ElasticLoadBalancers)); ok {
	//			var fn = func(name string) (*aws.LoadBalancer, error) {
	//				return &aws.LoadBalancer{
	//					Type: spotinst.String(strings.ToUpper(string(BalancerTypeClassic))),
	//					Name: spotinst.String(name),
	//				}, nil
	//			}
	//			if elbBalancers, err := expandBalancersContent(balNames, fn); err != nil {
	//				return err
	//			} else if elbBalancers != nil && len(elbBalancers) > 0 {
	//				existingBalancers := elastigroup.Integration.LoadBalancersConfig.LoadBalancers
	//				if existingBalancers != nil && len(existingBalancers) > 0 {
	//					elbBalancers = append(existingBalancers, elbBalancers...)
	//				}
	//				elastigroup.Integration.LoadBalancersConfig.SetLoadBalancers(elbBalancers)
	//			}
	//		}
	//		return nil
	//	},
	//	func(resourceObject interface{}, resourceData *schema.ResourceData, meta interface{}) error {
	//		egWrapper := resourceObject.(*commons.MangedInstanceAWSWrapper)
	//		if err := onBalancersUpdate(egWrapper, resourceData); err != nil {
	//			return err
	//		}
	//		return nil
	//	},
	//	nil,
	//)
}

//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
//            Utils
//-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
func expandAWSGroupRoute53Integration(data interface{}) (*aws.Route53Integration, error) {
	integration := &aws.Route53Integration{}
	list := data.([]interface{})

	if list != nil && list[0] != nil {
		m := list[0].(map[string]interface{})

		if v, ok := m[string(Domains)]; ok {
			domains, err := expandAWSGroupRoute53IntegrationDomains(v)

			if err != nil {
				return nil, err
			}
			integration.SetDomains(domains)
		}
	}
	return integration, nil
}

func expandAWSGroupRoute53IntegrationDomains(data interface{}) ([]*aws.Domain, error) {
	list := data.(*schema.Set).List()
	domains := make([]*aws.Domain, 0, len(list))

	for _, v := range list {
		attr, ok := v.(map[string]interface{})
		domain := &aws.Domain{}

		if !ok {
			continue
		}

		if v, ok := attr[string(HostedZoneId)].(string); ok && v != "" {
			domain.SetHostedZoneID(spotinst.String(v))
		}

		if v, ok := attr[string(SpotinstAcctID)].(string); ok && v != "" {
			domain.SetSpotinstAccountID(spotinst.String(v))
		}

		if r, ok := attr[string(RecordSets)]; ok {
			if recordSets, err := expandAWSGroupRoute53IntegrationDomainsRecordSets(r); err != nil {
				return nil, err
			} else {
				domain.SetRecordSets(recordSets)
			}
		}
		domains = append(domains, domain)
	}
	return domains, nil
}

func expandAWSGroupRoute53IntegrationDomainsRecordSets(data interface{}) ([]*aws.RecordSet, error) {
	list := data.(*schema.Set).List()
	recordSets := make([]*aws.RecordSet, 0, len(list))

	for _, v := range list {
		attr, ok := v.(map[string]interface{})

		if !ok {
			continue
		}

		if _, ok := attr[string(UsePublicIP)]; !ok {
			return nil, errors.New("invalid record set attributes: use_public_ip missing")
		}

		if _, ok := attr[string(Name)]; !ok {
			return nil, errors.New("invalid record set attributes: name missing")
		}

		recordSet := &aws.RecordSet{
			UsePublicIP: spotinst.Bool(attr[string(UsePublicIP)].(bool)),
			Name:        spotinst.String(attr[string(Name)].(string)),
		}

		recordSets = append(recordSets, recordSet)
	}
	return recordSets, nil
}

func expandBalancersContent(balancersIdentifiers interface{}, fn CreateBalancerObjFunc) ([]*aws.LoadBalancer, error) {
	if balancersIdentifiers == nil {
		return nil, nil
	}
	list := balancersIdentifiers.([]interface{})
	balancers := make([]*aws.LoadBalancer, 0, len(list))
	for _, str := range list {
		if id, ok := str.(string); ok && id != "" {
			if lb, err := fn(id); err != nil {
				return nil, err
			} else {
				balancers = append(balancers, lb)
			}
		}
	}
	return balancers, nil
}

type CreateBalancerObjFunc func(id string) (*aws.LoadBalancer, error)

func onBalancersUpdate(egWrapper *commons.ElastigroupWrapper, resourceData *schema.ResourceData) error {
	var elbNullify = false
	var tgNullify = false
	var mlbNullify = false

	elastigroup := egWrapper.GetElastigroup()

	//if !egWrapper.StatusElbUpdated {
	//	if elbBalancers, err := extractBalancers(BalancerTypeClassic, elastigroup, resourceData); err != nil {
	//		return err
	//	} else if elbBalancers != nil && len(elbBalancers) > 0 {
	//		existingBalancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers
	//		if existingBalancers != nil && len(existingBalancers) > 0 {
	//			elbBalancers = append(existingBalancers, elbBalancers...)
	//		}
	//		elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.SetLoadBalancers(elbBalancers)
	//	} else {
	//		elbNullify = true
	//	}
	//	egWrapper.StatusElbUpdated = true
	//}
	//if !egWrapper.StatusTgUpdated {
	//	if tgBalancers, err := extractBalancers(BalancerTypeTargetGroup, elastigroup, resourceData); err != nil {
	//		return err
	//	} else if tgBalancers != nil && len(tgBalancers) > 0 {
	//		existingBalancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers
	//		if existingBalancers != nil && len(existingBalancers) > 0 {
	//			tgBalancers = append(existingBalancers, tgBalancers...)
	//		}
	//		elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.SetLoadBalancers(tgBalancers)
	//	} else {
	//		tgNullify = true
	//	}
	//	egWrapper.StatusTgUpdated = true
	//}
	//if !egWrapper.StatusMlbUpdated {
	//	if mlbBalancers, err := extractBalancers(BalancerTypeMultaiTargetSet, elastigroup, resourceData); err != nil {
	//		return err
	//	} else if mlbBalancers != nil && len(mlbBalancers) > 0 {
	//		existingBalancers := elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.LoadBalancers
	//		if existingBalancers != nil && len(existingBalancers) > 0 {
	//			mlbBalancers = append(existingBalancers, mlbBalancers...)
	//		}
	//		elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.SetLoadBalancers(mlbBalancers)
	//	} else {
	//		mlbNullify = true
	//	}
	//	egWrapper.StatusMlbUpdated = true
	//}

	// All fields share the same object structure, we need to nullify if and only if there are no items
	// from all types
	if elbNullify && tgNullify && mlbNullify {
		elastigroup.Compute.LaunchSpecification.LoadBalancersConfig.SetLoadBalancers(nil)
	}
	return nil
}
