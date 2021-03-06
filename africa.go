package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"math"
)

const (
	empty = ""
	tab   = "\t"
)

//method to return the data as a list of map[string]interface{} of the commits per repository
func (repo Repositories) commit() []map[string]interface{}{
	var commit_list []interface{}
	var commit_data []map[string]interface{}
	for k := range repo {
		commit_numbers:=make(map[string]interface{})
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
		commit_numbers["repo_name"]=repo[k].Name
		commit_numbers["no_of_commits"]=len(commit_list)
		commit_data=append(commit_data,commit_numbers)
	}
	return commit_data
}

//function to return alist of  map[string]interface{} of the language percentages from the public repos data
func (repo Repositories) languages()[]map[string]interface{}{
	language_maping:=make(map[string]bool)
	var piechart_data []map[string]interface{}
	var languages []string
	for k := range repo {
		languages=append(languages,repo[k].Language)
	}

	//generating a list of the languages without repetition
	for i :=0;i<len(languages);i++{
		language_maping[languages[i]]=true
	}
	var non_repetitive []string
	for k :=range language_maping{
		non_repetitive=append(non_repetitive,k)
	}

	//counting the numbber of times a language occurs and append to the map
	language_frequency:=make(map[string]int)
	z:=0
	for _,v:= range non_repetitive{
		language_frequency[v]=z
			for item,_ := range languages{
				if languages[item]==v{
					z++
				}else{
					continue
				}
		}
	}
	//Calculating the percentages
	var total int = 0	
	for _,v := range language_frequency{
		total += v
	}
	for k,v := range language_frequency{
		if k==""{
			language_percentages:=make(map[string]interface{})
			language_percentages["percentage"]=math.Round((float64(v)/float64(total)*100)*100) / 100
			language_percentages["language_name"]="others"
			piechart_data=append(piechart_data,language_percentages)
		}else{
		language_percentages:=make(map[string]interface{})
		language_percentages["percentage"]=math.Round((float64(v)/float64(total)*100)*100) / 100
		language_percentages["language_name"]=k
		piechart_data=append(piechart_data,language_percentages)
		}
	}
	return piechart_data
}


//get graph data
// /bar/graph/

// /pie/chart/

//Functions to handle requests and for serving the computed data as json to the client side
func GetGraph(writer http.ResponseWriter, request *http.Request) {
	auto := Repositories{}
	response, err := http.Get("https://api.github.com/orgs/AfricasTalkingLtd/repos")
	if err != nil {
		fmt.Println("Error getting data from the API")
	}
	//using switch case statements to combine the two functionalities on one function to the server mux
	switch{
		//Bar graph data
	case request.Method=="GET" && request.URL.Path=="/bar/graph/" :
		{
			responsedata, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error getting data")
				fmt.Println(request.URL.Path,http.StatusInternalServerError)
				fmt.Println(err)
			}
			response.Body.Close()
			json.Unmarshal(responsedata, &auto)
			bar_graph:=auto.commit()
			if err != nil {
				{
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Println(request.URL.Path,http.StatusInternalServerError)
					fmt.Println(err)
				}
			} else {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				encoder := json.NewEncoder(writer)
				encoder.SetIndent(empty, tab)
				encoder.Encode(bar_graph)
				fmt.Println(request.URL.Path, http.StatusOK)
			}
		}
		//Pie chart graph data
	case request.Method=="GET" && request.URL.Path=="/pie/chart/":
		{
			responsedata, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Error getting data")
				fmt.Println(request.URL.Path,http.StatusInternalServerError)
				fmt.Println(err)
			}
			response.Body.Close()
			json.Unmarshal(responsedata, &auto)
			pie_chart:=auto.languages()
			if err != nil {
				{
					writer.Header().Set("Content-Type", "application/json")
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Println(request.URL.Path,http.StatusInternalServerError)
					fmt.Println(err)
				}
			} else {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusOK)
				encoder := json.NewEncoder(writer)
				encoder.SetIndent(empty, tab)
				encoder.Encode(pie_chart)
				fmt.Println(request.URL.Path, http.StatusOK)
			}
		}
	}
}
