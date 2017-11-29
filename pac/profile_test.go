package pac

import (
	"testing"
	"github.com/arces-wot/SEPA-Go/sepa"
	"strings"
)
var jsap_sample = `{
	"parameters":
	{
		"host" : "localhost",
		"ports":
		{
			"http":8000,
			"https":8443,
			"ws":9000,
			"wss":9443
		},
		"paths":
		{
			"query":"/query",
			"update":"/update",
			"subscribe":"/subscribe",
			"register":"/oauth/register",
			"tokenRequest":"/oauth/token",
			"securePath":"/secure"
		},
		"methods" : {
			"query" : "POST" ,
			"update" : "URL_ENCODED_POST"}
		 ,
		"formats" : {
			"query" : "JSON" ,
			"update" : "HTML"}
	},
	"namespaces":
	{
		"wot":"http://wot.arces.unibo.it/sepa#",
		"rdf":"http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	},
	"updates":
	{
		"TD_INIT":
		{
			"sparql":"DELETE {?thing wot:isDiscoverable ?discoverable . ?thing td:hasName ?oldName . ?thing wot:hasComponent ?component. ?component rdf:type td:Thing . ?thing td:hasProperty ?property. ?property td:hasName ?pName. ?property td:hasStability ?pStability. ?property td:isWritable ?pWrite. ?pValueType rdf:type ?pDataType . ?pValueType dul:hasDataValue ?pDataValue . ?thing td:hasEvent ?event. ?event td:hasName ?eName. ?event td:forProperty ?eProperty . ?thing td:hasAction ?action. ?action td:hasName ?aName. ?action wot:isAccessibleBy ?aProtocol. ?action td:forProperty ?aProperty} INSERT {?thing rdf:type td:Thing . ?thing td:hasName ?name. ?thing wot:isDiscoverable 'true'} WHERE { OPTIONAL {?thing rdf:type td:Thing. ?thing wot:isDiscoverable ?discoverable . ?thing td:hasName ?oldName} . OPTIONAL {?thing wot:hasComponent ?component. ?component rdf:type td:Thing} . OPTIONAL {?thing td:hasProperty ?property. ?property td:hasName ?pName. ?property td:hasStability ?pStability. ?property td:isWritable ?pWrite. ?pValueType rdf:type ?pDataType . ?pValueType dul:hasDataValue ?pDataValue} . OPTIONAL {?thing td:hasEvent ?event. ?event td:hasName ?eName. OPTIONAL {?event td:forProperty ?eProperty}} . OPTIONAL {?thing td:hasAction ?action. ?action td:hasName ?aName. ?action wot:isAccessibleBy ?aProtocol. OPTIONAL {?action td:forProperty ?aProperty}} }",
			"forcedBindings":
			{
				"thing":
				{
					"type":"uri",
					"value":""
				},
				"name":
				{
					"type":"literal",
					"value":""
				}
			}
		},
		 "UPDATE_DISCOVER" :
    {
      "sparql" : " DELETE { ?thing wot:isDiscoverable ?oldValue } INSERT {  ?thing wot:isDiscoverable ?value } WHERE {?thing wot:isDiscoverable ?oldValue}",
      "forcedBindings":
      {
        "value":
        {
          "type":"literal",
          "value":""
        },
        "thing":
        {
          "type":"uri",
          "value":""
        }
      }
    }
	},
	"queries":
	{
		"THING_DESCRIPTION":
		{
			"sparql":"SELECT ?name WHERE{ ?thing rdf:type td:Thing. ?thing wot:isDiscoverable ?discoverable . ?thing td:hasName ?name }",
			"forcedBindings":
			{
				"thing":
				{
					"type":"uri",
					"value":"Hello"
				}
			}
		},
		"THING_EVENT":
		{
			"sparql":"SELECT ?timeStamp ?value WHERE {?event rdf:type td:Event. ?event wot:hasInstance ?instance. ?instance wot:isGeneratedBy ?thing. ?instance wot:hasTimeStamp ?timeStamp. OPTIONAL {?instance td:hasOutput ?output. ?output dul:hasDataValue ?value}}",
			"forcedBindings":
			{
				"event":
				{
					"type":"uri",
					"value":""
				},
				"thing":
				{
					"type":"uri",
					"value":""
				}
			}
		}

	}
}
`
func TestProfileFromJson(t *testing.T) {
	data := strings.NewReader(jsap_sample)

	profile, e := ProfileFromJson(data)
	if e != nil {
		t.Error(e)
		t.FailNow()
	}

	checkEmbededConfiguration(profile.Configuration,t)

	querydata := struct {
		Event string
		Thing string
	}{"Impact","Table"}

	result,err := profile.GetQuery("THING_EVENT",querydata)
	
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	
	if result != "SELECT ?timeStamp ?value WHERE {Impact rdf:type td:Event. Impact wot:hasInstance ?instance. ?instance wot:isGeneratedBy Table. ?instance wot:hasTimeStamp ?timeStamp. OPTIONAL {?instance td:hasOutput ?output. ?output dul:hasDataValue ?value}}" {
		t.Error("Wrong loaded query",result)
		t.FailNow()
	}
}

func checkEmbededConfiguration(config sepa.Configuration,t *testing.T)  {
	t.Helper()
	if config.Host != "localhost" || config.Paths.Update != "/update" || config.Ports.Http != 8000 {
		t.Error("Wrong config data")
		t.FailNow()
	}

	if len(config.Extended) != 0 {
		t.Error("Wrong extended data")
		t.FailNow()
	}


}

func TestProfile_GetAdd(t *testing.T) {
	profile := newProfile(sepa.DefaultConfig())
	profile.AddQuery("Test1","SElECT ?name WHERE { ??.Id?? <hasName> ?name}")
	profile.AddQuery("Test2","SElECT ?name WHERE { ??.Id?? <hasName> ?name}")
	data := struct {
		Id string
	}{"<arces/Francesco>"}

	query, err := profile.GetQuery("Test1",data)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, err2 := profile.GetQuery("Test2",data)

	if err2 != nil {
		t.Error(err)
		t.FailNow()
	}

	if query != "SElECT ?name WHERE { <arces/Francesco> <hasName> ?name}" {
		t.Error("Wrong parsed query")
		t.FailNow()
	}

	e := profile.AddQuery("Test1", "just the same query")

	if e == nil {
		t.Error("The same query cannot be added")
		t.FailNow()
	}

}

func TestProfile_ForgetQuery(t *testing.T) {
	profile := newProfile(sepa.DefaultConfig())
	profile.AddQuery("Test","SElECT ?name WHERE { {{.Id}} <hasName> ?name}")
	profile.ForgetQuery("Test")
	_, e := profile.GetQuery("Test", nil)

	if e == nil {
		t.FailNow()
	}
}