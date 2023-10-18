package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/go-github/v41/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"

	xj "github.com/basgys/goxml2json"
)

type RequestStruct struct {
	UserName string
	Token    string
}

// type YourXMLStruct struct {
// 	YourTagName string `xml:"your_tag_name"`
// }

const prefix = "veryuniqueattrprefix-"

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/Mule/getRepo", GetRepository).Methods(http.MethodGet).Name("GetRepository")

	fmt.Println("Server is getting started...")

	fmt.Println("Listening at port 4000..")

	log.Fatal(http.ListenAndServe("localhost:8080", router))

	http.Handle("/", router)

}

func GetRepository(w http.ResponseWriter, r *http.Request) {
	// Your GitHub Personal Access Token (PAT)
	var request RequestStruct

	json.NewDecoder(r.Body).Decode(&request)

	token := request.Token

	// Create an OAuth2 token source with the PAT
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// Create an authenticated HTTP client
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)

	// Create a new GitHub client
	client := github.NewClient(oauthClient)

	// Replace "username" and "repositoryName" with the actual GitHub username and repository name
	// username := "Sneha-Jayakumar123"
	// repositoryName := "FirstRepo"

	// Get the repository
	var list *github.RepositoryListOptions
	repo, _, err := client.Repositories.List(context.Background(), request.UserName, list) // (context.Background(),&github.RepositoryListAllOptions{}) //Get(context.Background(), username, repositoryName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, repoName := range repo {
		fmt.Println(*repoName.Name)
	}

	var listContent *github.RepositoryContentGetOptions

	_, content, _, err := client.Repositories.GetContents(context.Background(), request.UserName, "FirstRepo", "", listContent) //(context.Background(),"Sneha-Jayakumar123",list)// (context.Background(),&github.RepositoryListAllOptions{}) //Get(context.Background(), username, repositoryName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// fmt.Printf("Error:%+v",content)

	for _, contentName := range content {
		fmt.Println(*contentName.Name)
	}

	contentA, _, _, err := client.Repositories.GetContents(context.Background(), request.UserName, "FirstRepo", "mule.xml", listContent) //(context.Background(),"Sneha-Jayakumar123",list)// (context.Background(),&github.RepositoryListAllOptions{}) //Get(context.Background(), username, repositoryName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("The content A", *contentA.Content)

	encodedString := *contentA.Content

	decodedData, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		fmt.Println("Error decoding the string:", err)
		return
	}

	// Convert the decoded data to a string
	decodedString := string(decodedData)

	// Print the decoded string
	fmt.Println(decodedString)

	type Root struct {
		XMLName xml.Name
	}

	var data Root
	err = xml.Unmarshal([]byte(decodedString), &data)
	if err != nil {
		fmt.Printf("Error unmarshaling XML: %v\n", err)
		return
	}

	// Print the XML tags
	response:=printXMLTags(data.XMLName, decodedString)

	var datas map[string]interface{}

	// Unmarshal the JSON string into the map
	err = json.Unmarshal([]byte(response), &datas)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("The data",datas)

	// Key you want to search for
	targetKey := "logger"

	// Call the recursive function to find the key in the JSON data
	instances := findAllInstances(datas, targetKey)

	var array []interface{}
	if len(instances) > 0 {
		// You can work with the found instances of the key here
		fmt.Printf("Found %d instances of '%s' data:\n", len(instances), targetKey)
		for i, instance := range instances {
			fmt.Printf("Instance %d: %v\n", i+1, instance)
			array=append(array,instance)
		}
	} else {
		fmt.Printf("No instances of '%s' data found\n", targetKey)
	}

	targetKey="-name"
	instances = findAllInstances(array, targetKey)
	fmt.Println("the instances",instances)

	// for index:=range array{
	// 	fmt.Println("The array",array[index])
	// 	data :=array[index]
	// 	if _, ok := data["-name"]; ok {
	// 		// The field "name" exists in the map
	// 		fmt.Println("The map has a field named 'name'.")
	// 	} else {
	// 		// The field "name" does not exist in the map
	// 		fmt.Println("The map does not have a field named 'name'.")
	// 	}
	// }

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(array)
}

func printXMLTags(name xml.Name, xmlString string) string{
	fmt.Println("Root Element:", name.Local)
	startArray := []string{}
	endArray := []string{}
	// var token xml.Token
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			{
				startArray = append(startArray, t.Name.Local)
				fmt.Println("Start Tag:", t.Name.Local)
			}
		case xml.EndElement:
			endArray = append(endArray, t.Name.Local)
			fmt.Println("End Tag:", t.Name.Local)
		}

		// if se, ok := token.(xml.StartElement); ok {
		// 	// Check if it's the element you're interested in
		// 	if se.Name.Local == "mule" {
		// 		var data YourXMLStruct
		// 		if decodeErr := decoder.DecodeElement(&data, &se); decodeErr == nil {
		// 			// You've successfully extracted the content
		// 			fmt.Println(data.YourTagName)
		// 		}
		// 	}
		// }
	}
	fmt.Printf("Start Tag:%+v \n", startArray)
	fmt.Printf("End Tag:%+v ", endArray)

	// xmlString =`xmlString`
	xml := strings.NewReader(xmlString)

	// Decode XML document
	root := &xj.Node{}
	err := xj.NewDecoder(xml).Decode(root)
	if err != nil {
		panic(err)
	}

	RemoveAttr(root)

	// Then encode it in JSON
	buf := new(bytes.Buffer)
	e := xj.NewEncoder(buf)
	err = e.Encode(root)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n",buf.String())
	return buf.String()
}

func RemoveAttr(n *xj.Node) {
	for k, v := range n.Children {
		if strings.HasPrefix(k, prefix) {
			delete(n.Children, k)
		} else {
			for _, n := range v {
				RemoveAttr(n)
			}
		}
	}
}

// func findKey(data map[string]interface{}, targetKey string) interface{} {
// 	for key, value := range data {
// 		if key == targetKey {
// 			return value
// 		}

// 		if nestedData, ok := value.(map[string]interface{}); ok {
// 			if result := findKey(nestedData, targetKey); result != nil {
// 				return result
// 			}
// 		}

// 		if nestedDataArray, ok := value.([]interface{}); ok {
// 			for _, nestedData := range nestedDataArray {
// 				if nestedDataMap, isMap := nestedData.(map[string]interface{}); isMap {
// 					if result := findKey(nestedDataMap, targetKey); result != nil {
// 						return result
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

func findAllInstances(data interface{}, targetKey string) []interface{} {
	var results []interface{}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		for key, value := range data.(map[string]interface{}) {
			if key == targetKey {
				results = append(results, value)
			}
			results = append(results, findAllInstances(value, targetKey)...)
		}
	case reflect.Slice:
		for _, value := range data.([]interface{}) {
			results = append(results, findAllInstances(value, targetKey)...)
		}
	}

	return results
}