package main

func workerProjects(
	managementProjectsEndpoint string,
) (map[string]projectData, error) {
	projects, projectsErr := projectsLoad(managementProjectsEndpoint)
	if projectsErr != nil {
		return nil, projectsErr
	}
	return projects, nil
}
