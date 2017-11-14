package sparql

import (
	"encoding/json"
	"errors"
	"github.com/knakk/rdf"
	"io"
)

type Results struct {
	Head    *header
	Results results
}

type Notification struct {
	Spuid    string
	Sequence int
	AddedResults  Results
	RemovedResultes Results
}

type header struct {
	Link []string
	Vars []string
}

type results struct {
	Distinct bool
	Ordered  bool
	Bindings []map[string]binding
}

type binding struct {
	Type     string // "uri", "literal", "typed-literal" or "bnode"
	Value    string
	Lang     string `json:"xml:lang"`
	DataType string
}

type notification struct {
	Spuid    string
	Sequence int
	Results  notificationResults
}


type notificationResults struct {
	Head           header
	AddedResults   results
	RemovedResults results
}

func ParseFromJson(reader io.Reader) (*Results, error) {
	var res Results
	err := json.NewDecoder(reader).Decode(&res)
	return &res, err
}

func ParseNotificationJson(reader io.Reader) (*Notification, error) {
	var not notification
	err := json.NewDecoder(reader).Decode(&not)

	addedResult := Results{
		Head: &not.Results.Head,
		Results: not.Results.AddedResults,
	}

	removedResult := Results{
		Head: &not.Results.Head,
		Results: not.Results.RemovedResults,
	}
	result := &Notification{
		Spuid:not.Spuid,
		Sequence:not.Sequence,
		AddedResults:addedResult,
		RemovedResultes:removedResult,
	}

	return result, err
}

// Bindings returns a map of the bound variables in the SPARQL response, where
// each variable points to one or more RDF terms.
func (r *Results) Bindings() map[string][]rdf.Term {
	rb := make(map[string][]rdf.Term)
	for _, v := range r.Head.Vars {
		for _, b := range r.Results.Bindings {
			t, err := termFromJSON(b[v])
			if err == nil {
				rb[v] = append(rb[v], t)
			}
		}
	}

	return rb
}

func (r *Results) Vars() []string {
	return r.Head.Vars
}

// Solutions returns a slice of the query solutions, each containing a map
// of all bindings to RDF terms.
func (r *Results) Solutions() []map[string]rdf.Term {
	var rs []map[string]rdf.Term

	for _, s := range r.Results.Bindings {
		solution := make(map[string]rdf.Term)
		for k, v := range s {
			term, err := termFromJSON(v)
			if err == nil {
				solution[k] = term
			}
		}
		rs = append(rs, solution)
	}

	return rs
}

// termFromJSON converts a SPARQL json result binding into a rdf.Term. Any
// parsing errors on typed-literal will result in a xsd:string-typed RDF term.
func termFromJSON(b binding) (rdf.Term, error) {
	switch b.Type {
	case "bnode":
		return rdf.NewBlank(b.Value)
	case "uri":
		return rdf.NewIRI(b.Value)
	case "literal":
		// Untyped literals are typed as xsd:string
		if b.Lang != "" {
			return rdf.NewLangLiteral(b.Value, b.Lang)
		}
		xsdString, _ := rdf.NewIRI("http://www.w3.org/2001/XMLSchema#string")
		return rdf.NewTypedLiteral(b.Value, xsdString), nil
	case "typed-literal":
		iri, err := rdf.NewIRI(b.DataType)
		if err != nil {
			return nil, err
		}
		return rdf.NewTypedLiteral(b.Value, iri), nil
	default:
		return nil, errors.New("unknown term type")
	}
}
