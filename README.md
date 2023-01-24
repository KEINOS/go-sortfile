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

    // Let the library auto detect the best way to sort the file.
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

func ExampleFromPathFunc() {
    exitOnError := func(err error) {
        if err != nil {
            log.Fatal(err)
        }
    }

    // Input and output file paths
    pathFileIn := filepath.Join("testdata", "sorted_chunks", "input_shuffled.txt")
    pathFileOut := filepath.Join(os.TempDir(), "pkg-sortfile_example_from_path.txt")

    // Clean up the output file after the test
    defer func() {
        exitOnError(os.Remove(pathFileOut))
    }()

    // User defined sort function. If isLess is nil, it will use the default sort
    // function which is equivalent to `sortfile.FromPath`.
    isLess := func(a, b string) bool {
        return a > b // reverse sort
    }

    // Let the library auto detect the best way to sort the file.
    forceExternalSort := false // auto detect

    // Sort file in-memory or external sort.
    err := sortfile.FromPathFunc(pathFileIn, pathFileOut, forceExternalSort, isLess)
    exitOnError(err)

    // Print the result
    data, err := os.ReadFile(pathFileOut)
    exitOnError(err)

    fmt.Println(string(data))
    // Output:
    // Zoe
    // Walter
    // Victor
    // Trudy
    // Trent
    // Steve
    // Peggy
    // Pat
    // Oscar
    // Matilda
    // Marvin
    // Mallory
    // Mallet
    // Justin
    // Ivan
    // Isaac
    // Frank
    // Eve
    // Ellen
    // Dave
    // Charlie
    // Carol
    // Bob
    // Alice
}
```

### Speed

Even a [simple implementation](./cmd/sortfile) is much faster than the ordinary `sort` command in linux/unix.
Though, we beleive it can be improved further.

```shellsession
$ # Around 1 GB of randomly shuffled data
$ ls -lah shuffled_huge.txt
-rw-r--r--  1 keinos  staff   985M  1 12 22:29 shuffled_huge.txt

$ # Ordinary sort command of linux/unix
$ time sort shuffled_huge.txt -o out_sort.txt
real    5m35.706s
user    11m52.320s
sys     0m40.690s

$ # Our sortfile command
$ time sortfile shuffled_huge.txt out_sortfile.txt
real    0m43.294s
user    0m36.283s
sys     0m5.751s

$ # Compare the result (no diff)
$ diff out_sort.txt out_sortfile.txt
$
```

## Contribute

- [PullRequest](https://github.com/KEINOS/go-sortfile/pulls)
  - Branch to PR: `main`
  - Any contribution for the better, faster, stronger implementation is welcome!
- [Issues](https://github.com/KEINOS/go-sortfile/issues)
  - Bug/vulnerability report: Please attach a reproducible simple example or a link to reference. It will help us alot to fix the issue faster.
  - Feature request: Please describe the feature you want to add and usecase. Although, we recommend to PR the feature, since it's more prioritized.
- [Help wanted](https://github.com/KEINOS/go-sortfile/issues)
