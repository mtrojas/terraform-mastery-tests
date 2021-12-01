package test

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestHelloWorldAppExample(t *testing.T) {
	t.Parallel()

	opts := &terraform.Options{
		TerraformDir: "/Users/mtrojas/Documents/terraform-mastery/modules/examples/hello-world-app",
		Vars: map[string]interface{}{
			"mysql_config": map[string]interface{}{
				"address": "mock-value-for-test",
				"port":    3306,
			},
		},
	}

	// Clean up everything at the end of the test
	defer terraform.Destroy(t, opts)
	terraform.InitAndApply(t, opts)

	albDnsName := terraform.OutputRequired(t, opts, "alb_dns_name")
	url := fmt.Sprintf("http://%s", albDnsName)

	maxRetries := 10
	timeBetweenRetries := 10 * time.Second
	tlsConfig := tls.Config{}

	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		url,
		&tlsConfig,
		maxRetries,
		timeBetweenRetries,
		func(status int, body string) bool {
			return status == 200 &&
				strings.Contains(body, "My Example Deployment v0.0.1 of Hello World App")
		},
	)
}
