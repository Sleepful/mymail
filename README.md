# Dev

First, create and env file: `cp env.example .env`, open the file and add the relevant secrets to your env file.
Use `make live` to run a server during development with `--watch` capability.
You may need to run `npm i` and `go mod tidy` to download dependencies.

Notes: 

- For Pulumi the only static ID is the SSL certificate from AWS, adjust this to your needs:
    * in file `pulumi/src/ecs.ts` adjust constant `const CERT_ARN`

- Besides `AWS`, a `postmarkapp.com` and a `turso.tech` account is necessary to fill the configuration.
    * For `postmark`, you need to perform proper configuration of your domain on their dashboard, in order to use the inbound webhook.

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
- Add persistent logs:
    * Currently the info and error logs are persisted to the filesystem, separate from the database provider to offset the load.
    * On a `prod` environment, it would be important to persist the logs in a service like AWS CloudWatch.

# Additional improvements for a long-term project

- Moving away from Fargate and implementing a CI/CD pipeline.
    * Fargate is useful for its simplicity as there is no need to manage Linux resources as there would be with EC2 or Kubernetes.
    * However, Fargate has a very slow deployment process, it may take up to 20 minutes to look at a live 



# todo

- excalidraw of system
- excalidraw of deployment

Saturday
<!-- - add adhoc auth: use password-with-auth enpoint on form submission, add a hook this route in order to store the token in cookie (view Templ example), then add a "everything middleware" that reads this cookie and appends it to the auth header of all requests that have the cookie. -->
<!--     * test that it works -->
<!--         + https://pocketbase.io/docs/api-records/#auth-with-password -->
<!-- - add the "Bind(apis.RequireAuth())" method to the custom routes that require auth -->
<!--     * test that it works -->
<!--         + https://pocketbase.io/docs/go-routing/#registering-new-routes -->
<!-- - later, add logout endpoint that just deletes the cookie -->
<!--     * test that it works -->

Sunday
<!-- - test pulumi deployment of progress so far -->
<!-- - start Postmark api integration once auth is working -->
<!--     * https://postmarkapp.com/developer/integration/community-libraries#google-go -->
<!-- - Store inbound emails in Pocketbase -->
<!--     * separate by user -->
<!-- - Add UI to look at inbound emails -->
<!--     * button to delete -->
<!--     * only own user sees own emails -->
<!-- - Add UI to create new emails -->
<!--     * test UI submit emails, verify Postmark delivery succeeds -->
<!--     * only own user sees own emails -->
<!--     * Store outbound emails in Pocketbase -->
    * add UI to look at outbound emails
        + (same as before):
            + button to delete
            + only own user sees own emails
- update Postmark inbound dashboard to use webhook from the deployment's domain instead of local dev domain

# Tech Stack

TODO: Add justification for these choices

## Infrastructure

- AWS 
    * Cloud provider.
- Docker
    * Containerization.
- [ Pulumi ]()
    * This is an IaC solution that can be configured from a variety of programming languages through an SDK.
    * Provides the capacity to use a familiar programming language to configure the AWS services, this allows easy iteration and powerful contructs within the configuration.
    * Controls the provisioning of the AWS resources from scratch, including: VPC (Subnets & Security Groups), ECS, Fargate.
- [ ECS ]()
    * AWS service to orchestrate containers.
    * The parent service to deploy Fargate instances.
    * Allows tight control over the network configuration through AWS VPC, such as:
        + The Subnet assignment to the resources.
        + The Security Groups assigned to the resources.
        + The routing tables within the VPC for its managed resources.
    * Manages a variety of services beyond the Fargate instance, such as:
        + The Internet Gateway to receive internet traffic.
        + The Load Balancer, to route internet traffic to the Fargate instances.
        + In case of private subnets, it manages the NAT Gateway.
        + In case of public subnets, the public IPs (Elastic IPs) for outbound traffic.
- [ Fargate ]()
    * The compute managed by ECS, Fargate is the "managed" version. ECS can also run on self-managed EC2. Fargate abstracts away the maintanance of the EC2 instances.
    * Allow easy deployment of managed Docker containers.
    * For an app that needs to run 24/7, the price point is slightly higher than equivalent EC2 instances, but not by much.
- [ Turso ]()
    * This SaaS offers a remote managed SQLite database for a great price. Similar to using RDS, but simpler.
    * Works like a charm with the chosen backend framework chosen for the app, as both are limited to the SQLite database.
- [ Postmark ]()
    * This SaaS offers inbound and outbound email delivery through a JSON API, this significantly simplifies email integration in comparison to common email protocols.
- Go programming language
    * The server app is written with Go.

## Reliability and Durability

### Server Reliability

Due to the usage of ECS, the reliability of the app ought to be fairly high, as ECS manages the recreation of the server on a different Availability Zone if the current Availability Zone goes down. The main issue is going to be the downtime that it takes for ECS to redeploy the task. In my experience, this downtime is anywhere between 10 and 30 minutes.

Additionally, a High Availability scheme for the server is not straight-forward to implement as it is and would require additional work. Specifically because the user sessions are stored on the server's memory, so it would be necessary to store the sessions on a durable persistance layer (such as the DB or distributed KV stores) to allow the horizontal scaling of the server across Availability Zones.

### Data Durability

Due to the usage of Turso SaaS for data, the data durability is incredibly high for any single write. All writes to the DB must be persisted to the Turso leader database before they are considered committed on the go server. As a fail-safe measure, periodic database backups may be performed.

## Server app

- Go libraries that power the app:
    * [ Pocketbase ]()
        + Framework to quickly iterate over the backend, managing many topics out of the box, such as:
            + Manages the database schema.
            + Manages a variety of operations on the database to consume the data.
            + Manages authentication for users and JWT tokens.
            + Manages authorization roles for each user, allowing them to access only their own data.
    * [ Templ ]()
        + HTML templating for Go apps. Like JSX but for Go.
    * [ SCS ](): Go session manager
        + TODO: explain cookies auth here
- Frontend libraries:
    * [ HTMX ]()
        + TODO: 
    * [ DaisyUi ]()
        + TODO: 

# Feature set

## Priorities when picking the feature-set of the solution

- Full implementation across the tech stack, including proper configuration of the cloud resources.
    * Including DNS and SSL configuration, which is not reflected in the code.
- Proper user authentication.
- Ease of use.
- Meets the requirements of receiving and sending email with an intuitive user experience.
- A balance to make an efficient use of time without cutting corners that ought to be considered cornerstones of online software.
    * Careful decisions when choosing the external dependencies to meet this balance:
        + Simplicity over complexity.
        + Efficient external tools where a hand-written solution provides little benefit.
        + Flexibility and freedom over smothering software dependencies.


TODO:
The features showcased include...
I picked these because....

# Limitations

Due to simple parsing logic for inbound emails:
- Only inbound emails with a single destination are taken into account.
    * Corollary: Emails directed at multiple accounts might not show up in the inbox.
- Only shows inbound emails with the exact email address, dots `.` and plus signs `+` are not supported.
- Only the following fields are shown in the UI:
    * Subject
    * To
    * From
    * Body in plain text
    * Date (which is not timezone-localized)

With additional time, these limitations would be amended with ease.

# Some learnings

## Things that went well 

## Things that I would have done differently a second time

Using Pocketbase with a go-based ServerSideRendering for the UI is not ideal, as the Pocketbase API is focused on a JavaScript SDK, and misses some of its user functionality when using it as a server framework. For example:
- there is no clear way to handle user auth from the Go code. 
- Pocketbase has its own idiom to write additional routes, the Pocketbase idiom is not compatible with the Templ idiom to render pages. Luckily Pocketbase offers a helper function `apis.WrapStdHandler(handler)` to break out of its idiom and use the standard `http.net` idiom, which Templ leverages.
