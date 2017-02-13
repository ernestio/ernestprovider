package real

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider"
	"github.com/tidwall/gjson"

	. "github.com/gucumber/gucumber"
)

var lastSubject string
var lastID string
var lastBody []byte
var key string
var subnetID string
var niID string
var blobEndpoint string

func init() {
	key = os.Getenv("ERNEST_CRYPTO_KEY")

	And(`^testing "(.+?)" does not exist$`, func(string) {
	})

	When(`^I request "(.+?)" with "(.+?)"$`, func(subject string, fileName string) {
		pwd, _ := os.Getwd()
		filePath := path.Join(pwd, "internal", "definitions", fileName)

		dat, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println(err.Error())
		}
		dat = []byte(strings.Replace(string(dat), "$(subnetID)", subnetID, -1))
		dat = []byte(strings.Replace(string(dat), "$(networkInterfaceID)", niID, -1))
		dat = []byte(strings.Replace(string(dat), "$(blobEndpoint)", blobEndpoint, -1))

		var j map[string]interface{}
		if err := json.Unmarshal(dat, &j); err != nil {
			T.Errorf("Could not unmarshal definition " + fileName)
			panic("problem with your json on file " + fileName)
		}

		crypto := aes.New()
		j["azure_client_id"], _ = crypto.Encrypt(os.Getenv("AZURE_CLIENT_ID"), key)
		j["azure_client_secret"], _ = crypto.Encrypt(os.Getenv("AZURE_CLIENT_SECRET"), key)
		j["azure_tenant_id"], _ = crypto.Encrypt(os.Getenv("AZURE_TENANT_ID"), key)
		j["azure_subscription_id"], _ = crypto.Encrypt(os.Getenv("AZURE_SUBSCRIPTION_ID"), key)
		j["environment"], _ = crypto.Encrypt(os.Getenv("AZURE_ENVIRONMENT"), key)
		j["id"] = lastID

		dat, _ = json.Marshal(j)

		lastSubject, lastBody = ernestprovider.GetAndHandle(subject, dat, key)
		var component struct {
			ID           string `json:"id"`
			BlobEndpoint string `json:"primary_blob_endpoint"`
		}
		_ = json.Unmarshal(lastBody, &component)
		parts := strings.Split(lastSubject, ".")
		if parts[1] == "create" {
			lastID = component.ID
			switch parts[0] {
			case "azure_subnet":
				subnetID = component.ID
			case "azure_network_interface":
				niID = component.ID
			case "azure_storage_account":
				blobEndpoint = component.BlobEndpoint
			}
		}
	})

	Then(`^I should get a "(.+?)" response with "(.+?)" containing "(.+?)"$`, func(subject string, k string, v string) {
		if lastSubject != subject {
			T.Errorf("Last subject was: \n" + lastSubject)
		}
		value := gjson.Get(string(lastBody), k).String()
		if strings.Contains(v, value) {
			fmt.Println(string(lastBody))
			T.Errorf("Value " + v + " is not equal to " + value)
		}
	})

	Then(`^I should get a "(.+?)" response with "(.+?)" as "(.+?)"$`, func(subject string, k string, v string) {
		if lastSubject != subject {
			T.Errorf("Last subject was: \n" + lastSubject)
		}
		value := gjson.Get(string(lastBody), k).String()
		if v != value {
			fmt.Println(string(lastBody))
			T.Errorf("Value " + v + " is not equal to " + value)
		}
	})

	And(`^I have no messages on the buffer$`, func() {
	})
	Then(`^I should get a "(.+?)" response with "(.+?)" as "(.+?)"$`, func(subject string, k string, v string) {
		if lastSubject != subject {
			T.Errorf("Last subject was: \n" + lastSubject)
		}
		value := gjson.Get(string(lastBody), k).String()
		if v != value {
			fmt.Println(string(lastBody))
			T.Errorf("Value " + v + " is not equal to " + value)
		}
	})

	And(`^I have no messages on the buffer$`, func() {
	})

	Then(`^I should get a "(.+?)" response with body "(.+?)"$`, func(subject string, body string) {
		if lastSubject != subject {
			T.Errorf("Last subject was: \n" + lastSubject)
		}
		if string(lastBody) != body {
			T.Errorf("Last body was: \n" + string(lastBody))
		}
	})

}
