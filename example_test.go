package sepa

import (
	"fmt"
	"log"
	"sync"
	"github.com/arces-wot/SEPA-Go/sepa/sparql"
)

func SimpleUsageExample()  {
		cli, _ := NewClient(Configuration{"wot.arces.unibo.it", "wot.arces.unibo.it:8000/update",
			"wot.arces.unibo.it:8000/query", "wot.arces.unibo.it:9000/subscribe"})

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
		for _, varable := range vars {
			for _, term := range res.Bindings()[varable] {
				fmt.Println(term)
			}
		}
		var wg sync.WaitGroup
		wg.Add(1)
		sub, _ := cli.Subscribe("Select * Where { ?a ?b ?c}", func(notification *sparql.Notification) {

			for _, solution := range notification.AddedResults.Solutions() {
				for key, term := range solution {
					fmt.Print(key,": ",term," ")
				}
				fmt.Println()
			}

			wg.Done()
		})
		fmt.Println(sub)
		err = cli.Update(`INSERT DATA
	{
		<http://example/lang> <http://example/thebest> "go".
		<http://example/lang> <http://example/theworst> "visual basic"
	}`)

		log.Println(err)
		wg.Wait()
		sub.Unsubscribe()
		fmt.Println("Unsubscribed")
}
