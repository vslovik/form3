form3 tests
===============

This directory contains additional test suites beyond the unit tests already in
[../form3](..). Whereas the unit tests run very quickly (since they
don't make any network calls), the tests in this directory are only run manually.

The test package is:

integration
-----------

Run tests using:

    go test -v -tags=integration ./integration