package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// Subscription ID for Azure
var subscriptionID string = "39b9a5ef-d664-4a6e-abd9-8a86328936b6"

func TestAzureNICConnection(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "shak0039",
		},
		MaxRetries:         10, // Increased retries for better stability
		TimeBetweenRetries: 30 * time.Second, // Wait between retries to handle Azure slowness
	}

	// Run `terraform init` and `terraform apply`
	fmt.Println("Running Terraform Apply...")
	terraform.InitAndApply(t, terraformOptions)

	// Retrieve Terraform output values
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name") // FIXED: Now retrieving NIC name

	// Confirm VM exists
	fmt.Println("Checking if the VM exists in Azure...")
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Confirm NIC exists
	fmt.Println("Checking if the NIC exists in Azure...")
	assert.True(t, azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID), "NIC should exist")

	// - Retrieve VM OS details from Terraform
	vmVersion := terraform.Output(t, terraformOptions, "vm_version")
	vmSKU := terraform.Output(t, terraformOptions, "vm_os")

	// Debugging: Print the retrieved values before asserting
	fmt.Println("Debugging VM OS Details:")
	fmt.Println("    - VM Version:", vmVersion)
	fmt.Println("    - VM SKU:", vmSKU)

	// - Validate vm version and sku
    assert.Equal(t, "latest", vmVersion)
    assert.Equal(t, "22_04-lts-gen2", vmSKU)

	// Destroy VM first to prevent NIC dependency issues
	fmt.Println("Destroying VM first before deleting the NIC...")
	terraform.RunTerraformCommand(t, terraformOptions, "destroy", "-target=azurerm_linux_virtual_machine.webserver", "-auto-approve")

	// Wait for 3 minutes (180 seconds) to allow Azure to fully detach the NIC
	fmt.Println("Waiting 180 seconds for Azure to fully release the NIC before deleting it...")
	time.Sleep(180 * time.Second)

	// Now destroy everything else
	fmt.Println("Running full Terraform Destroy for remaining resources...")
	terraform.Destroy(t, terraformOptions)

	fmt.Println("Test Completed Successfully!")
}