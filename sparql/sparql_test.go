package sparql

import (
	"testing"
	"strings"
	"github.com/knakk/sparql"
)

var queryResultJson string = `{
  "head": { "vars": [ "book" , "title" ]
  } ,
  "results": {
    "bindings": [
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book6" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Half-Blood Prince" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book7" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Deathly Hallows" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book5" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Order of the Phoenix" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book4" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Goblet of Fire" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book2" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Chamber of Secrets" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book3" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Prisoner Of Azkaban" }
      } ,
      {
        "book": { "type": "uri" , "value": "http://example.org/book/book1" } ,
        "title": { "type": "literal" , "value": "Harry Potter and the Philosopher's Stone" }
      }
    ]
  }
}`

var testNotificationJson string = `{
  "spuid": "sepa://subscription/0d057ca5-cc10-4e8a-a5d9-59d7b36f71d6",
  "sequence": 0,
  "results": {
    "head": {
      "vars": [ "book" , "title" ]
    },
    "addedResults": {
      "bindings": [
        {
          "book": { "type": "uri" , "value": "http://example.org/book/book6" } ,
          "title": { "type": "literal" , "value": "Harry Potter and the Half-Blood Prince" }
        } ,
        {
          "book": { "type": "uri" , "value": "http://example.org/book/book7" } ,
          "title": { "type": "literal" , "value": "Harry Potter and the Deathly Hallows" }
        }
      ]
    },
    "removedResults": {
      "bindings": [
        {
          "book": { "type": "uri" , "value": "http://example.org/book/book5" } ,
          "title": { "type": "literal" , "value": "Harry Potter and the Order of the Phoenix" }
        }
      ]
    }
  }
}
`

func TestParseResultsFromJson(t *testing.T)  {
	r := strings.NewReader(queryResultJson)
	result, err := sparql.ParseJSON(r)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(result.Head.Vars) != 2 {
		t.Error("Vars are not readed")
		t.Fail()
	}

}

func TestBindingFromJson(t *testing.T)  {
	r := strings.NewReader(queryResultJson)
	result, err := sparql.ParseJSON(r)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(result.Head.Vars) != 2 {
		t.Error("Vars are not readed")
		t.Fail()
	}

	if len(result.Bindings()["book"]) != 7 {
		t.Error("some book unreaded")
		t.Fail()
	}

	if result.Bindings()["title"][0].String() != "Harry Potter and the Half-Blood Prince" {
		t.Error("first bining wrong", result.Bindings()["book"][0].String())
		t.Fail()
	}
}

func TestParseNotificationJson(t *testing.T) {
	r := strings.NewReader(testNotificationJson)
	result, err := ParseNotificationJson(r)

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if result.Spuid != "sepa://subscription/0d057ca5-cc10-4e8a-a5d9-59d7b36f71d6" {
		t.Error("Spuid not valid")
		t.Fail()
	}

	if len(result.AddedResults.Vars()) != len(result.RemovedResultes.Vars()) &&
		len(result.AddedResults.Vars()) == 2 {
		t.Error("Malformed vars")
		t.Fail()
	}

	if len(result.AddedResults.Bindings()["book"]) != 2 {
		t.Error("Missing some added result")
		t.Fail()
	}

	if len(result.RemovedResultes.Bindings()["book"]) != 1 {
		t.Error("Missing some removed result")
		t.Fail()
	}

	if result.AddedResults.Bindings()["title"][0].String() != "Harry Potter and the Half-Blood Prince" {
		t.Error("first added binding wrong", result.AddedResults.Bindings()["title"][0])
		t.Fail()
	}



	if result.RemovedResultes.Bindings()["title"][0].String() != "Harry Potter and the Order of the Phoenix" {
		t.Error("first added binding wrong", result.AddedResults.Bindings()["title"][0])
		t.Fail()
	}

}
