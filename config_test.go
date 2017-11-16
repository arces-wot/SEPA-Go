package sepa

import (
	"strings"
	"testing"
)

var fullExample = ` { "parameters" : {
		"host" : "localhost" ,
		"ports" : {
			"http" : 8000 ,
			"https" : 8443 ,
			"ws" : 9000 ,
			"wss" : 9443}
		 ,
		"paths" : {
			"query" : "/query" ,
			"update" : "/update" ,
			"subscribe" : "/subscribe" ,
			"register" : "/oauth/register" ,
			"tokenRequest" : "/oauth/token" ,
			"securePath" : "/secure"}
		 ,
		"methods" : {
			"query" : "POST" ,
			"update" : "URL_ENCODED_POST"}
		 ,
		"formats" : {
			"query" : "JSON" ,
			"update" : "HTML"}
		 ,
		"security" : {
			"client_id" : "..." ,
			"client_secret" : "..." ,
			"jwt" : "..." ,
			"expires" : "..." ,
			"type" : "..."},

		"extended" :{
			"test" : 1,
			"data" : "Hello"
		}
		}
	}`

var onlySupported = ` { "parameters" : {
		"host" : "localhost" ,
		"ports" : {
			"http" : 8000 ,
			"ws" : 9000 }
		 ,
		"paths" : {
			"query" : "/query" ,
			"update" : "/update" ,
			"subscribe" : "/subscribe" }
		 ,
		"extended" :{
			"test" : 1,
			"data" : "Hello"
		}
		}
	}`

var missingValues = ` { "parameters" : {
		"host" : "localhost" ,
		"ports" : {
			"http" : 8000 }
		 ,
		"paths" : {
			"query" : "/query" ,
			"subscribe" : "/subscribe" }
		 ,
		"extended" :{
			"test" : 1,
			"data" : "Hello"
		}
		}
	}`

func TestConfigFromJson(t *testing.T) {
	rFull := strings.NewReader(fullExample)
	rOnlyS := strings.NewReader(onlySupported)
	rMiss := strings.NewReader(missingValues)

	config, err := ConfigFromJson(rFull)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if config.Host != "localhost" || config.Paths.Update != "/update" || config.Ports.Http != 8000 {
		t.Error("Wrong config data")
		t.FailNow()
	}

	if config.Extended["test"].(float64) != 1 || config.Extended["data"].(string) != "Hello" {
		t.Error("Wrog extended data readed")
		t.FailNow()
	}

	config, err = ConfigFromJson(rOnlyS)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	config, err = ConfigFromJson(rMiss)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if config.Ports.Ws != 9000 || config.Paths.Update != "/update" {
		t.Error("Wrong default values")
	}

}
