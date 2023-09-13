package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

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
	full_name := repository_owner + "/" + repository_name + ":"

	//check the API to see if the project imported
	imports, err := getImportStatus(so, orgId, integration, jobId[1])

	if err != nil {
		return nil, err
	}

	//once we get a successful import, look up the project ID so we can return it in our struct
	if imports {
		projects, err := GetProjectByName(so, orgId, full_name)
		if err != nil {
			return nil, err
		}
		returnData := &Target{
			Id:     projects.Id,
			Name:   projects.Attributes.Name,
			Owner:  repository_owner,
			Branch: branch,
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
