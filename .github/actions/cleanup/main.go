package main

import (
	"actions/cleanup/atlasclient"
	"actions/cleanup/project"
	"context"
	"log"
	"os"

	"go.mongodb.org/atlas/mongodbatlas"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("error getting config: %s", err)
	}

	client, err := atlasclient.SetupClient(config.PublicKey, config.PrivateKey, config.ManagerUrl)
	if err != nil {
		log.Fatalf("error creating atlas client: %s", err)
	}

	projects, _, err := client.Projects.GetAllProjects(context.Background(), &mongodbatlas.ListOptions{
		ItemsPerPage: 250,
	})
	if err != nil {
		log.Fatalf("error getting projects: %s", err)
	}
	projectList, err := project.PrepareListToDelete(projects, config.DeleteAll, config.Lifetime)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Total projects selected for deletion: ", len(projectList))
	ctx := context.Background()
	ok := project.DeleteProjects(ctx, client, projectList)
	if !ok {
		log.Printf("Not all project deleted. Please run cleaner again later")
		os.Exit(1)
	}
}
