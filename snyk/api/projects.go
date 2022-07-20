package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type OwnerDetails struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Projects struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
	Branch string `json:"branch"`
	Owner  OwnerDetails
}

func GetProjectById(so SnykOptions, orgId string, intType string) (*Projects, error) {

	path := fmt.Sprintf("/org/%s/project/%s", orgId, intType)

	log.Println(path)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var proj Projects
	json.NewDecoder(res.Body).Decode(&proj)

	project := &Projects{
		Id:     proj.Id,
		Name:   proj.Name,
		Origin: proj.Origin,
		Branch: proj.Branch,
	}

	return project, nil

}

func GetProjectOwner(so SnykOptions, orgId string, intType string) (*OwnerDetails, error) {
	file, err := os.OpenFile("/Users/jay/formstack/gotest/logs_owner.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	path := fmt.Sprintf("/org/%s/project/%s", orgId, intType)

	log.Println(path)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var proj Projects
	json.NewDecoder(res.Body).Decode(&proj)

	log.Println("proj= ", proj.Owner.Name)

	project := &OwnerDetails{
		Id:   proj.Owner.Id,
		Name: proj.Owner.Name,
	}

	return project, nil

}
