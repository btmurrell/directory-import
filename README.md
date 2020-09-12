# directory-import

This is a utility program written in [Go](https://golang.org/) to convert a large CSV file of all parents and students into multiple CSV files, grouped by grade and classroom, suitable for import into a my-pta website.  It was written as an exploratory project to learn the Go programming language.  It is effective at what it does.

## Building
For OSX, there is a strange, annoying error I don't understand that causes a built executable to fail with this ugly error:

````
# command-line-arguments
fatal error: runtime: bsdthread_register error
runtime stack:
runtime.throw(0x4576a0, 0x21)
/usr/local/go/src/runtime/panic.go:527 +0x90 fp=0x7fff5fbff680 sp=0x7fff5fbff668
runtime.goenvs()
/usr/local/go/src/runtime/os1_darwin.go:73 +0x8d fp=0x7fff5fbff6a8 sp=0x7fff5fbff680
runtime.schedinit()
/usr/local/go/src/runtime/proc1.go:60 +0x83 fp=0x7fff5fbff6f0 sp=0x7fff5fbff6a8
runtime.rt0_go(0x7fff5fbff728, 0x10, 0x7fff5fbff728, 0x0, 0x0, 0x10, 0x7fff5fbff910, 0x7fff5fbff93c, 0x7fff5fbff93f, 0x7fff5fbff99b, ...)
/usr/local/go/src/runtime/asm_amd64.s:109 +0x132 fp=0x7fff5fbff6f8 sp=0x7fff5fbff6f0
````

To avoid that error, build with this `-ldflags`. ([attribution](https://github.com/golang/go/issues/8801#issuecomment-66460009))
```
 go build -ldflags="-linkmode=external"
```



## Notes on running

A data file may be found in private repo [data.csv](https://github.com/btmurrell/data-repo/blob/master/directory-import/data.csv)

Run it this way:
```
directory-import filename.csv
```
where `filename.csv` is the CSV file that you want to process.  It will create a set of CSV files broken out by grade+classroom in a folder under you current folder titled `csv-output`.

You may optionally turn on logging output with the `-l` flag, possible values are:

 * `p` (panic) highest level (no logging output unless there's a Go runtime error)
 * `f` (fatal) will only report fatal errors (those that cause the app to crash, e.g. can't find file)
 * `e` (error) errors processing records of the file, but program continues (e.g. no email; record discarded)
 * `w` (warn) warnings while processing records, but record is accepted
 * `i` (info) information, such as processing summary
 * `d` (debug) lowest level detail logging output


### Configuration
By default, the input file's column order is expected to be:

TODO: FILL THIS OUT

FUTURE: you can specify the columns flag to indicate which spreadsheet column each

```
 -columns "parentName:J,studentName:A,parentEmail:N,parentEmailAlt:P,grade:G,room:B"
```
