# Dev

Use `make live` to run a server during development with `--watch` capability.

# Testing Dockerfile

Use `dockerbuild` and `dockerrun` to test the Docker image locally.

# Deploy

Deploy with Pulumi:

`cd pulumi && pulumi up`

This will automatically build the Docker image and push it to ECR for the new deployment.

# Making the project prod-ready

Create separate `dev` and `prod` environments.
Right now the same Turso database is accessed in any env, 
a better configuration would use a local database for `dev`.
To do so it would be necessary to add a CI/CD pipeline to handle 
the migrations on the `prod` database.


# todo

- excalidraw of system
- excalidraw of deployment
