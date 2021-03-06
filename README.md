# form3 client #

form3 is a Go client library for accessing the Form3 API.

## Usage ##
```go
package main

import "github.com/vslovik/form3/form3"	
```

Construct a new Form3 client, then use the account service on the client to
access account part of the Form3 API. For example:

```go
client := form3.NewClient(nil)

// list all accounts
accounts, _, _, err := client.Account.List(context.Background(), nil)
```

The services of a client can divide the API into logical chunks and correspond to
the structure of the Form3 API documentation at
https://api-docs.form3.tech/api.html.

Using the [context](https://godoc.org/context) package, one can easily
pass cancellation signals and deadlines to various services of the client for
handling a request. In case there is no context available, then `context.Background()`
can be used as a starting point.

### Creating Resources ###

All structs for Form3 resources use pointer values for all non-repeated fields.
This allows distinguishing between unset fields and those set to a zero-value.
For example:

```go
// create a new account
client := form3.NewClient(nil)

id := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
organizationId := "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"

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

_, _, _, err := client.Account.Create(context.Background(), id, organizationId, attr)
```

### Pagination ###

A request for resource collection (accounts) supports pagination. Pagination options are described in the
`form3.ListOptions` struct and passed to the list methods directly. 

```go
// get all accounts page by page
client := form3.NewClient(nil)

var allAccounts []*form3.Account
opt := &form3.ListOptions{PerPage: 2}
i := 0
for {
    opt.Page = i
    accounts, _, _, err := client.Account.List(context.Background(), opt)
    if len(accounts) == 0 {
        fmt.Printf("Account.List returned no accounts\n")
        break
    }
    if err != nil {
        fmt.Printf("Account.List returned error: %v\n", err)
    }
    allAccounts = append(allAccounts, accounts...)
    i += 1
}
```
## Tests ##

#### To run all tests in the form3 package: integration `integration_test.go` and unit tests `operations_test.go`, run

    $ cd interview-accountapi
    $ go test -v
    
with coverage:

    $ cd interview-accountapi
    $ go test -v --cover
    $ go test -coverprofile=coverage.out
    $ go tool cover -html=coverage.out
    
in docker-container:
   
    $ cd interview-accountapi
    $ docker-compose up
    
NOTE: BaseUrl is set in `client.go`:

```go
package form3

const defaultBaseURL = "http://localhost:8080/"
```
It should be `http://localhost:8080/` when running test on the host machine and `http://accountapi:8080/` to run them 
in docker container

#### To run separated integration tests `tests/integration/accounts_test.go`
Run fake form3 api with docker-compose, it will be accessible on http://localhost:8080
In form3/client.go set defaultBaseURL to http://loacalhost:8080/, then

    $ cd interview-accountapi/tests
    $ go get -t github.com/vslovik/form3/form3
    $ go test -v -tags=integration ./integration

See [README](interview-accountapi/tests/README.md).
    
#### To run Form3 SDK usage examples
Run fake form3 api with docker-compose, it will be accessible on http://localhost:8080
In form3/client.go set defaultBaseURL to http://loacalhost:8080/, then

    $ cd interview-accountapi/examples/accounts
    $ go get -t github.com/vslovik/form3/form3
    $ go run main.go

## Author

Valeriya Slovikovskaya, vslovik@gmail.com.

![test screen 1](https://github.com/vslovik/form3/blob/main/test_screens_1.png)
![test screen 2](https://github.com/vslovik/form3/blob/main/test_screens_2.png)
![test screen 3](https://github.com/vslovik/form3/blob/main/test_screens_3.png)
![test screen 4](https://github.com/vslovik/form3/blob/main/test_screens_4.png)


