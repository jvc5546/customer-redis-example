# Customer Redis Example

This repo is used for the customer engineer exercise.
Instructions can be found at https://gist.github.com/jmacelroy/583831a3643eaac769f0d4cda2fd837d.

The prompt for the exercise can be found [here](PROMPT.md)

You can answer the questions from the prompt and complete the demo by creating a fork where you'll make your modifications before sharing.

## Prompt

The exercise consists of the following parts.

### How would you initially respond to the customer?
Immediately after reading the email:
"Hi Alice!

Let me try to get the demo going locally and I can update you with the necessary steps as soon as possible.

Thank you!
-Jonathan Caballero"

### What would your next steps be after reading their email and responding?
    1) Fork the repo and clone it locally.
    2) Run a redis server and the go application locally, without containerization simply to see if the code to be contained is properly compiling.
    3) Resolve any code bugs to get the application up and running with the redis server locally.
    4) Create a docker-compose.yml file to be able to test the container cluster locally and make sure redis and the app run as expected in the cluster environment. Based the docker compose file on this documentation: https://www.okteto.com/docs/reference/compose/ 
    5) Set up my local Okteto environment for the project: CLI installation, Okteto context, and Okteto manifest based on the docker-compose.yml file.
        - Ran into an error specifying that a host volume was needed. I followed the reference here as well https://www.okteto.com/docs/reference/compose/#volumes-string-optional
    6) Deploy the container cluster to the Okteto Cloud in and make sure the application works from the namespace's endpoint.
        - At this point, I realized that the syncing of files, recompiling the code, and restarting the server were not all done by the Okteto CLI alone.
        - Okteto was properly syncing the files, but the Go application was not updated on the Cloud cluster. Figured out that a third party library was necessary (like Nodemon for Node) that will recompile, and restart the Go server when a change occurs. Found that CompileDaemon (https://github.com/githubnemo/CompileDaemon) is a common tool for Go applications, but did see that there were a few different tools and would be happy to learn about those as well.
    7) Update the build commands to make sure that the application restarts on file changes in the Okteto Cloud.

### Summarize what work was needed for completing the demo for the customer.
"Hi Alice,

I was able to get the demo running in Okteto with the following steps:
1) The Go application code needs a few changes to properly work with Redis. Specifically, the Redis API requires the application's context as one of the arguments and it was missing on the following Redis client API: Ping, Incr, Decr, and Get.
2) 

### What improvments can Okteto make to better support this, or similar, customers in the future.
