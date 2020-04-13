<h4 align="center">Gmail Email Enumeration</h4>
<p align="center">
  <a href="https://twitter.com/sho_luv">
    <img src="https://img.shields.io/badge/Twitter-%40sho_luv-blue.svg">
  </a>
</p>

## Gmail Enum

I forked this project from https://github.com/H3LL0WORLD/Gmail-Enum and simply made some modifications to suite my needs. 
A fairly descent/fast go program to enumerate gmail accounts using a glitch by [@x0rz](https://twitter.com/x0rz) as described [here](https://blog.0day.rocks/abusing-gmail-to-get-previously-unlisted-e-mail-addresses-41544b62b2)

<img src="https://github.com/sho-luv/gmailprob/blob/master/img/gmailenum.gif" alt="gmailenum" />

### Requirements:
- [Golang](https://golang.org)

## How to Use
```
# Clone this repository:
$ git clone https://github.com/sho-luv/gmailprob

# Go into the repository:
$ cd gmailprob

# build the app
$ go build 

# run the app
$ ./emailprob 
Usage of ./emailprob:
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

