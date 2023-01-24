# `sortfile` command

`sortfile` is a simple command line tool to sort a file in-memory or external sort.

```shell
sortfile <input file> <output file>
```

It is much faster than the ordinary `sort` command in linux/unix. Though, we beleive it can be improved further.

```shellsession
$ # Around 1 GB of random data
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
