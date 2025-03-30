# Dev

Use `make live` to run a server during development with `--watch` capability.

# Testing Dockerfile

Use `make dockerbuild` and `make dockerrun` to test the Docker image locally.

# Deploy

First make sure all the assets have been generated in your local machine (Docker will copy them):

`make live` in the root of the project

Deploy with Pulumi:

`cd pulumi && pulumi up`

This will automatically build the Docker image and push it to ECR for the new deployment.

You need to configure your aws cli credentials so Pulumi may pull them.

# Making the project prod-ready

These are some aspects to take into account if this project were to be deployed in a prod environment.

- Create separate `dev` and `prod` environments.
Right now the same Turso database is accessed in any env, 
a better configuration would use a local database for `dev`.
To do so it would be necessary to add a CI/CD pipeline to handle 
the migrations on the `prod` database.
- Add CSRF protection to the forms, given that the project is using cookie-based authentication.




# todo

- excalidraw of system
- excalidraw of deployment

- make subgroups middleware work
- make auth work (trigger webapi without http call?)

- add adhoc auth: use password-with-auth enpoint on form submission, add a hook this route in order to store the token in cookie (view Templ example), then add a "everything middleware" that reads this cookie and appends it to the auth header of all requests that have the cookie.
    * test that it works
        + https://pocketbase.io/docs/api-records/#auth-with-password
- add the "Bind(apis.RequireAuth())" method to the custom routes that require auth
    * test that it works
        + https://pocketbase.io/docs/go-routing/#registering-new-routes
- later, add logout endpoint that just deletes the cookie
    * test that it works

- start Postmark api integration once auth is working
    * https://postmarkapp.com/developer/integration/community-libraries#google-go

# Some learnings

Using Pocketbase with a go-based ServerSideRendering for the UI is not ideal, as the Pocketbase API is focused on a JavaScript SDK, and misses some of its user functionality when using it as a server framework. For example:
- there is no clear way to handle user auth from the Go code. 
- Pocketbase has its own idiom to write additional routes, the Pocketbase idiom is not compatible with the Templ idiom to render pages. Luckily Pocketbase offers a helper function `apis.WrapStdHandler(handler)` to break out of its idiom and use the standard `http.net` idiom, which Templ leverages.
