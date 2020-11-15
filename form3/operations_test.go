package form3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

const baseURLPath = "/localhost:8080"

var attr = &AccountCreateRequestAttributes{
	Country:               "GB",
	BaseCurrency:          "GBP",
	BankID:                "400302",
	BankIDCode:            "GBDSC",
	AccountNumber:         "10000004",
	CustomerID:            "234",
	Iban:                  "GB28NWBK40030212764204",
	Bic:                   "NWBKGB42",
	AccountClassification: "Personal",
}

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	client = NewClient(nil)
	u, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = u
	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func TestAccountService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
		v := &AccountCreateRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Error(err)
		}

		testMethod(t, r, "POST")
		testHeader(t, r, "Content-Type", "application/json")
		testHeader(t, r, "Accept", "application/vnd.api+json")
		w.WriteHeader(http.StatusCreated)

		want := &AccountCreateRequest{&AccountCreateRequestData{
			attr,
			"d91afcdb-62d2-4185-b23d-71c98eaab812",
			"a68eddcd-6eec-4b5e-846d-97b1161248e2",
			"accounts",
		}}
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Request body = %+v, want %+v", v, want)
		}
		_, e := fmt.Fprint(w, `{
								   "data":{
									  "attributes":{
										 "account_classification":"Personal",
										 "account_number":"10000004",
										 "alternative_bank_account_names":null,
										 "bank_id":"400302",
										 "bank_id_code":"GBDSC",
										 "base_currency":"GBP",
										 "bic":"NWBKGB42",
										 "country":"GB",
										 "customer_id":"234",
										 "iban":"GB28NWBK40030212764204"
									  },
									  "created_on":"2020-11-11T10:40:44.709Z",
									  "id":"a68eddcd-6eec-4b5e-846d-97b1161248e2",
									  "modified_on":"2020-11-11T10:40:44.709Z",
									  "organisation_id":"d91afcdb-62d2-4185-b23d-71c98eaab812",
									  "type":"accounts",
									  "version":0
								   },
								   "links":{
									  "self":"/v1/organisation/accounts/a68eddcd-6eec-4b5e-846d-97b1161248e2"
								   }
								}`)
		if e != nil {
			t.Error(e)
		}
	})

	account, _, _, err := client.Account.Create(context.Background(),
		"a68eddcd-6eec-4b5e-846d-97b1161248e2",
		"d91afcdb-62d2-4185-b23d-71c98eaab812", attr)
	if err != nil {
		t.Errorf("Account.Create returned error: %v", err)
	}

	want := &Account{
		Type:           "accounts",
		ID:             "a68eddcd-6eec-4b5e-846d-97b1161248e2",
		OrganisationID: "d91afcdb-62d2-4185-b23d-71c98eaab812",
		Version:        0,
		CreatedOn:      time.Date(2020, time.November, 11, 10, 40, 44, 709000000, time.UTC),
		ModifiedOn:     time.Date(2020, time.November, 11, 10, 40, 44, 709000000, time.UTC),
		Attributes: &AccountAttributes{
			Country:                     "GB",
			BaseCurrency:                "GBP",
			BankID:                      "400302",
			BankIDCode:                  "GBDSC",
			AccountNumber:               "10000004",
			CustomerID:                  "234",
			Iban:                        "GB28NWBK40030212764204",
			AccountClassification:       "Personal",
			Bic:                         "NWBKGB42",
			AlternativeBankAccountNames: nil,
		}}
	if !reflect.DeepEqual(account, want) {
		t.Errorf("Account.Create returned %+v, want %+v", account, want)
	}
}

func TestAccountService_Fetch(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, e := fmt.Fprint(w, `{
									  "data": {
										"attributes": {
										  "account_classification": "Personal",
										  "account_number": "41426819",
										  "alternative_bank_account_names": null,
										  "bank_id": "400300",
										  "bank_id_code": "GBDSC",
										  "base_currency": "GBP",
										  "bic": "NWBKGB22",
										  "country": "GB",
										  "customer_id": "234",
										  "iban": "GB28NWBK40030212764204"
										},
										"created_on": "2020-11-11T10:40:44.709Z",
										"id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
										"modified_on": "2020-11-11T10:40:44.709Z",
										"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
										"type": "accounts",
										"version": 0
									  },
									  "links": {
										"self": "/v1/organisation/accounts/a68eddcd-6eec-4b5e-846d-97b1161248e2"
									  }
									}`)
		if e != nil {
			t.Error(e)
		}
	})

	account, _, _, err := client.Account.Fetch(context.Background(), "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
	if err != nil {
		t.Errorf("Account.Fetch returned error: %v", err)
	}

	want := &Account{
		Type:           "accounts",
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Version:        0,
		CreatedOn:      time.Date(2020, time.November, 11, 10, 40, 44, 709000000, time.UTC),
		ModifiedOn:     time.Date(2020, time.November, 11, 10, 40, 44, 709000000, time.UTC),
		Attributes: &AccountAttributes{
			AccountClassification:       "Personal",
			AccountNumber:               "41426819",
			AlternativeBankAccountNames: nil,
			BankID:                      "400300",
			BankIDCode:                  "GBDSC",
			BaseCurrency:                "GBP",
			Bic:                         "NWBKGB22",
			Country:                     "GB",
			CustomerID:                  "234",
			Iban:                        "GB28NWBK40030212764204",
		}}
	if !reflect.DeepEqual(account, want) {
		t.Errorf("Account.Fetch returned %+v, want %+v", account, want)
	}
}

func TestAccountService_List_NoPages(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/organisation/accounts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, `{
									  "data": [
										{
										  "attributes": {
											"account_classification": "Personal",
											"account_number": "10000004",
											"alternative_bank_account_names": null,
											"bank_id": "400302",
											"bank_id_code": "GBDSC",
											"base_currency": "GBP",
											"bic": "NWBKGB42",
											"country": "GB",
											"customer_id": "234",
											"iban": "GB28NWBK40030212764204"
										  },
										  "created_on": "2020-11-11T10:40:44.709Z",
										  "id": "a68eddcd-6eec-4b5e-846d-97b1161248e2",
										  "modified_on": "2020-11-11T10:40:44.709Z",
										  "organisation_id": "d91afcdb-62d2-4185-b23d-71c98eaab812",
										  "type": "accounts",
										  "version": 0
										}
									  ],
									  "links": {
										"first": "/v1/organisation/accounts?page%5Bnumber%5D=first",
										"last": "/v1/organisation/accounts?page%5Bnumber%5D=last",
										"self": "/v1/organisation/accounts"
									  }
									}`)
		if err != nil {
			t.Error(err)
		}
	})

	accounts, _, _, err := client.Account.List(context.Background(), nil)
	if err != nil {
		t.Errorf("Account.List returned error: %v", err)
	}

	account := accounts[0]

	want := &Account{
		Type:           "accounts",
		ID:             "a68eddcd-6eec-4b5e-846d-97b1161248e2",
		OrganisationID: "d91afcdb-62d2-4185-b23d-71c98eaab812",
		Version:        0,
		CreatedOn:      time.Date(2020, time.November, 11, 10, 40, 44, 709000000, time.UTC),
		ModifiedOn:     time.Date(2020, time.November, 11, 10, 40, 44, 709000000, time.UTC),
		Attributes: &AccountAttributes{
			AccountClassification:       "Personal",
			AccountNumber:               "10000004",
			AlternativeBankAccountNames: nil,
			BankID:                      "400302",
			BankIDCode:                  "GBDSC",
			BaseCurrency:                "GBP",
			Bic:                         "NWBKGB42",
			Country:                     "GB",
			CustomerID:                  "234",
			Iban:                        "GB28NWBK40030212764204",
		}}

	if !reflect.DeepEqual(account, want) {
		t.Errorf("Account.List returned %+v, want %+v", account, want)
	}
}

func TestActiveService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/v1/organisation/accounts/ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Account.Delete(context.Background(), "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", 0)
	if err != nil {
		t.Errorf("Account.Delete returned error: %v", err)
	}

	if _, err := client.Account.Delete(context.Background(), "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", 0); err != nil {
		t.Errorf("Account.Delete returned error: %v", err)
	}
}
