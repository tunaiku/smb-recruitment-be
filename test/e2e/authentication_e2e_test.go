package e2e

import (
	"net/http"
	"testing"

	httpexpect "github.com/gavv/httpexpect/v2"
	"github.com/tunaiku/mobilebanking/test/e2e/setup"
	"go.uber.org/dig"
)

var (
	container = dig.New()
)

func TestAuthenticateEndpoint_Should_ReturnAccessToken_When_CredentialValid(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		e.POST("/auth/authenticate").WithJSON(map[string]interface{}{
			"username": "john",
			"password": "123456",
		}).Expect().JSON().Path("$.access_token").NotNull()
	})
}

func TestAuthenticateEndpoint_Should_ReturnHttpStatusBadRequest_When_PasswordInvalid(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		e.POST("/auth/authenticate").WithJSON(map[string]interface{}{
			"username": "john",
			"password": "12345",
		}).Expect().Status(http.StatusBadRequest)
	})
}

func TestAuthenticateEndpoint_Should_ReturnHttpStatusBadRequest_When_UsernameInvalid(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		e.POST("/auth/authenticate").WithJSON(map[string]interface{}{
			"username": "johe",
			"password": "123456",
		}).Expect().Status(http.StatusBadRequest)
	})
}
