package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"

	// "strings"
	"github.com/google/go-github/v41/github"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/Mule/getRepo", GetRepository).Methods(http.MethodGet).Name("GetRepository")

	fmt.Println("Server is getting started...")

	fmt.Println("Listening at port 4000..")

	log.Fatal(http.ListenAndServe("localhost:8080", router))

	http.Handle("/", router)

}

type RequestStruct struct {
	UserName string
	Token    string
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
	printXMLTags(data.XMLName, decodedString)
}

func printXMLTags(name xml.Name, xmlString string) {
	fmt.Println("Root Element:", name.Local)
	startArray := []string{}
	endArray := []string{}
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
	}
	fmt.Printf("Start Tag:%+v \n", startArray)
	fmt.Printf("End Tag:%+v ", endArray)
}
