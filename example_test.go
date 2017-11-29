package sepa

import (
	"fmt"
	"github.com/arces-wot/SEPA-Go/sepa/sparql"
	"log"
	"sort"
	"sync"
)

func Example() {
	cli := NewDefaultClient("localhost")

	err := cli.Update(`INSERT DATA
	{
		<http://example/lang> <http://example/thebest> "go"
	}`)

	err = cli.Update(`DELETE DATA
	{
		<http://example/lang> <http://example/thebest> "go".
		<http://example/lang> <http://example/theworst> "visual basic"
	}`)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := cli.Query("Select * Where { ?a ?b ?c}")

	if err != nil {
		fmt.Println(err)
		return
	}
	vars := res.Vars()
	fmt.Println("Printing vars:")
	for _, varable := range vars {
		for _, term := range res.Bindings()[varable] {
			fmt.Println(term)
		}
	}
	fmt.Println("")
	var wg sync.WaitGroup
	wg.Add(1)
	sub, _ := cli.Subscribe("Select * Where { ?a ?b ?c}", func(notification *sparql.Notification) {
		fmt.Println("Printing notification:")
		for _, solution := range notification.AddedResults.Solutions() {

			//Sorting keys to get a predictable output
			// see https://blog.golang.org/go-maps-in-action#TOC_7.
			var keys []string
			for k := range solution {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, key := range keys {
				fmt.Print(key, ": ", solution[key], " ")
			}
			fmt.Println(".")
		}

		wg.Done()
	})

	err = cli.Update(`INSERT DATA
	{
		<http://example/lang> <http://example/thebest> "go".
		<http://example/lang> <http://example/theworst> "visual basic"
	}`)

	log.Println(err)
	wg.Wait()
	sub.Unsubscribe()
	fmt.Println("")
	fmt.Println("Unsubscribed")

	// Output:
	//Printing vars:
	//http://example/lang
	//http://example/thebest
	//Go
	//
	//Printing notification:
	//a: http://example/lang b: http://example/thebest c: go .
	//a: http://example/lang b: http://example/theworst c: visual basic .
	//
	//Unsubscribed
}
