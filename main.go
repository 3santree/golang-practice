package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
)

func main() {

	// Flag handling, must specify the url
	var domain = flag.String("d", "", "Specify domain url")
	var output = flag.String("o", "", "output file location")
	flag.Parse()
	if *domain == "" {
		flag.PrintDefaults()
		return
	}

	// fetch json from API
	fmt.Println("Target:", *domain)
	api := "https://crt.sh/?q=" + *domain + "&output=json"
	resp, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}
	resp.Header.Set("User-Agent", "Mozilla/5.0")
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// Get subdomain from json -> slice -> make slice uniq
	// gjson play ground: https://gjson.dev/
	value := gjson.Get(string(body), "#.common_name")
	var subList []string
	for i := 0; i < int(value.Get("#").Num); i++ {
		path := strconv.Itoa(i)
		subList = append(subList, value.Get(path).String())
	}
	subList = removeDuplicates(subList)

	// get ip from url, conbine to a map
	subMap := make(map[string]string)
	for _, url := range subList {
		ips, err := net.LookupIP(url)
		if err != nil {
			continue
		}

		fmt.Println(ips[len(ips)-1].String(), url)
		subMap[url] = ips[len(ips)-1].String()
	}
	fmt.Printf("Find %d subdomain", len(subMap))

	// output to a file
	if isFlagPassed("o") {
		f, err := os.Create(*output)
		if err != nil {
			panic(err)
		}
		for url, ip := range subMap {
			f.WriteString(url + " " + ip + "\n")
		}
	}

}
func removeDuplicates(elements []string) []string { // change string to int here if required
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{} // change string to int here if required
	result := []string{}             // change string to int here if required

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
