package project

import (
	"actions/cleanup/deployment"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/atlas/mongodbatlas"
)

func PrepareListToDelete(projects *mongodbatlas.Projects, deleteAll bool, lifetimeInHours int) ([]*mongodbatlas.Project, error) {
	var projectList []*mongodbatlas.Project
	if deleteAll {
		projectList = projects.Results
	} else {
		for _, project := range projects.Results {
			createdTime, err := time.Parse(time.RFC3339, project.Created) // check format
			if err != nil {
				return nil, fmt.Errorf("error parsing project creation time: %w", err)
			}
			if time.Since(createdTime) > time.Duration(lifetimeInHours)*time.Hour {
				projectList = append(projectList, project)
			}
		}
	}
	return projectList, nil
}

func DeleteProjects(ctx context.Context, client mongodbatlas.Client, projectList []*mongodbatlas.Project) bool {
	ok := true
	for _, project := range projectList {
		log.Printf("Deleting project with name %s and id %s", project.Name, project.ID)

		err := deleteAllPE(ctx, client.PrivateEndpoints, project.ID)
		if err != nil {
			log.Printf("error deleting private endpoints: %s", err)
		}
		err = DeleteAllNetworkPeers(ctx, client.Peers, project.ID)
		if err != nil {
			log.Printf("error deleting network peers: %s", err)
		}
		err = deployment.DeleteAllDeployments(ctx, client.Clusters, project.ID)
		if err != nil {
			log.Printf("error deleting deployments: %s", err)
		}
		err = deployment.DeleteAllServerless(ctx, client.ServerlessInstances, project.ID)
		if err != nil {
			log.Printf("error deleting serverless: %s", err)
		}
		err = deployment.DeleteAllDataFederationInstances(ctx, client.DataFederation, project.ID)
		if err != nil {
			log.Printf("error deleting DataFederation: %s", err)
		}
		err = deployment.DeleteAllAdvancedClusters(ctx, client.AdvancedClusters, project.ID)
		if err != nil {
			log.Printf("error deleting advanced clusters: %s", err)
		}
		_, err = client.Projects.Delete(context.Background(), project.ID)
		if err != nil {
			ok = false
			log.Printf("error deleting project: %s", err)
		} else {
			log.Printf("Project successufully deleted")
		}
	}
	return ok
}
