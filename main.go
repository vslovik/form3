package main

import (
	"./form3"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

const (
	id          = "d91afcdb-62d2-4185-b23d-71c98eaab815"
	operationId = "d91afcdb-62d2-4185-b23d-71c98eaab812"
)

var attr = &form3.AccountCreateRequestAttributes{
	BankID:                  "400300",
	BankIDCode:              "GBDSC",
	BaseCurrency:            "GBP",
	Bic:                     "NWBKGB22",
	Country:                 "GB",
	AccountNumber:           "10000004",
	CustomerID:              "234",
	Iban:                    "GB28NWBK40030212764204",
	AccountClassification:   "Personal",
	JointAccount:            true,
	Switched:                "X",
	SecondaryIdentification: "X",
	AccountMatchingOptOut:   false,
	AlternativeNames:        false,
}

func resetAttr() {
	attr = &form3.AccountCreateRequestAttributes{
		BankID:                  "400300",
		BankIDCode:              "GBDSC",
		BaseCurrency:            "GBP",
		Bic:                     "NWBKGB22",
		Country:                 "GB",
		AccountNumber:           "10000004",
		CustomerID:              "234",
		Iban:                    "GB28NWBK40030212764204",
		AccountClassification:   "Personal",
		JointAccount:            true,
		Switched:                "X",
		SecondaryIdentification: "X",
		AccountMatchingOptOut:   false,
		AlternativeNames:        false,
	}
}

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

func createAccount(id string, check bool) {
	uid, err := uuid()
	if err != nil {
		fmt.Printf("uuid generation error\n")
		return
	}
	operationId := strings.TrimSuffix(string(uid), "\n")

	acc, _, _, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err != nil {
		panic(fmt.Sprintf("Account.Create returned error: %v\n", err))
	}
	if acc == nil {
		panic(fmt.Sprintf("Account %v does not exists after creation.\n", id))
	}
	fmt.Printf("OK\n")
	if check {
		fmt.Printf("Fetching account %v, checking all properties are correctly set...\n", id)
		acc, _, _, err := client.Account.Fetch(context.Background(), id)
		if err != nil {
			panic(fmt.Sprintf("Account.Fetch returned error: %v\n", err))
		}
		if acc.OrganisationID != operationId {
			panic(fmt.Sprintf("Invalid account OrganisationID: %v\n", acc.OrganisationID))
		}
		if acc.Type != "accounts" {
			panic(fmt.Sprintf("Invalid account Type: %v\n", acc.Type))
		}
		if acc.Attributes.BankID != "400300" {
			panic(fmt.Sprintf("Invalid account BankID: %v\n", acc.Attributes.BankID))
		}
		if acc.Attributes.BankIDCode != "GBDSC" {
			panic(fmt.Sprintf("Invalid account BankIDCode: %v\n", acc.Attributes.BankIDCode))
		}
		if acc.Attributes.BaseCurrency != "GBP" {
			panic(fmt.Sprintf("Invalid account BaseCurrency: %v\n", acc.Attributes.BaseCurrency))
		}
		if acc.Attributes.Bic != "NWBKGB22" {
			panic(fmt.Sprintf("Invalid account Bic: %v\n", acc.Attributes.Bic))
		}
		if acc.Attributes.Country != "GB" {
			panic(fmt.Sprintf("Invalid account Country: %v\n", acc.Attributes.Country))
		}
		if acc.Attributes.AccountNumber != "10000004" {
			panic(fmt.Sprintf("Invalid account AccountNumber: %v\n", acc.Attributes.AccountNumber))
		}
		if acc.Attributes.CustomerID != "234" {
			panic(fmt.Sprintf("Invalid account CustomerID: %v\n", acc.Attributes.CustomerID))
		}
		if acc.Attributes.Iban != "GB28NWBK40030212764204" {
			panic(fmt.Sprintf("Invalid account Iban: %v\n", acc.Attributes.Iban))
		}
		if acc.Attributes.AccountClassification != "Personal" {
			panic(fmt.Sprintf("Invalid account AccountClassification: %v\n", acc.Attributes.AccountClassification))
		}
		fmt.Printf("OK\n")
	}
}

func deleteAccount(id string) {
	acc, _, _, err := client.Account.Fetch(context.Background(), id)
	if err != nil {
		panic(fmt.Sprintf("Account.Fetch returned error: %v\n", err))
	}

	_, err = client.Account.Delete(context.Background(), id, 0)
	if err != nil {
		panic(fmt.Sprintf("Account.Delete returned error: %v\n", err))
	}

	// check again and verify not exists
	acc, _, resp, err := client.Account.Fetch(context.Background(), id)
	if err == nil {
		panic(fmt.Sprintf("Account.Fetch does not returned error on not existent account (%v) fetch\n", id))
	}
	if resp != nil && resp.StatusCode != 404 {
		panic(fmt.Sprintf("Account.Fetch does not returned 404 staus code on not existent account (%v) fetch\n", id))
	}
	if acc != nil {
		panic(fmt.Sprintf("Still exists %v after deleting.\n", id))
	}
	fmt.Printf("OK\n")
}

func createAccountBunch(number int) {
	for i := 0; i < number; i++ {
		uid, err := uuid()
		if err != nil {
			fmt.Printf("uuid generation error\n")
			return
		}
		id := strings.TrimSuffix(string(uid), "\n")
		fmt.Printf("%v: Creating account %v...\n", i, id)
		createAccount(id, i == 0)
	}
}

func getPage(page int, opt *form3.ListOptions) ([]*form3.Account, error) {
	opt.Page = page
	accounts, _, _, err := client.Account.List(context.Background(), opt)
	for _, elem := range accounts {
		fmt.Printf("Got account %s \n", elem.ID)
	}
	if err != nil {
		panic(fmt.Sprintf("Account.List returned error: %v\n", err))
	}

	return accounts, nil
}

// get all pages of results
func getAllPages(perPage int) ([]*form3.Account, error) {
	var allAccounts []*form3.Account
	opt := &form3.ListOptions{PerPage: perPage}
	i := 0
	for {
		fmt.Printf("Retrieving Page %v of %v accounts...\n", i, opt.PerPage)
		accounts, err := getPage(i, opt)
		if err != nil {
			return nil, err
		}
		if len(accounts) == 0 {
			fmt.Printf("Account.List returned no accounts\n")
			break
		}
		fmt.Printf("Account.List returned %v accounts\n", len(accounts))
		allAccounts = append(allAccounts, accounts...)
		i += 1
	}
	fmt.Printf("Account.List retrieved  %v accounts\n", len(allAccounts))
	return allAccounts, nil
}

func deleteAll(accounts []*form3.Account) {
	for i, elem := range accounts {
		fmt.Printf("%v: Deleting account %s...\n", i+1, elem.ID)
		deleteAccount(elem.ID)
	}
}

func positiveTestCases() {
	n := 10
	fmt.Printf("I want to create %v accounts...\n", n)
	createAccountBunch(n)

	perPage := 2
	fmt.Printf("I want to retrive all accounts page by page, %v account per page\n", perPage)
	allAccounts, _ := getAllPages(perPage)

	fmt.Printf("I want to retrive all accounts in one request\n")
	accounts, _, _, err := client.Account.List(context.Background(), nil)
	if err != nil {
		panic(fmt.Sprintf("Account.List returned error: %v\n", err))
	}
	fmt.Printf("%v accounts retrieved\n", len(accounts))
	fmt.Printf("I check that the number of accounts retrived in both operations is the same\n")
	if len(accounts) != len(allAccounts) {
		panic(fmt.Sprintf("Wrong number of accounts retrieved\n"))
	} else {
		fmt.Printf("OK\n")
	}
	fmt.Printf("I want to delete all accounts\n")
	deleteAll(accounts)

	fmt.Printf("I want to check that there is no accounts left\n")
	accounts, _, _, err = client.Account.List(context.Background(), nil)
	if err != nil {
		panic(fmt.Sprintf("Account.List returned error: %v\n", err))
	}
	if len(accounts) > 0 {
		panic(fmt.Sprintf("%v accounts retrieved\n", len(accounts)))
	}
	fmt.Printf("OK\n")
}

func testCreateBadRequestEmptyAccountId() {
	fmt.Print("Test: Create: bad request: empty account id\n")
	_, _, resp, err := client.Account.Create(context.Background(), "", operationId, attr)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
}

func testCreateBadRequestEmptyOperationId() {
	fmt.Print("Test: Create: bad request: empty operation id\n")
	_, _, resp, err := client.Account.Create(context.Background(), id, "", attr)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
}

func testCreateBadRequestEmptyAttributesList() {
	fmt.Print("Test: Create: bad request: empty attributes list\n")
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, &form3.AccountCreateRequestAttributes{})
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
}

func testCreateBadRequestInvalidId() {
	fmt.Print("Test: Create: bad request: invalid id\n")
	_, _, resp, err := client.Account.Create(context.Background(), "x", operationId, attr)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
}

func testCreateBadRequestInvalidOperationId() {
	fmt.Print("Test: Create: bad request: operation id\n")
	_, _, resp, err := client.Account.Create(context.Background(), id, "x", attr)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
}

func testCreateBadRequestInvalidCountry() {
	fmt.Print("Test: Create: bad request: invalid country\n")
	attr.Country = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
	resetAttr()
}

func testCreateBadRequestInvalidBaseCurrency() {
	fmt.Print("Test: Create: bad request: invalid base currency\n")
	attr.BaseCurrency = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
	resetAttr()
}

func testCreateBadRequestInvalidIban() {
	fmt.Print("Test: Create: bad request: invalid IBAN\n")
	attr.Iban = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
	resetAttr()
}

func testCreateBadRequestInvalidBic() {
	fmt.Print("Test: Create: bad request: invalid BIC\n")
	attr.Bic = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
	resetAttr()
}

func testCreateBadRequestInvalidAccountClassification() {
	fmt.Print("Test: Create: bad request: invalid account classification\n")
	attr.AccountClassification = "x"
	_, _, resp, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
	resetAttr()
}

func testFetchNotExistentAccount() {
	fmt.Print("Test: Fetch: bad request: account does not exist\n")
	_, _, resp, err := client.Account.Fetch(context.Background(), "d91afcdb-xxxx-4185-b23d-71c98eaab815")
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
	resetAttr()
}

func testDeleteNotExistentAccount() {
	fmt.Print("Test: Delete: bad request: account does not exist\n")
	resp, err := client.Account.Delete(context.Background(), "d91afcdb-xxxx-4185-b23d-71c98eaab815", 0)
	if err == nil {
		panic("Test failed\n")
	}
	if resp == nil {
		panic("Test failed\n")
	}

	if resp.StatusCode != 400 {
		panic("Test failed\n")
	}
	fmt.Print("OK\n")
	resetAttr()
}

func main() {
	positiveTestCases()
	testCreateBadRequestEmptyAccountId()
	testCreateBadRequestEmptyOperationId()
	testCreateBadRequestEmptyAttributesList()
	testCreateBadRequestInvalidId()
	testCreateBadRequestInvalidOperationId()
	testCreateBadRequestInvalidCountry()
	testCreateBadRequestInvalidBaseCurrency()
	testCreateBadRequestInvalidIban()
	testCreateBadRequestInvalidBic()
	testCreateBadRequestInvalidAccountClassification()
	testFetchNotExistentAccount()
	testDeleteNotExistentAccount()
}
