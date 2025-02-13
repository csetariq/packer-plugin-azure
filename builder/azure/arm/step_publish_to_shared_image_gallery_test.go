// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arm

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"

	"github.com/hashicorp/packer-plugin-azure/builder/azure/common"
	"github.com/hashicorp/packer-plugin-azure/builder/azure/common/constants"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

func TestStepPublishToSharedImageGalleryShouldNotPublishForVhd(t *testing.T) {
	var testSubject = &StepPublishToSharedImageGallery{
		publish: func(context.Context, PublishArgs) (string, error) {
			return "test", nil
		},
		say:   func(message string) {},
		error: func(e error) {},
		toSIG: func() bool { return false },
	}

	stateBag := createTestStateBagStepPublishToSharedImageGalleryForVhd()
	var result = testSubject.Run(context.Background(), stateBag)
	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == true {
		t.Fatalf("Expected the step to not set stateBag['%s'], but it was.", constants.Error)
	}
}

func TestStepPublishToSharedImageGalleryShouldPublishForManagedImageWithSig(t *testing.T) {
	var testSubject = &StepPublishToSharedImageGallery{
		publish: func(context.Context, PublishArgs) (string, error) {
			return "", nil
		},
		say:   func(message string) {},
		error: func(e error) {},
		toSIG: func() bool { return true },
	}

	stateBag := createTestStateBagStepPublishToSharedImageGallery(true)
	var result = testSubject.Run(context.Background(), stateBag)
	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == true {
		t.Fatalf("Expected the step to not set stateBag['%s'], but it was.", constants.Error)
	}
}

func TestStepPublishToSharedImageGalleryShouldPublishForNonManagedImageWithSig(t *testing.T) {
	var actualPublishArgs PublishArgs
	expectedPublishArgs := PublishArgs{
		SubscriptionID:  "Unit Test: ManagedImageSubscription",
		ReplicaCount:    1,
		ReplicationMode: "Full",
		Location:        "Unit Test: Location",
		SourceID:        "Unit Test: VM ID",
		SharedImageGallery: SharedImageGalleryDestination{
			SigDestinationGalleryName:   "Unit Test: ManagedImageSharedGalleryName",
			SigDestinationImageName:     "Unit Test: ManagedImageSharedGalleryImageName",
			SigDestinationSubscription:  "Unit Test: ManagedImageSubscription",
			SigDestinationImageVersion:  "Unit Test: ManagedImageSharedGalleryImageVersion",
			SigDestinationResourceGroup: "Unit Test: ManagedImageSigPublishResourceGroup",
			SigDestinationReplicationRegions: []string{
				"ManagedImageSharedGalleryReplicationRegionA",
				"ManagedImageSharedGalleryReplicationRegionB",
			},
			SigDestinationStorageAccountType: "Standard_LRS",
		},
		Tags: map[string]string{"tag01": "Unit Test: Tags"},
	}
	var testSubject = &StepPublishToSharedImageGallery{
		publish: func(ctx context.Context, args PublishArgs) (string, error) {
			actualPublishArgs = args
			return "", nil
		},
		say:   func(message string) {},
		error: func(e error) {},
		toSIG: func() bool { return true },
	}

	stateBag := createTestStateBagStepPublishToSharedImageGallery(false)
	var result = testSubject.Run(context.Background(), stateBag)
	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == true {
		t.Fatalf("Expected the step to not set stateBag['%s'], but it was.", constants.Error)
	}

	if diff := cmp.Diff(actualPublishArgs, expectedPublishArgs, []cmp.Option{
		cmpopts.IgnoreUnexported(PublishArgs{}),
	}...); diff != "" {
		t.Fatalf("Unexpected diff %s", diff)
	}
}

func TestStepPublishToSharedImageGalleryShouldPublishWithShallowReplication(t *testing.T) {
	var actualPublishArgs PublishArgs
	expectedPublishArgs := PublishArgs{
		SubscriptionID:  "Unit Test: ManagedImageSubscription",
		ReplicaCount:    1,
		ReplicationMode: "Shallow",
		Location:        "Unit Test: Location",
		SourceID:        "Unit Test: VM ID",
		SharedImageGallery: SharedImageGalleryDestination{
			SigDestinationGalleryName:   "Unit Test: ManagedImageSharedGalleryName",
			SigDestinationImageName:     "Unit Test: ManagedImageSharedGalleryImageName",
			SigDestinationSubscription:  "Unit Test: ManagedImageSubscription",
			SigDestinationImageVersion:  "Unit Test: ManagedImageSharedGalleryImageVersion",
			SigDestinationResourceGroup: "Unit Test: ManagedImageSigPublishResourceGroup",
			SigDestinationReplicationRegions: []string{
				"ManagedImageSharedGalleryReplicationRegionA",
				"ManagedImageSharedGalleryReplicationRegionB",
			},
			SigDestinationStorageAccountType: "Standard_LRS",
		},
		Tags: map[string]string{"tag01": "Unit Test: Tags"},
	}
	var testSubject = &StepPublishToSharedImageGallery{
		publish: func(ctx context.Context, args PublishArgs) (string, error) {
			actualPublishArgs = args
			return "", nil
		},
		say:   func(message string) {},
		error: func(e error) {},
		toSIG: func() bool { return true },
	}

	stateBag := createTestStateBagStepPublishToSharedImageGallery(false)
	stateBag.Put(constants.ArmSharedImageGalleryDestinationShallowReplication, true)
	var result = testSubject.Run(context.Background(), stateBag)
	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == true {
		t.Fatalf("Expected the step to not set stateBag['%s'], but it was.", constants.Error)
	}

	if diff := cmp.Diff(actualPublishArgs, expectedPublishArgs, []cmp.Option{
		cmpopts.IgnoreUnexported(PublishArgs{}),
	}...); diff != "" {
		t.Fatalf("Unexpected diff %s", diff)
	}
}

func TestStepPublishToSharedImageGalleryShouldPublishWithReplicationCount(t *testing.T) {
	var actualPublishArgs PublishArgs
	expectedPublishArgs := PublishArgs{
		SubscriptionID:  "Unit Test: ManagedImageSubscription",
		ReplicaCount:    5,
		ReplicationMode: "Full",
		Location:        "Unit Test: Location",
		SourceID:        "Unit Test: VM ID",
		SharedImageGallery: SharedImageGalleryDestination{
			SigDestinationGalleryName:   "Unit Test: ManagedImageSharedGalleryName",
			SigDestinationImageName:     "Unit Test: ManagedImageSharedGalleryImageName",
			SigDestinationSubscription:  "Unit Test: ManagedImageSubscription",
			SigDestinationImageVersion:  "Unit Test: ManagedImageSharedGalleryImageVersion",
			SigDestinationResourceGroup: "Unit Test: ManagedImageSigPublishResourceGroup",
			SigDestinationReplicationRegions: []string{
				"ManagedImageSharedGalleryReplicationRegionA",
				"ManagedImageSharedGalleryReplicationRegionB",
			},
			SigDestinationStorageAccountType: "Standard_LRS",
		},
		Tags: map[string]string{"tag01": "Unit Test: Tags"},
	}
	var testSubject = &StepPublishToSharedImageGallery{
		publish: func(ctx context.Context, args PublishArgs) (string, error) {
			actualPublishArgs = args
			return "", nil
		},
		say:   func(message string) {},
		error: func(e error) {},
		toSIG: func() bool { return true },
	}

	stateBag := createTestStateBagStepPublishToSharedImageGallery(false)
	stateBag.Put(constants.ArmManagedImageSharedGalleryImageVersionReplicaCount, int64(5))
	var result = testSubject.Run(context.Background(), stateBag)
	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == true {
		t.Fatalf("Expected the step to not set stateBag['%s'], but it was.", constants.Error)
	}

	if diff := cmp.Diff(actualPublishArgs, expectedPublishArgs, []cmp.Option{
		cmpopts.IgnoreUnexported(PublishArgs{}),
	}...); diff != "" {
		t.Fatalf("Unexpected diff %s", diff)
	}
}

func createTestStateBagStepPublishToSharedImageGallery(managed bool) multistep.StateBag {
	stateBag := new(multistep.BasicStateBag)

	stateBag.Put(constants.ArmManagedImageSigPublishResourceGroup, "Unit Test: ManagedImageSigPublishResourceGroup")
	stateBag.Put(constants.ArmManagedImageSharedGalleryName, "Unit Test: ManagedImageSharedGalleryName")
	stateBag.Put(constants.ArmManagedImageSharedGalleryImageName, "Unit Test: ManagedImageSharedGalleryImageName")
	stateBag.Put(constants.ArmManagedImageSharedGalleryImageVersion, "Unit Test: ManagedImageSharedGalleryImageVersion")
	stateBag.Put(constants.ArmLocation, "Unit Test: Location")
	value := "Unit Test: Tags"
	tags := map[string]string{
		"tag01": value,
	}
	stateBag.Put(constants.ArmTags, tags)
	stateBag.Put(constants.ArmManagedImageSharedGalleryReplicationRegions, []string{"ManagedImageSharedGalleryReplicationRegionA", "ManagedImageSharedGalleryReplicationRegionB"})
	stateBag.Put(constants.ArmManagedImageSharedGalleryImageVersionStorageAccountType, "Standard_LRS")
	if managed {
		stateBag.Put(constants.ArmManagedImageResourceGroupName, "Unit Test: ManagedImageResourceGroupName")
		stateBag.Put(constants.ArmManagedImageName, "Unit Test: ManagedImageName")
	} else {
		stateBag.Put(constants.ArmImageParameters, &images.Image{Properties: &images.ImageProperties{
			SourceVirtualMachine: &images.SubResource{Id: common.StringPtr("Unit Test: VM ID")},
		}})
	}
	stateBag.Put(constants.ArmManagedImageSubscription, "Unit Test: ManagedImageSubscription")
	stateBag.Put(constants.ArmSharedImageGalleryDestinationSubscription, "Unit Test: ManagedImageSubscription")
	stateBag.Put(constants.ArmIsManagedImage, managed)
	stateBag.Put(constants.ArmIsSIGImage, true)
	stateBag.Put(constants.ArmSharedImageGalleryDestinationShallowReplication, false)

	return stateBag
}

func createTestStateBagStepPublishToSharedImageGalleryForVhd() multistep.StateBag {
	stateBag := new(multistep.BasicStateBag)

	stateBag.Put(constants.ArmLocation, "Unit Test: Location")
	value := "Unit Test: Tags"
	tags := map[string]*string{
		"tag01": &value,
	}
	stateBag.Put(constants.ArmTags, tags)

	return stateBag
}
