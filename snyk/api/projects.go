package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

type ProjectsRest struct {
	Id         string `json:"id"`
	Attributes AttributesRest
}

type AttributesRest struct {
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

func GetProjectById(so SnykOptions, orgId string, intType string) (*Projects, error) {
	path := fmt.Sprintf("/org/%s/project/%s", orgId, intType)

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
	path := fmt.Sprintf("/org/%s/project/%s", orgId, intType)

	res, err := clientDo(so, "GET", path, nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var proj Projects
	json.NewDecoder(res.Body).Decode(&proj)

	project := &OwnerDetails{
		Id:   proj.Owner.Id,
		Name: proj.Owner.Name,
	}

	return project, nil

}

func GetAllProjects(so SnykOptions, orgId string) ([]ProjectsRest, error) {
	path := fmt.Sprintf("/orgs/%s/projects", orgId)

	res, err := clientDoRest(so, "GET", path, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	projects := map[string]json.RawMessage{}
	err = json.NewDecoder(res.Body).Decode(&projects)

	if err != nil {
		return nil, err
	}

	var proj []ProjectsRest
	json.Unmarshal(projects["data"], &proj)
	return proj, nil

}

func GetProjectByName(so SnykOptions, orgId string, name string) (*ProjectsRest, error) {
	projects, err := GetAllProjects(so, orgId)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if strings.Contains(project.Attributes.Name, name) {
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

func UpdateProject(so SnykOptions, id string, orgId string, integration string, repository_owner string, repository_name string, branch string) (*Projects, error) {
	path := fmt.Sprintf("/org/%s/project/%s", orgId, id)

	t := Projects{
		Branch: branch,
	}
	body, _ := json.Marshal(t)

	res, err := clientDo(so, "PUT", path, body)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	projects, err := GetProjectById(so, orgId, id)

	if err != nil {
		return nil, err
	}
	return projects, nil
}
