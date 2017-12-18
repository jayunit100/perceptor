package scanner

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/bdsengineering/go-hub-client/hubclient"
	"github.com/prometheus/common/log"
)

// HitHubAPI is an example; don't use it in production
func HitHubAPI() {
	baseURL := "https://localhost"
	username := "sysadmin"
	password := "blackduck"
	pf, err := NewProjectFetcher(username, password, baseURL)
	if err != nil {
		panic("unable to instantiate ProjectFetcher: " + err.Error())
	}
	project := pf.FetchProjectOfName("openshift/origin-docker-registry")
	bytes, _ := json.Marshal(project)
	log.Infof("fetched project: %v \n\nwith json: %v", project, string(bytes[:]))
	log.Infof("bytes: %d", len(bytes))
}

func exampleHubAPI() {
	baseURL := "https://localhost"
	username := "sysadmin"
	password := "blackduck"
	client, err := hubclient.NewWithSession(baseURL, hubclient.HubClientDebugTimings)
	if err != nil {
		log.Fatalf("unable to create hub client %v", err)
		panic("oops, unable to create hub client " + err.Error())
	}
	err = client.Login(username, password)
	if err == nil {
		log.Info("success logging in!")
		projects, _ := client.ListProjects()
		log.Info("projects: %v", projects)
	} else {
		log.Errorf("unable to log in, %v", err)
	}

	projs, err := client.ListProjects()
	if err != nil {
		panic(fmt.Sprintf("error fetching project list: %v", err))
	}
	for _, p := range projs.Items {
		log.Info("proj: ", p)
		log.Info("proj href: ", p.Meta.Href)
		link, err := p.GetProjectVersionsLink()
		if err != nil {
			panic(fmt.Sprintf("error getting project versions link: %v", err))
		}
		versions, err := client.ListProjectVersions(*link)
		if err != nil {
			panic(fmt.Sprintf("error fetching project version: %v", err))
		}
		log.Info("project versions for url: ", link.Href, ": ", versions, "\n\n")

		for _, v := range versions.Items {
			log.Info("version: ", v)
			log.Info("version href: ", v.Meta.Href)
			codeLocationsLink, err := v.GetCodeLocationsLink()
			if err != nil {
				panic(fmt.Sprintf("error getting code locations link: %v", err))
			}
			//codeLocations, err := client.GetCodeLocation(*codeLocationsLink)
			codeLocations, err := client.ListCodeLocations(*codeLocationsLink)
			//			client.
			if err != nil {
				panic(fmt.Sprintf("error fetching code locations: %v", err))
			}
			log.Info("code locations: ", codeLocations)
			for _, codeLocation := range codeLocations.Items {
				scanSummariesLink, err := codeLocation.GetScanSummariesLink()
				if err != nil {
					panic(fmt.Sprintf("error getting scan summaries link: %v", err))
				}
				scanSummaries, err := client.ListScanSummaries(*scanSummariesLink)
				if err != nil {
					panic(fmt.Sprintf("error fetching scan summaries: %v", err))
				}
				for _, scanSummary := range scanSummaries.Items {
					log.Info("scan summary: ", scanSummary)
				}
			}

			riskProfileLink, err := v.GetProjectVersionRiskProfileLink()
			if err != nil {
				panic(fmt.Sprintf("error getting risk profile link: %v", err))
			}
			riskProfile, err := client.GetProjectVersionRiskProfile(*riskProfileLink)
			if err != nil {
				panic(fmt.Sprintf("error fetching project version risk profile: %v", err))
			}
			log.Info("project version risk profile: ", riskProfile)

			// TODO can't get PolicyStatus for now
			// v.GetPolicyStatusLink()

			//scanSummaryLink, err := v.
			log.Info("\n\n")
		}
		log.Info("\n\n\n")
	}
	//	log.Info("projs", projs)
}
