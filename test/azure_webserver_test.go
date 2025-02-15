package test 

import (
	"testing"
	"time"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// Subscription ID for Azure
var subscriptionID string = "39b9a5ef-d664-4a6e-abd9-8a86328936b6"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": "shak0039",
		},
		MaxRetries:         10, // Increased retries for better stability
		TimeBetweenRetries: 30 * time.Second, // Wait between retries to handle Azure slowness
	}

	// Run `terraform init` and `terraform apply`
	fmt.Println("🚀 Running Terraform Apply...")
	terraform.InitAndApply(t, terraformOptions)

	// Retrieve Terraform output values
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	fmt.Println("✅ Checking if the VM exists in Azure...")
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	// Destroy VM first to prevent NIC dependency issues
	fmt.Println("🛑 Destroying VM first before deleting the NIC...")
	terraform.RunTerraformCommand(t, terraformOptions, "destroy", "-target=azurerm_linux_virtual_machine.webserver", "-auto-approve")

	// Wait for 3 minutes (180 seconds) to allow Azure to fully detach the NIC
	fmt.Println("⏳ Waiting 180 seconds for Azure to fully release the NIC before deleting it...")
	time.Sleep(180 * time.Second)

	// Now destroy everything else
	fmt.Println("🔥 Running full Terraform Destroy for remaining resources...")
	terraform.Destroy(t, terraformOptions)

	fmt.Println("✅ Test Completed Successfully!")
}



/*package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "39b9a5ef-d664-4a6e-abd9-8a86328936b6"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "shak0039",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
}
*/