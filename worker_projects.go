package main

func workerProjects(
	managementProjectsEndpoint string,
	projectsManager *projects,
) error {
	projects, projectsErr := projectsLoad(managementProjectsEndpoint)
	if projectsErr != nil {
		return projectsErr
	}

	projectsManagerErr := projectsManager.load(projects)
	if projectsManagerErr != nil {
		return projectsManagerErr
	}

	return nil
}
