# Test Data Directory

This directory contains test data for the sortfile package.

To generate the large test data file for benchmark, run the following command:

    go generate ./...

**Note:** The large test data file is not included in the repository. It is
generated on demand by the `go generate` command. Don't forget to remove the
generated file after benchmarking. It will take up a lot of space (around 1 GB).
