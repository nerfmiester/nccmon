package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"

	as "github.com/aerospike/aerospike-client-go"
)

type nccTestRes struct {
	Version string `json:"Version"`
	Request struct {
		Return           string `json:"Return"`
		AccountID        string `json:"AccountId"`
		ID               string `json:"Id"`
		StartDate        string `json:"StartDate"`
		EndDate          string `json:"EndDate"`
		LimitTestResults string `json:"LimitTestResults"`
		Format           string `json:"Format"`
	} `json:"Request"`
	Response struct {
		Status  string `json:"Status"`
		Code    int    `json:"Code"`
		Message string `json:"Message"`
		Account struct {
			Pages struct {
				Page struct {
					TestResults struct {
						TestResult []struct {
							LocalDateTime     string `json:"LocalDateTime"`
							TestResultDetails struct {
								ResultDetail []struct {
									ObjectURL            string  `json:"ObjectUrl"`
									TransferredBytes     string  `json:"TransferredBytes"`
									ContentSeconds       float64 `json:"ContentSeconds"`
									TotalSeconds         float64 `json:"TotalSeconds"`
									GzipSavingPercentage float64 `json:"GzipSavingPercentage"`
									StatusCode           string  `json:"StatusCode"`
								} `json:"ResultDetail"`
							} `json:"TestResultDetails"`
						} `json:"TestResult"`
					} `json:"TestResults"`
				} `json:"Page"`
			} `json:"Pages"`
		} `json:"Account"`
	} `json:"Response"`
}

type store struct {
	RunDate         string
	RunBytes        string
	ContentSecs     float64
	RunResultSecs   float64
	GzipSavePercent float64
	RunStatusCode   string
}

type siteConfidenceAPI struct {
	Version  string `json:"Version"`
	Request  string `json:"Request"`
	Response struct {
		Status  string `json:"Status"`
		Code    int    `json:"Code"`
		Message string `json:"Message"`
		APIKey  struct {
			Lifetime int    `json:"Lifetime"`
			Value    string `json:"Value"`
		} `json:"ApiKey"`
	} `json:"Response"`
}

type stores []store

// Host is .
var Host = flag.String("h", "127.0.0.1", "Aerospike server seed hostnames or IP addresses")

// Port is .
var Port = flag.Int("p", 3000, "Aerospike server seed hostname or IP address port number.")

// Namespace is .
//var Namespace = flag.String("n", "test", "Aerospike namespace.")
var storeB bytes.Buffer

// Set is .
//var Set = flag.String("s", "testset", "Aerospike set name.")

func init() {

}

// one to one
func main() {
	flag.Parse()
	_, err := initDb("127.0.0.1", 3000)
	panicOnError(err)
	fmt.Println("About to encode")
	enc := gob.NewEncoder(&storeB) // Will write to network.
	//	dec := gob.NewDecoder(&storeB) // Will read from network.
	fmt.Println("Woo hoo")
	var url string

	strre := store{}
	//	var keys []*as.Key
	runOutput := make(map[string][]store)

	apiKey := &siteConfidenceAPI{}
	apiKeyURL := "https://api.siteconfidence.co.uk/current/username/adrian.jackson@specialistholidays.com/password/JDw!QcC6&lY45Zh/Format/json"

	getJSON(apiKeyURL, apiKey)
	fmt.Printf("the api key is %v \n", apiKey.Response.APIKey.Value)
	qry := fmt.Sprintf("https://api.siteconfidence.co.uk/current/%v/Return/[Account[Pages[Page[TestResults[TestResult[LocalDateTime,TestResultDetails[ResultDetail[ObjectUrl,TransferredBytes,ContentSeconds,TotalSeconds,GzipSavingPercentage,StatusCode]]]]]]]]/AccountId/MN1A6642/Id/MN1PG24667/StartDate/2016-09-01/EndDate/2016-09-02/ShowSteps/1/LimitTestResults/9999/Format/json/", apiKey.Response.APIKey.Value)

	res := &nccTestRes{}

	getJSON(qry, res)

	for _, x := range res.Response.Account.Pages.Page.TestResults.TestResult {
		strre.RunDate = x.LocalDateTime

		for _, y := range x.TestResultDetails.ResultDetail {
			url = y.ObjectURL
			strre.RunResultSecs = y.TotalSeconds
			strre.RunStatusCode = y.StatusCode
			strre.ContentSecs = y.ContentSeconds
			strre.GzipSavePercent = y.GzipSavingPercentage
			strre.RunBytes = y.TransferredBytes
			//err = enc.Encode(strre)
			//if err != nil {
			//log.Fatal("encode error:", err)
			//}
			//fmt.Println("strre--->", strre)
			runOutput[url] = append(runOutput[url], strre)

		}

	}
	// Closures that order the Planet structure.
	totSecs := func(s1, s2 *store) bool {
		return s1.RunResultSecs < s2.RunResultSecs
	}
	for _, v := range runOutput {

		//	fmt.Println("K--->", k)
		fmt.Println("V--->", v)

		//reader := bytes.NewReader(v)
		//	dec := gob.NewDecoder(&storeB)
		//		t := make([]store, 10)

		//	err = dec.Decode(&t)
		//		if err != nil {
		//			log.Fatal("decode error 1x:", err)
		//		}
		fmt.Println("t before sort --->", v)
		By(totSecs).Sort(v)
		fmt.Println("t after sort --->", v)
		err = enc.Encode(v)
		panicOnError(err)
		fmt.Println("after encoding ---->", string(storeB.Bytes()))

		/* if err != nil {
			log.Fatal("encode error:", err)
		}
		fmt.Println("before put rec--->")
		//k1, err := putRec(k, storeB.Bytes(), client)
		fmt.Println("Past put rec--->")
		panicOnError(err)
		keys = append(keys, k1) */

		/*
			fmt.Println("K--->", k)
			fmt.Println("V--->", v)
		*/
		//	}
		//for _, v := range keys {
		/*	val, err := getRec(v, client)
			panicOnError(err)

			fmt.Println("V--->", val) */
	}
	//fmt.Printf("%v\n", runOutput)
}

func getJSON(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error reading toml file")
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)

}

func (slice *storeSorter) Len() int {
	return len(slice.stores)
}

func (slice *storeSorter) Less(i, j int) bool {
	return slice.stores[i].RunResultSecs < slice.stores[j].RunResultSecs
}

func (slice *storeSorter) Swap(i, j int) {
	slice.stores[i], slice.stores[j] = slice.stores[j], slice.stores[i]
}

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(p1, p2 *store) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(stores []store) {
	ps := &storeSorter{
		stores: stores,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// planetSorter joins a By function and a slice of Planets to be sorted.
type storeSorter struct {
	stores []store
	by     func(p1, p2 *store) bool // Closure used in the Less method.
}

func initDb(host string, port int) (*as.Client, error) {
	if Client, err := as.NewClient(host, port); err != nil {
		fmt.Println("err -> ", err)
		return nil, err
	} else {
		return Client, nil
	}
}

func putRec(k string, v []byte, client *as.Client) (*as.Key, error) {

	key, _ := as.NewKey("tester", "testSet", k)
	fmt.Println("V fputrec --->", v)
	bin := as.NewBin("bin1", v)

	log.Printf("Single Put: namespace=%s set=%s key=%s value=%s",
		key.Namespace(), key.SetName(), key.Value(), bin.Value)

	client.PutBins(as.NewWritePolicy(0, 0), key, bin)

	log.Printf("Single Get: namespace=%s set=%s key=%s", key.Namespace(), key.SetName(), key.Value())

	return key, nil
}

func getRec(k *as.Key, client *as.Client) (string, error) {

	if record, err := client.Get(as.NewPolicy(), k); err != nil {
		fmt.Println("err -> ", err)
		return "", err
	} else {

		t := []store{}
		dec := gob.NewDecoder(&storeB)
		err = dec.Decode(&t)
		if err != nil {
			log.Fatal("decode error 1:", err)
		}
		return record.String(), nil
	}

}

func panicOnError(err error) {
	if err != nil {
		fmt.Println("we know")
		panic(err)
	}
}
