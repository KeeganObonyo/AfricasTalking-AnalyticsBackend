package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	empty = ""
	tab   = "\t"
)

//method to return the data as a map[string]int of the commits per repository
func (repo Repositories) commit() map[string]int{
	commit_numbers:=make(map[string]int)
	var commit_list []interface{}
	for k := range repo {
		commitURL := repo[k].CommitsURL
		commitURL = strings.Replace(commitURL, "{/sha}", "",-1)
		commits, err := http.Get(commitURL)
		if err != nil {
			fmt.Println("Error getting data")
		}
		commitsdata, err := ioutil.ReadAll(commits.Body)
		if err != nil {
			fmt.Println("Error getting data")
		}
		commits.Body.Close()
		json.Unmarshal(commitsdata, &commit_list)
		commit_numbers[repo[k].Name]=len(commit_list)
	}
	return commit_numbers
}
//function to return map[string]string of the languages from the public repos
func (repo Repositories) languages() map[string]string{
	repo_languages:=make(map[string]string)
	for k := range repo {
		repo_languages[repo[k].Name]=repo[k].Language
	}
	return repo_languages
}
//get graph data
// /bar/graph/
func GetBarGraph(writer http.ResponseWriter, request *http.Request) {
	auto := Repositories{}
	response, err := http.Get("https://api.github.com/orgs/AfricasTalkingLtd/repos")
	if err != nil {
		fmt.Println("Error getting data from alphavantage")
	}
	responsedata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error getting data from alphavantage")
	}
	response.Body.Close()
	json.Unmarshal(responsedata, &auto)
	news:=auto.commit()
	if err != nil {
		{
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Println(http.StatusInternalServerError)
			fmt.Println(err)
		}
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(news)
		fmt.Println(request.URL.Path, http.StatusOK)
	}
}


// get pie chart data
// /pie/chart/
func GetPieChart(writer http.ResponseWriter, request *http.Request) {
// 	auto := Repositories{}
// 	response, err := http.Get("https://api.github.com/orgs/AfricasTalkingLtd/repos")
// 	if err != nil {
// 		fmt.Println("Error getting data from alphavantage")
// 	}
// 	responsedata, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		fmt.Println("Error getting data from alphavantage")
// 	}
// 	response.Body.Close()
// 	json.Unmarshal(responsedata, &auto)
// 	languages:=auto.languages()
// 	if err != nil {
// 		{
// 			writer.Header().Set("Content-Type", "application/json")
// 			writer.WriteHeader(http.StatusInternalServerError)
// 			fmt.Println(http.StatusInternalServerError)
// 			fmt.Println(err)
// 		}
// 	} else {
// 		writer.Header().Set("Content-Type", "application/json")
// 		writer.WriteHeader(http.StatusOK)
// 		encoder := json.NewEncoder(writer)
// 		encoder.SetIndent(empty, tab)
// 		encoder.Encode(languages)
// 		fmt.Println(request.URL.Path, http.StatusOK)
// 	}
}