package machine

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"
	powervsproviderv1 "github.com/openshift/cluster-api-provider-powervs/pkg/apis/powervsprovider/v1alpha1"
	awsclient "github.com/openshift/cluster-api-provider-powervs/pkg/client"
	machinev1 "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	mapierrors "github.com/openshift/machine-api-operator/pkg/controller/machine"
	"k8s.io/klog/v2"
)

// removeStoppedMachine removes all instances of a specific machine that are in a stopped state.
func removeStoppedMachine(machine *machinev1.Machine, client awsclient.Client) error {
	instance, err := client.GetInstanceByName(machine.Name)
	if err != nil && err != awsclient.ErrorInstanceNotFound {
		klog.Errorf("Error getting instance by name: %s, err: %v", machine.Name, err)
		return fmt.Errorf("error getting instance by name: %s, err: %v", machine.Name, err)
	} else if err == awsclient.ErrorInstanceNotFound {
		klog.Infof("Instance not found with name: %n", machine.Name)
		return nil
	}

	if instance != nil && *instance.Status == awsclient.InstanceStateNameShutoff {
		return client.DeleteInstance(*instance.PvmInstanceID)
	}
	return nil
}

func launchInstance(machine *machinev1.Machine, machineProviderConfig *powervsproviderv1.PowerVSMachineProviderConfig, userData []byte, client awsclient.Client) (*models.PVMInstance, error) {
	// code for powervs

	memory, err := strconv.ParseFloat(machineProviderConfig.Memory, 64)
	if err != nil {
		return nil, mapierrors.InvalidMachineConfiguration("failed to convert memory(%s) to float64", machineProviderConfig.Memory)
	}
	cores, err := strconv.ParseFloat(machineProviderConfig.Cores, 64)
	if err != nil {
		return nil, mapierrors.InvalidMachineConfiguration("failed to convert Cores(%s) to float64", machineProviderConfig.Cores)
	}

	var nets []*models.PVMInstanceAddNetwork

	for _, net := range machineProviderConfig.Subnets {
		nets = append(nets, &models.PVMInstanceAddNetwork{NetworkID: &net})
	}

	params := &p_cloud_p_vm_instances.PcloudPvminstancesPostParams{
		Body: &models.PVMInstanceCreate{
			ImageID:     &machineProviderConfig.ImageID,
			KeyPairName: *machineProviderConfig.KeyName,
			Networks:    nets,
			ServerName:  &machine.Name,
			Memory:      &memory,
			Processors:  &cores,
			ProcType:    &machineProviderConfig.ProcessorType,
			SysType:     machineProviderConfig.MachineType,
			UserData:    base64.StdEncoding.EncodeToString(userData),
		},
	}
	instances, err := client.CreateInstance(params)
	if err != nil {
		return nil, mapierrors.CreateMachine("error creating powervs instance: %v", err)
	}

	insIDs := make([]string, 0)
	for _, in := range *instances {
		insID := in.PvmInstanceID
		insIDs = append(insIDs, *insID)
	}

	if len(insIDs) == 0 {
		return nil, mapierrors.CreateMachine("error getting the instance ID post deployment for: %s", machine.Name)
	}

	instance, err := client.GetInstance(insIDs[0])
	if err != nil {
		return nil, mapierrors.CreateMachine("error getting the instance for ID: %s", insIDs[0])
	}
	return instance, nil

	//machineKey := runtimeclient.ObjectKey{
	//	Name:      machine.Name,
	//	Namespace: machine.Namespace,
	//}
	//amiID, err := getAMI(machineKey, machineProviderConfig.AMI, client)
	//if err != nil {
	//	return nil, mapierrors.InvalidMachineConfiguration("error getting AMI: %v", err)
	//}
	//
	//securityGroupsIDs, err := getSecurityGroupsIDs(machineProviderConfig.SecurityGroups, client)
	//if err != nil {
	//	return nil, mapierrors.InvalidMachineConfiguration("error getting security groups IDs: %v", err)
	//}
	//subnetIDs, err := getSubnetIDs(machineKey, machineProviderConfig.Subnet, machineProviderConfig.Placement.AvailabilityZone, client)
	//if err != nil {
	//	return nil, mapierrors.InvalidMachineConfiguration("error getting subnet IDs: %v", err)
	//}
	//if len(subnetIDs) > 1 {
	//	klog.Warningf("More than one subnet id returned, only first one will be used")
	//}
	//
	//// build list of networkInterfaces (just 1 for now)
	//var networkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
	//	{
	//		DeviceIndex:              aws.Int64(machineProviderConfig.DeviceIndex),
	//		AssociatePublicIpAddress: machineProviderConfig.PublicIP,
	//		SubnetId:                 subnetIDs[0],
	//		Groups:                   securityGroupsIDs,
	//	},
	//}
	//
	//blockDeviceMappings, err := getBlockDeviceMappings(machineKey, machineProviderConfig.BlockDevices, *amiID, client)
	//if err != nil {
	//	return nil, mapierrors.InvalidMachineConfiguration("error getting blockDeviceMappings: %v", err)
	//}
	//
	//clusterID, ok := getClusterID(machine)
	//if !ok {
	//	klog.Errorf("Unable to get cluster ID for machine: %q", machine.Name)
	//	return nil, mapierrors.InvalidMachineConfiguration("Unable to get cluster ID for machine: %q", machine.Name)
	//}
	//// Add tags to the created machine
	//tagList := buildTagList(machine.Name, clusterID, machineProviderConfig.Tags)
	//
	//tagInstance := &ec2.TagSpecification{
	//	ResourceType: aws.String("instance"),
	//	Tags:         tagList,
	//}
	//tagVolume := &ec2.TagSpecification{
	//	ResourceType: aws.String("volume"),
	//	Tags:         tagList,
	//}
	//
	//userDataEnc := base64.StdEncoding.EncodeToString(userData)
	//
	//var iamInstanceProfile *ec2.IamInstanceProfileSpecification
	//if machineProviderConfig.IAMInstanceProfile != nil && machineProviderConfig.IAMInstanceProfile.ID != nil {
	//	iamInstanceProfile = &ec2.IamInstanceProfileSpecification{
	//		Name: aws.String(*machineProviderConfig.IAMInstanceProfile.ID),
	//	}
	//}
	//
	//var placement *ec2.Placement
	//if machineProviderConfig.Placement.AvailabilityZone != "" && machineProviderConfig.Subnet.ID == nil {
	//	placement = &ec2.Placement{
	//		AvailabilityZone: aws.String(machineProviderConfig.Placement.AvailabilityZone),
	//	}
	//}
	//
	//instanceTenancy := machineProviderConfig.Placement.Tenancy
	//
	//switch instanceTenancy {
	//case "":
	//	// Do nothing when not set
	//case awsproviderv1.DefaultTenancy, awsproviderv1.DedicatedTenancy, awsproviderv1.HostTenancy:
	//	if placement == nil {
	//		placement = &ec2.Placement{}
	//	}
	//	tenancy := string(instanceTenancy)
	//	placement.Tenancy = &tenancy
	//default:
	//	return nil, mapierrors.CreateMachine("invalid instance tenancy: %s. Allowed options are: %s,%s,%s",
	//		instanceTenancy,
	//		awsproviderv1.DefaultTenancy,
	//		awsproviderv1.DedicatedTenancy,
	//		awsproviderv1.HostTenancy)
	//}
	//
	//inputConfig := ec2.RunInstancesInput{
	//	ImageId:      amiID,
	//	InstanceType: aws.String(machineProviderConfig.InstanceType),
	//	// Only a single instance of the AWS instance allowed
	//	MinCount:              aws.Int64(1),
	//	MaxCount:              aws.Int64(1),
	//	KeyName:               machineProviderConfig.KeyName,
	//	IamInstanceProfile:    iamInstanceProfile,
	//	TagSpecifications:     []*ec2.TagSpecification{tagInstance, tagVolume},
	//	NetworkInterfaces:     networkInterfaces,
	//	UserData:              &userDataEnc,
	//	Placement:             placement,
	//	InstanceMarketOptions: getInstanceMarketOptionsRequest(machineProviderConfig),
	//}
	//
	//if len(blockDeviceMappings) > 0 {
	//	inputConfig.BlockDeviceMappings = blockDeviceMappings
	//}
	//runResult, err := client.RunInstances(&inputConfig)
	//if err != nil {
	//	metrics.RegisterFailedInstanceCreate(&metrics.MachineLabels{
	//		Name:      machine.Name,
	//		Namespace: machine.Namespace,
	//		Reason:    err.Error(),
	//	})
	//	// we return InvalidMachineConfiguration for 4xx errors which by convention signal client misconfiguration
	//	// https://tools.ietf.org/html/rfc2616#section-6.1.1
	//	// https: //docs.aws.amazon.com/AWSEC2/latest/APIReference/errors-overview.html
	//	// https://docs.aws.amazon.com/sdk-for-go/api/aws/awserr/
	//	if _, ok := err.(awserr.Error); ok {
	//		if reqErr, ok := err.(awserr.RequestFailure); ok {
	//			if strings.HasPrefix(strconv.Itoa(reqErr.StatusCode()), "4") {
	//				klog.Infof("Error launching instance: %v", reqErr)
	//				return nil, mapierrors.InvalidMachineConfiguration("error launching instance: %v", reqErr.Message())
	//			}
	//		}
	//	}
	//	klog.Errorf("Error creating EC2 instance: %v", err)
	//	return nil, mapierrors.CreateMachine("error creating EC2 instance: %v", err)
	//}
	//
	//if runResult == nil || len(runResult.Instances) != 1 {
	//	klog.Errorf("Unexpected reservation creating instances: %v", runResult)
	//	return nil, mapierrors.CreateMachine("unexpected reservation creating instance")
	//}
	//
	//return runResult.Instances[0], nil
}
