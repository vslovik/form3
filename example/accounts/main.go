package main

import (
	"context"
	"fmt"

	"github.com/vslovik/form3"
)

const id = "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

var client = form3.NewClient(nil)

func createAccount() {
	operationId := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"

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

	_, _, _, err := client.Account.Create(context.Background(), id, operationId, attr)
	if err != nil {
		fmt.Printf("Account.Create returned error: %v", err)
	}

	// check again and verify exists
	acc, _, _, err := client.Account.Fetch(context.Background(), id)
	if err != nil {
		fmt.Printf("Account.Fetch returned error: %v", err)
	}
	if acc == nil {
		fmt.Printf("Account %v does not exists after creation.", id)
	}
}

func deleteAccount() {
	_, err := client.Account.Delete(context.Background(), id, 0)
	if err != nil {
		fmt.Printf("Account.Delete returned error: %v", err)
	}

	// check again and verify not exists
	acc, _, _, e := client.Account.Fetch(context.Background(), id)
	if e != nil {
		fmt.Printf("Account.Fetch returned error: %v", e)
	}
	if acc != nil {
		fmt.Printf("Still exists %v after deleting.", id)
	}
}

func test() {
	accounts, _, _, err := client.Account.List(context.Background(), nil)
	if err != nil {
		fmt.Printf("Account.List returned error: %v", err)
	}

	if len(accounts) == 0 {
		fmt.Printf("Account.List returned no accounts")
	}

	acc, _, _, err := client.Account.Fetch(context.Background(), id)
	if err != nil {
		fmt.Printf("Account.Fetch returned error: %v", err)
	}

	if acc != nil { // If already exists, delete then recreate account.
		deleteAccount()
		createAccount()
	} else { // Otherwise, create account and then delete it.
		createAccount()
		deleteAccount()
	}
}

func main() {
	test()
}
