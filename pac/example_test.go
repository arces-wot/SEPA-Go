package pac

import (
	"github.com/arces-wot/SEPA-Go/sepa"
	"fmt"
	"github.com/arces-wot/SEPA-Go/sepa/sparql"
	"sync"
)

const queryName = "FINDLANG"
const updateName = "PUTLANG"
const updateClear = "REMOVELANG"

var wg sync.WaitGroup

func Example() {

	application := initApp()

	cleaner, ec := application.newProducer(updateClear)

	if ec != nil{
		fmt.Println("Cannot create clear:",ec)
		return
	}

	//Clean the previous data
	cleaner.Produce(EMPTYDATA)

	//Create a consumer
	consumer, e := application.newConsumer(queryName)

	if e != nil{
		fmt.Println("Cannot create consumer:",e)
		return
	}



	//Used for sync
	wg.Add(1)

	if _, err := consumer.Consume(listen,EMPTYDATA); err != nil{
		fmt.Println("Cannot consume query: ",err)
		return
	}

	producer, e := application.newProducer(updateName)

	if e != nil{
		fmt.Println("Cannot create producer:",e)
		return
	}

	lang := struct {
		Lang string
	}{"Go"}
	fmt.Println("Produce: ",lang)
	producer.Produce(lang)
	//Wait for notification
	wg.Wait()

	// Output:
	// Produce:  {Go}
	// Yes!
}

func initApp()Application  {
	config := sepa.DefaultConfig()
	config.Host ="localhost"
	profile := newProfile(config)

	clearGraph := `DELETE DATA
	{
		<http://example/lang> <http://example/thebest> "GO"
	}
	`

	simpleUpdate := `INSERT DATA
	{
		<http://example/lang> <http://example/thebest> "??.Lang??"
	}
	`

	simpleQuery := ` SELECT ?c WHERE { <http://example/lang> <http://example/thebest> ?c }`



	profile.AddQuery(queryName,simpleQuery)
	profile.AddUpdate(updateName,simpleUpdate)
	profile.AddUpdate(updateClear,clearGraph)

	return newApplication(profile)
}

func listen(not *sparql.Notification)  {
	fmt.Println("Yes!")
	wg.Done()
}
