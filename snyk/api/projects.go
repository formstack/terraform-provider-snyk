package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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

type TargetImport struct {
	Target Target `json:"target"`
}

type Target struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
	Id     string `json:"id,omitempty"`
}

type ImportStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
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

func ImportProject(so SnykOptions, orgId string, integration string, repository_owner string, repository_name string, branch string) (*Target, error) {
	t := TargetImport{
		Target{
			Owner:  repository_owner,
			Name:   repository_name,
			Branch: branch,
		},
	}

	path := fmt.Sprintf("/org/%s/integrations/%s/import", orgId, integration)
	body, _ := json.Marshal(t)
	res, err := clientDo(so, "POST", path, body)

	if err != nil {
		return nil, err
	}

	//get the import job url so we can check the status
	jobUrl := fmt.Sprintf("%v", res.Header["Location"])
	jobUrl = strings.Trim(jobUrl, "]")
	jobId := strings.Split(jobUrl, "import/")

	defer res.Body.Close()

	//this needs to be fixed.  it only works on the repository root
	full_name := repository_owner + "/" + repository_name

	//check the API to see if the project imported
	imports, err := getImportStatus(so, orgId, integration, jobId[1])

	if err != nil {
		return nil, err
	}

	//once we get a successful import, look up the project ID so we can return it in our struct
	if imports == true {
		projects, err := GetProjectByName(so, orgId, full_name)
		if err != nil {
		}
		returnData := &Target{
			Id:     projects.Id,
			Name:   projects.Name,
			Branch: projects.Branch,
		}

		return returnData, nil

	}

	return nil, err

}

func getImportStatus(so SnykOptions, orgId string, integration string, jobId string) (bool, error) {
	path := fmt.Sprintf("/org/%s/integrations/%s/import/%s", orgId, integration, jobId)

	//check api every 10 seconds, 20 times to get success message on porject import
	for i := 0; i < 20; i++ {
		res, err := clientDo(so, "GET", path, nil)

		if err != nil {

			return false, err
		}

		var stat ImportStatus
		json.NewDecoder(res.Body).Decode(&stat)
		if stat.Status == "complete" {
			return true, nil
		}
		time.Sleep(10 * time.Second)

	}
	return false, nil
}
