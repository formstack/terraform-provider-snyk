package api

import (
	"encoding/json"
	"fmt"
	"log"
)

type Projects struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Origin string `json:"origin"`
	Branch string `json:"branch"`
	Owner  string `json:"owner"`
}

func GetProjectById(so SnykOptions, orgId string, intType string) (*Projects, error) {

	path := fmt.Sprintf("/org/%s/project/%s", orgId, intType)

	log.Println(path)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var proj map[string]string
	json.NewDecoder(res.Body).Decode(&proj)

	project := &Projects{
		Id:     proj["id"],
		Name:   proj["name"],
		Origin: proj["origin"],
		Branch: proj["branch"],
		Owner:  proj["owner"],
	}

	return project, nil

}
