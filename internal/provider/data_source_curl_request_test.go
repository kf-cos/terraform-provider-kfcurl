package provider

import (
	"bytes"
	"fmt"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccdataSourceCurlRequest(t *testing.T) {

	rName := sdkacctest.RandomWithPrefix("devopsrob")
	json := `{"name": "` + rName + `"}`

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"POST",
		"https://example.com/create",
		httpmock.NewStringResponder(200, `{"name": "devopsrob"}`),
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccdataSourceCurlRequest(json),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.kfcurl_request.test", "name", regexp.MustCompile("^leader")),
					resource.TestMatchResourceAttr(
						"data.kfcurl_request.test", "response", regexp.MustCompile(`^{"name": "devopsrob"}`)),
					resource.TestMatchResourceAttr(
						"data.kfcurl_request.test", "url", regexp.MustCompile("^https://example.com")),
					resource.TestMatchResourceAttr(
						"data.kfcurl_request.test", "request_body", regexp.MustCompile(json)),
					resource.TestMatchResourceAttr(
						"data.kfcurl_request.test", "method", regexp.MustCompile("^POST")),
				),
			},
		},
	})
}

func testAccdataSourceCurlRequest(body string) string {
	return fmt.Sprintf(`
data "kfcurl_request" "test" {
  name           = "leader"
  url            = "https://example.com/create"
  request_body = <<EOF
%s
EOF
  method         = "POST"
  response_codes = [200]
}
`, body)
}

func TestAccdataSourceRetriesOnFailure(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix("devopsrob")
	json := `{"name": "` + rName + `"}`

	mc := &mockClient{}

	resp1 := &http.Response{}
	resp1.StatusCode = http.StatusInternalServerError
	resp1.Body = ioutil.NopCloser(bytes.NewReader([]byte("boom")))

	resp2 := &http.Response{}
	resp2.StatusCode = http.StatusOK
	resp2.Body = ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mc.On("Do", mock.Anything).Return(resp1, nil).Once()
	mc.On("Do", mock.Anything).Return(resp2, nil)

	Client = mc

	// ensure default client is replaced after test
	defer func() {
		Client = &http.Client{}
	}()

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				PlanOnly: true, // Read is called for the plan and apply phase, we only need to test read once
				Config:   testAccdataSourceCurlBodyWithRetry(json),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						if len(mc.Calls) != 2 {
							return fmt.Errorf("expected http request to be made 2 times. It was made %v times", len(mc.Calls))
						}
						return nil
					},
				),
			},
		},
	})
}

func testAccdataSourceCurlBodyWithRetry(body string) string {
	return fmt.Sprintf(`
data "kfcurl_request" "test" {
 name           = "leader"
 url            = "https://example.com/create"
 response_codes = ["200"]

 request_body = <<EOF
%s
EOF

 retry_interval = 1
 max_retry 	 	= 1
 method         = "POST"
}
`, body)

}
