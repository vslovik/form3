package integration

import (
	"context"
	"fmt"
	"github.com/vslovik/form3/form3"
	"os/exec"
	"strings"
	"testing"
)

var client = form3.NewClient(nil)

func uuid() (string, error) {
	cmd := exec.Command("uuidgen")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return string(stdout), nil
}

func createAccount(t *testing.T, id string, check bool) {
	uid, err := uuid()
	if err != nil {
		t.Fatalf("uuid generation error\n")
	}
	operationId := strings.TrimSuffix(string(uid), "\n")

	attr := &form3.AccountCreateRequestAttributes{
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
		uid, err := uuid()
		if err != nil {
			t.Fatalf("uuid generation error\n")
		}
		id := strings.TrimSuffix(string(uid), "\n")
		fmt.Printf("%v: Creating account %v...\n", i, id)
		createAccount(t, id, i == 0)
	}
}

func getPage(t *testing.T, page int, opt *form3.ListOptions) []*form3.Account {
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
func getAllPages(t *testing.T, perPage int) []*form3.Account {
	var allAccounts []*form3.Account
	opt := &form3.ListOptions{PerPage: perPage}
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

func deleteAll(t *testing.T, accounts []*form3.Account) {
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
