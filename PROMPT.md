# Customer Engineer Exercise

## Instructions
As part of the interview process we propose this offline exercise.
It should take roughly 3-4 hours and you have one week to provide your submission.
If you have questions or issues during its development, do not hesitate to reach out.
My email is jacob@okteto.com.

You can deliver the solution as a fork of this GitHub repository, answering the questions in the README and completing the configuration needed.
Additionally, provide a link to your Okteto namespace with the application running.
If you prefer to keep the repo private, you can share it with GitHub user jmacelroy (https://github.com/jmacelroy).

If you have not already shared your GitHub id with me, please send it to me at jacob@okteto.com, and I will set you up with the ability to use the Okteto deployment at https://okteto.assessment.jdm.okteto.net/.
After I have provisioned you with the ability to use that deployment, you can login with your GitHub account.
After installing the Okteto CLI, you can run okteto context use https://okteto.assessment.jdm.okteto.net, to point your cli to that deployment.
This is how our customers use their self-hosted or dedicated installations.

## Prompt

The exercise consists of the following parts.

1 - An existing customer of Okteto has been using the platform for nine months.
You remember helping them answer some easier getting started questions and remember that they love using Okteto but aren't necessarily fluent in containerization and kubernetes.
Today they have sent the following support email.

```
Hello,

I am trying to get my team setup to work on a new service but am having some issues.

I would like to deploy two services.
One service is a http api that acts as a redis client.
The code for this service is included in the repo along with a Dockerfile.
The other service is the redis server that would run alongside the http api in Okteto.
I would like to configure the api service that consumes redis for remote development in Okteto.

I am sure that we can work through the issues if I can get this demo working, https://github.com/okteto/customer-redis-example.
The dockerfile works but ran into issues creating the Okteto manifest and the docker compose file.

We used compose for our last project and would like to continue to use it.
Can you help us get the project setup to work on Okteto with Redis?

Kind Regards,
--
Alice Barr (she/her)
Platform Engineering Lead
Fake Company, Inc.
```

Read the email carefully for understanding.

2 - The SLA for first contact is 24 hours in this situation.

How would you initially respond?

Immediately after reading the email:
"Hi Alice!

Let me try to get the demo going locally and I can update you with the necessary steps as soon as possible.

Thank you!
-Jonathan Caballero"

What would your next steps be?
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

Please provide the answer to these questions by cloning the repo linked in the email and adding them
to the README.md

3 - Complete the demo for the customer.
For the sake of our assessment we have decided as a team that you should complete the demo for them.
Please complete the demo by adding a docker compose file and an okteto manifest.
Ensure that the demo runs and you can use the application as intended on the cluster where your access has been provisioned.

4 - Summarize for our customer.
Now that you have completed the demo please do a short write up of what it took to accomplish this.
Include any links to our [documentation](https://www.okteto.com/docs/welcome/overview/).
The goal would be to succinctly relay to them what it took in hopes that they may understand and repeat it.

5 - Suggest improvements for Okteto.
Provide any suggestions that you have for how Okteto can better support this customer in the future.
