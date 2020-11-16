package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/vslovik/form3/form3"
	"log"
	"strings"
)

var client = form3.NewClient(nil)

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

func createAccount(id string, check bool) {
	uid := uuid()
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
		log.Fatal(fmt.Sprintf("Account.Create returned error: %v\n", err))
	}
	if acc == nil {
		fmt.Printf("Account %v does not exists after creation.\n", id)
	}
	fmt.Printf("OK\n")
	if check {
		fmt.Printf("Fetching account %v, checking all properties are correctly set...\n", id)
		acc, _, _, err := client.Account.Fetch(context.Background(), id)
		if err != nil {
			log.Fatal(fmt.Sprintf("Account.Fetch returned error: %v\n", err))
		}
		if acc.OrganisationID != operationId {
			log.Fatal(fmt.Sprintf("Invalid account OrganisationID: %v\n", acc.OrganisationID))
		}
		if acc.Type != "accounts" {
			log.Fatal(fmt.Sprintf("Invalid account Type: %v\n", acc.Type))
		}
		if acc.Attributes.BankID != "400300" {
			log.Fatal(fmt.Sprintf("Invalid account BankID: %v\n", acc.Attributes.BankID))
		}
		if acc.Attributes.BankIDCode != "GBDSC" {
			log.Fatal(fmt.Sprintf("Invalid account BankIDCode: %v\n", acc.Attributes.BankIDCode))
		}
		if acc.Attributes.BaseCurrency != "GBP" {
			log.Fatal(fmt.Sprintf("Invalid account BaseCurrency: %v\n", acc.Attributes.BaseCurrency))
		}
		if acc.Attributes.Bic != "NWBKGB22" {
			log.Fatal(fmt.Sprintf("Invalid account Bic: %v\n", acc.Attributes.Bic))
		}
		if acc.Attributes.Country != "GB" {
			log.Fatal(fmt.Sprintf("Invalid account Country: %v\n", acc.Attributes.Country))
		}
		if acc.Attributes.AccountNumber != "10000004" {
			log.Fatal(fmt.Sprintf("Invalid account AccountNumber: %v\n", acc.Attributes.AccountNumber))
		}
		if acc.Attributes.CustomerID != "234" {
			log.Fatal(fmt.Sprintf("Invalid account CustomerID: %v\n", acc.Attributes.CustomerID))
		}
		if acc.Attributes.Iban != "GB28NWBK40030212764204" {
			log.Fatal(fmt.Sprintf("Invalid account Iban: %v\n", acc.Attributes.Iban))
		}
		if acc.Attributes.AccountClassification != "Personal" {
			log.Fatal(fmt.Sprintf("Invalid account AccountClassification: %v\n", acc.Attributes.AccountClassification))
		}
		fmt.Printf("OK\n")
	}
}

func deleteAccount(id string) {
	acc, _, _, err := client.Account.Fetch(context.Background(), id)
	if err != nil {
		log.Fatal(fmt.Sprintf("Account.Fetch returned error: %v\n", err))
	}

	_, err = client.Account.Delete(context.Background(), id, 0)
	if err != nil {
		log.Fatal(fmt.Sprintf("Account.Delete returned error: %v\n", err))
	}

	// check again and verify not exists
	acc, _, _, err = client.Account.Fetch(context.Background(), id)
	if err != nil {
		log.Fatal(fmt.Sprintf("Account.Fetch returned error: %v\n", err))
	}
	if acc != nil {
		log.Fatal(fmt.Sprintf("Still exists %v after deleting.\n", id))
	}
	fmt.Printf("OK\n")
}

func createAccountBunch(number int) {
	for i := 0; i < number; i++ {
		uid := uuid()
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
		log.Fatal(fmt.Sprintf("Account.List returned error: %v\n", err))
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

func main() {
	n := 10
	fmt.Printf("I want to create %v accounts...\n", n)
	createAccountBunch(n)

	perPage := 2
	fmt.Printf("I want to retrive all accounts page by page, %v account per page\n", perPage)
	allAccounts, _ := getAllPages(perPage)

	fmt.Printf("I want to retrive all accounts in one request\n")
	accounts, _, _, err := client.Account.List(context.Background(), nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Account.List returned error: %v\n", err))
	}
	fmt.Printf("%v accounts retrieved\n", len(accounts))
	fmt.Printf("I check that the number of accounts retrived in both operations is the same\n")
	if len(accounts) != len(allAccounts) {
		log.Fatal(fmt.Sprintf("Wrong number of accounts retrieved\n"))
	} else {
		fmt.Printf("OK")
	}
	fmt.Printf("I want to delete all accounts\n")
	deleteAll(accounts)

	fmt.Printf("I want to check that there is no accounts left\n")
	accounts, _, _, err = client.Account.List(context.Background(), nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Account.List returned error: %v\n", err))
	}
	if len(accounts) > 0 {
		log.Fatal(fmt.Sprintf("%v accounts retrieved\n", len(accounts)))
	}
	fmt.Printf("OK\n")
}
