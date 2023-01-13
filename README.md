# go-sortfile

[go-sortfile](https://github.com/KEINOS/go-sortfile) is a simple Go library for sorting large files.

If the file size is smaller than the available memory, sort in-memory; if the file size is larger than the available memory, perform an external sort with a K-way merge sort.

## Usage

```go
go get "github.com/KEINOS/go-sortfile"
```

### Example

```go
import "github.com/KEINOS/go-sortfile/sortfile"

func ExampleFromPath() {
    // Input and output file paths
    pathFileIn := filepath.Join("path", "to", "large_file.txt")
    pathFileOut := filepath.Join("path", "to", "large_file-sorted.txt")

    forceExternalSort := false // false ==> auto detect sort method

    // Sort file in-memory or external sort.
    //
    // If the 3rd argument is false, the library will auto detect the best way
    // to sort the file. In-memory sorting or external sorting. If true it will
    // force to use external sort.
    err := sortfile.FromPath(pathFileIn, pathFileOut, forceExternalSort)
    if err != nil {
        log.Fatal(err)
    }

    // Print the result
    data, err := os.ReadFile(pathFileOut)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(data))
    // Output:
    // Alice
    // Bob
    // Carol
    // Charlie
    // Dave
    // Ellen
    // Eve
    // Frank
    // Isaac
    // Ivan
    // Justin
    // Mallet
    // Mallory
    // Marvin
    // Matilda
    // Oscar
    // Pat
    // Peggy
    // Steve
    // Trent
    // Trudy
    // Victor
    // Walter
    // Zoe
}
```

## Contribute

- [PullRequest](https://github.com/KEINOS/go-sortfile/pulls)
  - Branch to PR: `main`
  - Any contribution for the better, faster, stronger implementation is welcome!
- [Issues](https://github.com/KEINOS/go-sortfile/issues)
  - Bug/vulnerability report: Please attach a reproducible simple example or a link to reference. It will help us alot to fix the issue faster.
  - Feature request: Please describe the feature you want to add and usecase. Although, we recommend to PR the feature, since it's more prioritized.
- [Help wanted](https://github.com/KEINOS/go-sortfile/issues)
