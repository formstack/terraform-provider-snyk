package api

import (
	"encoding/json"
	"fmt"
	"log"
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
	log.Println("proj_name=", proj.Name)
	project := &Projects{
		Id:     proj.Id,
		Name:   proj.Name,
		Origin: proj.Origin,
		Branch: proj.Branch,
	}

	return project, nil

}

func GetProjectOwner(so SnykOptions, orgId string, intType string) (*OwnerDetails, error) {
	path := fmt.Sprintf("/org/%s/project/%s", orgId, intType)

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

func GetAllProjects(so SnykOptions, orgId string) ([]Projects, error) {
	path := fmt.Sprintf("/org/%s/projects", orgId)
	res, err := clientDo(so, "GET", path, nil)

	defer res.Body.Close()

	projects := map[string]json.RawMessage{}
	err = json.NewDecoder(res.Body).Decode(&projects)

	if err != nil {
		return nil, err
	}

	var proj []Projects
	json.Unmarshal(projects["projects"], &proj)
	return proj, nil

}

func GetProjectByName(so SnykOptions, orgId string, name string) (*Projects, error) {
	projects, err := GetAllProjects(so, orgId)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == name {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("project with name %s not found", name)

}

func DeleteProject(so SnykOptions, orgId string, intType string) error {
	path := fmt.Sprintf("/org/%s/project/%s", orgId, intType)

	_, err := clientDo(so, "DELETE", path, nil)

	return err

}
