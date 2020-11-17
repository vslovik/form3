package form3

import (
	"context"
	"fmt"
	"time"
)

type AccountService service

type AccountAttributes struct {
	AccountClassification       string      `json:"account_classification"`
	AccountNumber               string      `json:"account_number"`
	AlternativeBankAccountNames interface{} `json:"alternative_bank_account_names"`
	BankID                      string      `json:"bank_id"`
	BankIDCode                  string      `json:"bank_id_code"`
	BaseCurrency                string      `json:"base_currency"`
	Bic                         string      `json:"bic"`
	Country                     string      `json:"country"`
	CustomerID                  string      `json:"customer_id"`
	Iban                        string      `json:"iban"`
}

type Account struct {
	Attributes     *AccountAttributes `json:"attributes"`
	CreatedOn      time.Time          `json:"created_on"`
	ID             string             `json:"id"`
	ModifiedOn     time.Time          `json:"modified_on"`
	OrganisationID string             `json:"organisation_id"`
	Type           string             `json:"type"`
	Version        int                `json:"version"`
}

type AccountListLinks struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Self  string `json:"self"`
}

type AccountCreateLinks struct {
	Self string `json:"self"`
}

type AccountFetchLinks struct {
	Self string `json:"self"`
}

type AccountFetchResponse struct {
	Data  *Account           `json:"data"`
	Links *AccountFetchLinks `json:"links"`
}

type AccountListResponse struct {
	Data  []*Account        `json:"data"`
	Links *AccountListLinks `json:"links"`
}

type AccountCreateResponse struct {
	Data  *Account            `json:"data"`
	Links *AccountCreateLinks `json:"links"`
}

// Request
type AccountCreateRequestAttributes struct {
	AccountClassification   string `json:"account_classification"`
	AccountNumber           string `json:"account_number"`
	BankID                  string `json:"bank_id"`
	BankIDCode              string `json:"bank_id_code"`
	BaseCurrency            string `json:"base_currency"`
	Bic                     string `json:"bic"`
	Country                 string `json:"country"`
	CustomerID              string `json:"customer_id"`
	Iban                    string `json:"iban"`
	JointAccount            bool   `json:"joint_account"`
	Switched                string `json:"switched"`
	SecondaryIdentification string `json:"secondary_identification"`
	AccountMatchingOptOut   bool   `json:"account_matching_opt_out"`
	AlternativeNames        bool   `json:"alternative_names"`
}

type AccountCreateRequestData struct {
	Attributes     *AccountCreateRequestAttributes `json:"attributes"`
	OrganisationID string                          `json:"organisation_id"`
	ID             string                          `json:"id"`
	Type           string                          `json:"type"`
}

type AccountCreateRequest struct {
	Data *AccountCreateRequestData `json:"data"`
}

/* ToDo
0. !!! Validation: Country is mandatory, depending on the country bank_id and bic are mandatory
1. !!! do fake first. If no account number or IBAN is provided, Form3 generates a valid account number (see below). If supported by the country, an IBAN is also generated. https://github.com/digitorus/iso7064
2. If an account number is provided but the IBAN is empty, Form3 generates an IBAN if supported by the country.
3. If only an IBAN is provided, the account number will be left empty.
4. !!! Validate country (belongs to list) https://www.iso.org/obp/ui/#search https://www.iso.org/iso-3166-country-codes.html
5. !!! Validate currency (belongs to list) list_one.xls
*/
// Form3 API docs: https://api-docs.form3.tech/api.html#organisation-accounts-create
func (s *AccountService) Create(ctx context.Context, id string, operationId string, attributes *AccountCreateRequestAttributes) (*Account, *AccountCreateLinks, *Response, error) {
	req, err := s.client.NewRequest("POST", "/v1/organisation/accounts",
		&AccountCreateRequest{&AccountCreateRequestData{
			attributes,
			operationId,
			id,
			"accounts"}})
	if err != nil {
		return nil, nil, nil, err
	}

	var r *AccountCreateResponse
	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, nil, resp, err
	}

	account := r.Data
	links := r.Links

	return account, links, resp, nil
}

/*
ToDo !!!: handle the situation when some fields in response are optional
*/
// Form3 API docs: https://api-docs.form3.tech/api.html#organisation-accounts-fetch
func (s *AccountService) Fetch(ctx context.Context, id string) (*Account, *AccountFetchLinks, *Response, error) {
	u := fmt.Sprintf("/v1/organisation/accounts/%s", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var r *AccountFetchResponse
	resp, err := s.client.Do(ctx, req, &r)
	if err != nil {
		return nil, nil, resp, err
	}

	account := r.Data
	links := r.Links

	return account, links, resp, nil
}

/*
ToDo!!!: support paging (and do not support filter)
*/
// Form3 API docs: https://api-docs.form3.tech/api.html#organisation-accounts-list
func (s *AccountService) List(ctx context.Context, opts *ListOptions) ([]*Account, *AccountListLinks, *Response, error) {
	var u string
	u = "/v1/organisation/accounts"
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var apiResp *AccountListResponse
	resp, err := s.client.Do(ctx, req, &apiResp)
	if err != nil {
		return nil, nil, resp, err
	}

	accounts := apiResp.Data
	links := apiResp.Links

	return accounts, links, resp, nil
}

/*
ToDo:
0. !!! Validate: account_id and number are mandatory
1. !!!Errors 204, 404, 409
*/
// Form3 API docs: https://api-docs.form3.tech/api.html#organisation-accounts-delete
func (s *AccountService) Delete(ctx context.Context, id string, version int) (*Response, error) {
	u := fmt.Sprintf("/v1/organisation/accounts/%s?version=%d", id, version)
	req, err := s.client.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(ctx, req, nil)
}
