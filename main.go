package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Found   bool
	Address string
}

var (
	inputAccount string
	inputFile  string
	outputFile string
	domain     string
	stdin      bool
	validChars bool
	threads    int
	client     *http.Client
	headers    map[string]string
)

const (
        InfoColor    = "\033[1;34m%s\033[0m"
        NoticeColor  = "\033[1;36m%s\033[0m"
        WarningColor = "\033[1;33m%s\033[0m"
        ErrorColor   = "\033[1;31m%s\033[0m"
        DebugColor   = "\033[0;36m%s\033[0m"
)

func init() {
	flag.StringVar(&inputFile, "I", "", "File of accounts to test")
	flag.StringVar(&inputAccount, "i", "", "accounts to test")
	flag.StringVar(&outputFile, "o", "", "Output file (default: Stdout)")
	flag.StringVar(&domain, "d", "", "Append domain to every address (empty to no append)")
	flag.BoolVar(&stdin, "stdin", false, "Read accounts from stdin")
	flag.BoolVar(&validChars, "r", false, "Remove gmail address' invalid chars")
	flag.IntVar(&threads, "t", 10, "Number of threads")
	flag.Parse()

	if inputFile == "" && inputAccount == "" && !stdin {
		flag.Usage()
		os.Exit(1)
	}

	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			// MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}

	headers = map[string]string{
		"User-Agent":      `Mozilla/5.0 (Windows NT 6.1; rv:61.0) Gecko/20100101 Firefox/61.0`,
		"Accept-Language": `en-US,en;q=0.5`,
	}
}

// TestAddress checks if a given address is valid using the glitch described here: https://blog.0day.rocks/abusing-gmail-to-get-previously-unlisted-e-mail-addresses-41544b62b2
func TestAddress(addr string, resChan chan<- Result) {
	URL := fmt.Sprintf("https://mail.google.com/mail/gxlu?email=%s", url.QueryEscape(addr))
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return
	}

	// Add headers
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()

	found := len(resp.Cookies()) > 0
	//fmt.Print(found,"\n")
	//fmt.Print(addr,"\n")
	resChan <- Result{found, addr}
	//fmt.Print(addr,"\n")
	return
}

func main() {

	addrChan := make(chan string, threads)
	resultsChan := make(chan Result)

	// Group to wait for all threads (routines) to finish
	threadsG := new(sync.WaitGroup)

	var input *os.File
	if stdin {
		input = os.Stdin
		inputFile = "stdin"
	} else if inputFile != "" {
		f, err := os.Open(inputFile)
		if err != nil {
			fmt.Printf("[!] Error opening file '%s'\n", inputFile)
			return
		}
		input = f
		defer f.Close()
	} else {
		inputFile = "single"

	}

	// if an output file is provided redirect output to file
	out, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE, os.ModeAppend)
	if err != nil {
		// if no file provided redirect to stdout
		out = os.Stdout
	}
	defer out.Close()

	// TODO: Put some fancy ascii art here??
	fmt.Println("--- Starting bruteforce --")
	fmt.Printf("| Input:   %s\n", inputFile)
	fmt.Printf("| Threads: %d\n\n", threads)

	// Start all threads (routines)
	for i := 0; i < threads; i++ {
		go func() {
			for addr := range addrChan {
				// Append domain to address
				if domain != "" {
					addr += "@" + domain
				}

				if validChars {
					addr = RemoveInvalidChars(addr)
				}

				TestAddress(addr, resultsChan)
			}
			threadsG.Done()
		}()
		threadsG.Add(1)
	}

	// if not single then iterate through stdin or file
	if inputFile != "single" {
		scanner := bufio.NewScanner(input)
		scanner.Split(bufio.ScanLines)

		go func() {
			for scanner.Scan() {
				addr := strings.TrimSpace(scanner.Text())
				// Skip comments and empty lines
				if !strings.HasPrefix(addr, "#") && addr != "" {
					addrChan <- addr
				}
			}

			close(addrChan)
			threadsG.Wait()
			close(resultsChan)
		}()
	// else single email test
	} else if inputAccount != "" {
		go func() {
			addrChan <- inputAccount
			close(addrChan)
			threadsG.Wait()
			close(resultsChan)
	}()
	}

	tested, found := 0, 0
	for result := range resultsChan {
		tested++
		if result.Found {
			found++
			if out == os.Stdout {
				// 'Flush' stdout
				fmt.Printf("%100s\r", "")
			}
			// print found emails in yellow
			fmt.Fprintln(out, fmt.Sprintf(WarningColor, result.Address))
		}

		fmt.Printf("[*] Tested: %d, Found: %d\r", tested, found)
	}
	fmt.Printf("[*] Tested: %d, Found: %d\n", tested, found)

}
