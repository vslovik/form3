package form3

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"testing"
)

const (
	id          = "d91afcdb-62d2-4185-b23d-71c98eaab815"
	operationId = "d91afcdb-62d2-4185-b23d-71c98eaab812"
)

var client = NewClient(nil)

func uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func createAccount(t *testing.T, id string, check bool) {
	uid := uuid()

	operationId := strings.TrimSuffix(uid, "\n")

	attr := &AccountCreateRequestAttributes{
		BankID:                "400300",
		BankIDCode:            "GBDSC",
		BaseCurrency:          "GBP",
		Bic:                   "NWBKGB22",
		Country:               "GB",
		AccountNumber:         "10000004",
		CustomerID:            "234",
		Iban:                  "GB28NWBK40030212764204",
		AccountClassification: "Personal",
	}

	acc, _, _, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err != nil {
		t.Fatalf("Account.Create returned error: %v\n", err)
	}
	if acc == nil {
		t.Fatalf("Account %v does not exists after creation.\n", id)
	}
	fmt.Printf("OK\n")
	if check {
		fmt.Printf("Fetching account %v, checking all properties are correctly set...\n", id)
		acc, _, _, err := client.Account.Fetch(context.Background(), id)
		if err != nil {
			t.Fatalf("Account.Fetch returned error: %v\n", err)
		}
		if acc.OrganisationID != operationId {
			t.Fatalf("Invalid account OrganisationID: %v\n", acc.OrganisationID)
		}
		if acc.Type != "accounts" {
			t.Fatalf("Invalid account Type: %v\n", acc.Type)
		}
		if acc.Attributes.BankID != "400300" {
			t.Fatalf("Invalid account BankID: %v\n", acc.Attributes.BankID)
		}
		if acc.Attributes.BankIDCode != "GBDSC" {
			t.Fatalf("Invalid account BankIDCode: %v\n", acc.Attributes.BankIDCode)
		}
		if acc.Attributes.BaseCurrency != "GBP" {
			t.Fatalf("Invalid account BaseCurrency: %v\n", acc.Attributes.BaseCurrency)
		}
		if acc.Attributes.Bic != "NWBKGB22" {
			t.Fatalf("Invalid account Bic: %v\n", acc.Attributes.Bic)
		}
		if acc.Attributes.Country != "GB" {
			t.Fatalf("Invalid account Country: %v\n", acc.Attributes.Country)
		}
		if acc.Attributes.AccountNumber != "10000004" {
			t.Fatalf("Invalid account AccountNumber: %v\n", acc.Attributes.AccountNumber)
		}
		if acc.Attributes.CustomerID != "234" {
			t.Fatalf("Invalid account CustomerID: %v\n", acc.Attributes.CustomerID)
		}
		if acc.Attributes.Iban != "GB28NWBK40030212764204" {
			t.Fatalf("Invalid account Iban: %v\n", acc.Attributes.Iban)
		}
		if acc.Attributes.AccountClassification != "Personal" {
			t.Fatalf("Invalid account AccountClassification: %v\n", acc.Attributes.AccountClassification)
		}
		fmt.Printf("OK\n")
	}
}

func deleteAccount(t *testing.T, id string) {
	acc, _, _, err := client.Account.Fetch(context.Background(), id)
	if err != nil {
		t.Fatalf("Account.Fetch returned error: %v\n", err)
	}

	_, err = client.Account.Delete(context.Background(), id, 0)
	if err != nil {
		t.Fatalf("Account.Delete returned error: %v\n", err)
	}

	// check again and verify not exists
	acc, _, _, err = client.Account.Fetch(context.Background(), id)
	if err != nil {
		t.Fatalf("Account.Fetch returned error: %v\n", err)
	}
	if acc != nil {
		t.Fatalf("Still exists %v after deleting.\n", id)
	}
	fmt.Printf("OK\n")
}

func createAccountBunch(t *testing.T, number int) {
	for i := 0; i < number; i++ {
		uid := uuid()
		id := strings.TrimSuffix(uid, "\n")
		fmt.Printf("%v: Creating account %v...\n", i, id)
		createAccount(t, id, i == 0)
	}
}

func getPage(t *testing.T, page int, opt *ListOptions) []*Account {
	opt.Page = page
	accounts, _, _, err := client.Account.List(context.Background(), opt)
	for _, elem := range accounts {
		fmt.Printf("Got account %s \n", elem.ID)
	}
	if err != nil {
		t.Fatalf("Account.List returned error: %v\n", err)
	}
	return accounts
}

// get all pages of results
func getAllPages(t *testing.T, perPage int) []*Account {
	var allAccounts []*Account
	opt := &ListOptions{PerPage: perPage}
	i := 0
	for {
		fmt.Printf("Retrieving Page %v of %v accounts...\n", i, opt.PerPage)
		accounts := getPage(t, i, opt)
		if len(accounts) == 0 {
			fmt.Printf("Account.List returned no accounts\n")
			break
		}
		fmt.Printf("Account.List returned %v accounts\n", len(accounts))
		allAccounts = append(allAccounts, accounts...)
		i += 1
	}
	fmt.Printf("Account.List retrieved  %v accounts\n", len(allAccounts))
	return allAccounts
}

func deleteAll(t *testing.T, accounts []*Account) {
	for i, elem := range accounts {
		fmt.Printf("%v: Deleting account %s...\n", i+1, elem.ID)
		deleteAccount(t, elem.ID)
	}
}

func TestAccount_ListFetchCreateDelete(t *testing.T) {
	n := 10
	fmt.Printf("I want to create %v accounts...\n", n)
	createAccountBunch(t, n)

	perPage := 2
	fmt.Printf("I want to retrive all accounts page by page, %v account per page\n", perPage)
	allAccounts := getAllPages(t, perPage)

	fmt.Printf("I want to retrive all accounts in one request\n")
	accounts, _, _, err := client.Account.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Account.List returned error: %v\n", err)
	}
	fmt.Printf("%v accounts retrieved\n", len(accounts))
	fmt.Printf("I check that the number of accounts retrived in both operations is the same\n")
	if len(accounts) != len(allAccounts) {
		t.Fatalf("Wrong number of accounts retrieved\n")
	} else {
		fmt.Printf("OK")
	}
	fmt.Printf("I want to delete all accounts\n")
	deleteAll(t, accounts)

	fmt.Printf("I want to check that there is no accounts left\n")
	accounts, _, _, err = client.Account.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Account.List returned error: %v\n", err)
	}
	if len(accounts) > 0 {
		t.Fatalf("%v accounts retrieved\n", len(accounts))
	}
	fmt.Printf("OK\n")
}

func TestCreate_BadRequestEmptyAccountId(t *testing.T) {
	fmt.Print("Test: Create: bad request: empty account id\n")
	_, _, resp, err := client.Account.Create(context.Background(), "", operationId, attr)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestCreate_BadRequestEmptyOperationId(t *testing.T) {
	fmt.Print("Test: Create: bad request: empty operation id\n")
	_, _, resp, err := client.Account.Create(context.Background(), id, "", attr)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestCreate_BadRequestEmptyAttributesList(t *testing.T) {
	fmt.Print("Test: Create: bad request: empty attributes list\n")
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, &AccountCreateRequestAttributes{})
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestCreate_BadRequestInvalidId(t *testing.T) {
	fmt.Print("Test: Create: bad request: invalid id\n")
	_, _, resp, err := client.Account.Create(context.Background(), "x", operationId, attr)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestCreate_BadRequestInvalidOperationId(t *testing.T) {
	fmt.Print("Test: Create: bad request: operation id\n")
	_, _, resp, err := client.Account.Create(context.Background(), id, "x", attr)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestCreate_BadRequestInvalidCountry(t *testing.T) {
	fmt.Print("Test: Create: bad request: invalid country\n")
	attr.Country = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestCreate_BadRequestInvalidBaseCurrency(t *testing.T) {
	fmt.Print("Test: Create: bad request: invalid base currency\n")
	attr.BaseCurrency = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestCreate_BadRequestInvalidIban(t *testing.T) {
	fmt.Print("Test: Create: bad request: invalid IBAN\n")
	attr.Iban = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestCreate_BadRequestInvalidBic(t *testing.T) {
	fmt.Print("Test: Create: bad request: invalid BIC\n")
	attr.Bic = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestCreate_BadRequestInvalidAccountClassification(t *testing.T) {
	fmt.Print("Test: Create: bad request: invalid account classification\n")
	attr.AccountClassification = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestFetch_NotExistentAccount(t *testing.T) {
	fmt.Print("Test: Fetch: bad request: account does not exist\n")
	_, _, resp, err := client.Account.Fetch(context.Background(), "d91afcdb-xxxx-4185-b23d-71c98eaab815")
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}

func TestDelete_NotExistentAccount(t *testing.T) {
	fmt.Print("Test: Delete: bad request: account does not exist\n")
	resp, err := client.Account.Delete(context.Background(), "d91afcdb-xxxx-4185-b23d-71c98eaab815", 0)
	if err == nil {
		t.Fatalf("Test failed: no error\n")
	}
	if resp == nil {
		t.Fatalf("Test failed: response is nil\n")
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Test failed: response status code (%v) is not 400\n", resp.StatusCode)
	}
	fmt.Print("OK\n")
}
