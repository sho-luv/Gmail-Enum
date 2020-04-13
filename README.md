## Gmail Enum
A fairly descent/fast go program to enumerate gmail accounts using a glitch by [@x0rz](https://twitter.com/x0rz) as described [here](https://blog.0day.rocks/abusing-gmail-to-get-previously-unlisted-e-mail-addresses-41544b62b2)

### Requirements:
- [Golang](https://golang.org)

### Usage:
```sh
$ go build

$ ./Gmail-Enum 
Usage of ./Gmail-Enum:
  -I string
    	File of accounts to test
  -d string
    	Append domain to every address (empty to no append)
  -i string
    	account to test
  -o string
    	Output file (default: Stdout)
  -r	Remove gmail address' invalid chars
  -stdin
    	Read accounts from stdin
  -t int
    	Number of threads (default 10)
```
