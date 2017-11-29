package pac

import (
	"github.com/arces-wot/SEPA-Go/sepa"
	"fmt"
	"github.com/arces-wot/SEPA-Go/sepa/sparql"
	"sync"
)

var wg sync.WaitGroup

func Example() {
	config := sepa.DefaultConfig()
	config.Host ="localhost"
	profile := newProfile(config)

	clearGraph := `DELETE DATA
	{
		<http://example/lang> <http://example/thebest> "??.Lang??"
	}
	`

	simpleUpdate := `INSERT DATA
	{
		<http://example/lang> <http://example/thebest> "??.Lang??"
	}
	`

	simpleQuery := ` SELECT ?c WHERE { ??.A?? ??.B?? ?c }`

	const queryName = "FINDLANG"
	const updateName = "PUTLANG"
	const updateClear = "REMOVELANG"

	profile.AddQuery(queryName,simpleQuery)
	profile.AddUpdate(updateName,simpleUpdate)
	profile.AddUpdate(updateClear,clearGraph)

	application := newApplication(profile)

	cleaner, ec := application.newProducer(updateClear)

	if ec != nil{
		fmt.Println("Cannot create clear:",ec)
		return
	}
	//Clean the previous data
	lang := struct {
		Lang string
	}{"Go"}

	cleaner.Produce(lang)

	//Create a consumer
	consumer, e := application.newConsumer(queryName)

	if e != nil{
		fmt.Println("Cannot create consumer:",e)
		return
	}
	//TODO: use no data for query
	data := struct {
		A string
		B string
	}{"<http://example/lang>","<http://example/thebest>"}

	//Used for sync
	wg.Add(1)

	if _, err := consumer.Consume(listen,data); err != nil{
		fmt.Println("Cannot consume query: ",err)
		return
	}

	producer, e := application.newProducer(updateName)

	if e != nil{
		fmt.Println("Cannot create producer:",e)
		return
	}

	fmt.Println("Produce: ",lang)
	producer.Produce(lang)
	//Wait for notification
	wg.Wait()

	// Output:
	// Produce:  {Go}
	// Yes!
}

func listen(not *sparql.Notification)  {
	fmt.Println("Yes!")
	wg.Done()
}