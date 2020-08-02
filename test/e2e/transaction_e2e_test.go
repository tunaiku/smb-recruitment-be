package e2e_test

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"text/template"

	httpexpect "github.com/gavv/httpexpect/v2"
	"github.com/tunaiku/mobilebanking/test/e2e/setup"
)

type testCase struct {
	desc                 string
	payload              map[string]interface{}
	pathVariables        map[string]interface{}
	responseHTTPStatus   int
	responseBodyExpecter func(*httpexpect.Response)
}

type transactionEndpointTestTable struct {
	testCases  []testCase
	endpoint   string
	httpExpect *httpexpect.Expect
	httpMethod string
}

func (tbl *transactionEndpointTestTable) runTests(t *testing.T) {

	for _, tC := range tbl.testCases {
		tmpl, err := template.New("request").Parse(tbl.endpoint)
		if err != nil {
			t.Error(err)
		}
		var output bytes.Buffer
		err = tmpl.Execute(&output, tC.pathVariables)
		if err != nil {
			panic(err)
		}
		requestEndpoint := string(output.Bytes())
		method := strings.ToUpper(tbl.httpMethod)
		testCaseName := fmt.Sprintln(method, " ", requestEndpoint, " ", tC.desc)
		t.Run(testCaseName, func(t *testing.T) {

			accessToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmYzU1ZTNhOC1jMGZiLTQwYzctYWI4YS05Y2RhM2ZjYTQwZDQifQ.2oM9B0sTpIlgN-zvDGyrnaNJDiIIU6eIgiko7NxZj2s"

			resp := tbl.httpExpect.Request(method, requestEndpoint).WithJSON(tC.payload).WithHeader("Authorization", accessToken).Expect()
			resp.Status(tC.responseHTTPStatus)
			expect := tC.responseBodyExpecter
			if expect != nil {
				expect(resp)
			}
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		testTable := transactionEndpointTestTable{
			endpoint:   "/transaction",
			httpMethod: "post",
			httpExpect: e,
			testCases: []testCase{
				{
					desc: " should create transaction when auth_method sets to `otp` and the parameter valid",
					payload: map[string]interface{}{
						"auth_method":         "otp",
						"amount":              "3000",
						"transaction_code":    "T001",
						"destination_account": "10001",
					},
					responseHTTPStatus: http.StatusCreated,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ContainsKey("transaction_id").NotEmpty()
					},
				},
				{
					desc: " should created transaction when auth_method sets to `pin` and the parameter valid",
					payload: map[string]interface{}{
						"auth_method":         "pin",
						"amount":              "3000",
						"transaction_code":    "T001",
						"destination_account": "10001",
					},
					responseHTTPStatus: http.StatusCreated,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ContainsKey("transaction_id").NotEmpty()
					},
				},
				{
					desc: " should be failed with '400' as http status code and {\"message\":\"amount not match the minimum transaction amount\"} when the amount not match the minimum transaction amount",
					payload: map[string]interface{}{
						"auth_method":         "pin",
						"amount":              "2000",
						"transaction_code":    "T001",
						"destination_account": "10001",
					},
					responseHTTPStatus: http.StatusBadRequest,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "amount not match the minimum transaction amount")
					},
				},
				{
					desc: " should be failed with '400' as http status code and {\"message\":\"unsupported authorization method\"} when auth_method sets to 'password'",
					payload: map[string]interface{}{
						"auth_method":         "password",
						"amount":              "2000",
						"transaction_code":    "T001",
						"destination_account": "10001",
					},
					responseHTTPStatus: http.StatusBadRequest,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "unsupported authorization method")
					},
				},
				{
					desc: " should be failed with '400' as http status code and {\"message\":\"authorization method not configured\"} when auth_method sets to 'otp' but the user not configure it",
					payload: map[string]interface{}{
						"auth_method":         "otp",
						"amount":              "2000",
						"transaction_code":    "T001",
						"destination_account": "10001",
					},
					responseHTTPStatus: http.StatusBadRequest,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "authorization method not configured")
					},
				},
				{
					desc: " should be failed with '400' as http status code and {\"message\":\"authorization method not configured\"} when auth_method sets to 'pin' but the user not configure it",
					payload: map[string]interface{}{
						"auth_method":         "pin",
						"amount":              "2000",
						"transaction_code":    "T001",
						"destination_account": "10001",
					},
					responseHTTPStatus: http.StatusBadRequest,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "authorization method not configured")
					},
				},
				{
					desc: " should be failed with '400' as http status code and {\"message\":\"transaction code not found\"} when transaction_code sets to 'T003' but that transaction is not found",
					payload: map[string]interface{}{
						"auth_method":         "pin",
						"amount":              "2000",
						"transaction_code":    "T001",
						"destination_account": "10001",
					},
					responseHTTPStatus: http.StatusBadRequest,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "transaction code not found")
					},
				},
				{
					desc: " should be failed with '400' as http status code and {\"message\":\"transaction not allowed\"} when transaction_code sets to 'T001' but the user has no previllage to do the transaction",
					payload: map[string]interface{}{
						"auth_method":         "pin",
						"amount":              "2000",
						"transaction_code":    "T001",
						"destination_account": "10001",
					},
					responseHTTPStatus: http.StatusBadRequest,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "transaction not allowed")
					},
				},
				{
					desc: " should be failed with '400' as http status code and {\"message\":\"destination account not found\"} when destination_account sets to '10003' but the that account not found",
					payload: map[string]interface{}{
						"auth_method":         "pin",
						"amount":              "2000",
						"transaction_code":    "T001",
						"destination_account": "10001",
					},
					responseHTTPStatus: http.StatusBadRequest,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "transaction not allowed")
					},
				},
			}}
		testTable.runTests(t)
	})
}

func TestVerifyTransaction(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		testTable := transactionEndpointTestTable{
			endpoint:   "/transaction/{{.ID}}/verify",
			httpMethod: "put",
			httpExpect: e,
			testCases: []testCase{
				{
					desc: " transaction should verified with `201-Accepted` when transaction id is valid and the credential is matches ",
					pathVariables: map[string]interface{}{
						"ID": "1111",
					},
					payload: map[string]interface{}{
						"credential": "123456",
					},
					responseHTTPStatus: http.StatusAccepted,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ContainsKey("transaction_id")
					},
				},
				{
					desc: " transaction should be failed with `400-Bad Request` and `{\"message\":\"invalid credential\"}`  when valid credential is invalid ",
					pathVariables: map[string]interface{}{
						"ID": "1111",
					},
					payload: map[string]interface{}{
						"credential": "1234",
					},
					responseHTTPStatus: http.StatusBadRequest,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "invalid credential")
					},
				},
				{
					desc: " transaction should be failed with `400-Bad Request` and `{\"message\":\"invalid transaction\"}`  when the transaction state is not `WaitAuthorization` ",
					pathVariables: map[string]interface{}{
						"ID": "1111",
					},
					payload: map[string]interface{}{
						"credential": "123456",
					},
					responseHTTPStatus: http.StatusBadRequest,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "invalid transaction")
					},
				},
				{
					desc: " transaction should be failed with `404-Not Found` and `{\"message\":\"transaction not found\"}`  when there is no transaction with belong to path variable id ",
					pathVariables: map[string]interface{}{
						"ID": "1112",
					},
					payload: map[string]interface{}{
						"credential": "1234",
					},
					responseHTTPStatus: http.StatusNotFound,
					responseBodyExpecter: func(resp *httpexpect.Response) {
						resp.JSON().Object().ValueEqual("message", "transaction not found")
					},
				},
			},
		}
		testTable.runTests(t)
	})
}
