# Customer Redis Example

This repo is used for the customer engineer exercise.
Instructions can be found at https://gist.github.com/jmacelroy/583831a3643eaac769f0d4cda2fd837d.

The prompt for the exercise can be found [here](PROMPT.md)

You can answer the questions from the prompt and complete the demo by creating a fork where you'll make your modifications before sharing.

## Prompt

The exercise consists of the following parts.

### How would you initially respond to the customer?
```
As soon as I finish reading the email:
"Hi Alice!

Let me try to get the demo going locally and I can update you with the necessary steps as soon as possible.

Thank you!
-Jonathan Caballero"
```
### What would your next steps be after reading their email and responding?

1) Fork the repo and clone it locally.
2) Run a redis server and the go application locally without containerization. This is to simply see if the code to be contained and deployed is able to compile.
3) Resolve any code bugs to get the application up and running with the redis server locally.
4) Create a ```docker-compose.yml``` file to be able to test the container cluster locally and make sure redis and the app run as expected in the cluster environment. Based the docker compose file on [this documentation](https://www.okteto.com/docs/reference/compose/ ).
5) Set up my local Okteto environment for the project: CLI installation, Okteto context, and Okteto manifest based on the ```docker-compose.yml``` file.
    - Ran into an error specifying that a host volume was needed. I followed the reference [here](https://www.okteto.com/docs/reference/compose/#volumes-string-optional) as well.
6) Deploy the container cluster to the Okteto Cloud with ```okteto up``` and confirm the application works from the namespace's endpoint.
    - At this point, I realized that the syncing of files, recompiling the code, and restarting the server were not all done by the Okteto CLI alone.
    - Okteto was properly syncing the files, but the Go server was not updated on the Cloud cluster. I Figured out that a third party library was necessary (like Nodemon for Node) that will recompile, and restart the Go server when a change occurs. Found that [CompileDaemon](https://github.com/githubnemo/CompileDaemon) is a common tool for Go applications, but did see that there were a few different tools and would be happy to learn about those as well.
7) Update the commands to start up the container with CompileDaemon based on this [documentation](https://www.okteto.com/docs/reference/compose/#command-string-optional). And make sure that the application restarts on file changes in the Okteto Cloud.

### Summarize what work was needed for completing the demo for the customer.
"Hi Alice,

I was able to get the demo running in Okteto with the following steps:
1) The Go application code needs a few changes to properly work with Redis. Specifically, the Redis client package was updated to **github.com/redis/go-redis/v9**. So you need to install that package with ```go get github.com/redis/go-redis/v9```

2) The redis client installed on step 1 requires the application's 'context' as one of the arguments for the following API: **Ping, Incr, Decr, and Get**. So the endpoint handlers were slighly updated so they could receive the 'context' as an argument when invoked by the main function. For example, in the following code block, you can see the variable **ctx** as the context being passed to the **Ping** API and down to the handlers to pass on to the rest of the redis API.
    ```
    func main() {
        ctx := context.Background()
        .
        .
        .

        pong, err := client.Ping(ctx).Result()
        .
        .
        .
        http.HandleFunc("/increment", func(w http.ResponseWriter, r *http.Request) {
            handlers.Increment(ctx, w, r)
        })
        })
        .
        .
        .
    }
    ```

3) Create a ```docker-compose.yml``` file in your project root directory. We will set up the two services, the api and redis server, following the Okteto Docker Compose Reference documentation [here](https://www.okteto.com/docs/reference/compose/)
    - ```version```: The first field can be the docker-compose version you are trying to use. It is optional but it is always a good practice to lock in a version of the different libraries and software you use to prevent breaking changes with any unexpected future updates. In this case version 3.8 worked well so I locked it in the file. 
    - ```services```: Then we weill begin the services block which you can learn more about [here](https://www.okteto.com/docs/reference/compose/#services-object-optional). Here is where we will have two separate blocks, each block for one of our services. The api block, and the redis block.

    For the api service:
    - ```build```: For the api service we need to specify that it will be built using a Dockerfile, so we need to add the [build field](https://www.okteto.com/docs/reference/compose/#build-stringobject-optional). Since our Dockerfile is in the same directory as the ockerfile-compose.yml file, we can simply use the relative path approach. By default, docker-compose looks for the file named Dockerfile so we do not need to specify a file name.
    - ```depend_on```: Since the api service depends on the redis service which we will define in a bit, we can specify that with the [depends_on field](https://www.okteto.com/docs/reference/compose/#depends_on-stringobject-optional). We can rely on the default condition of service_started and leave this field blank, but you can change it if you need a more specific dependency condition.
    - ```ports```: We also want access to the api service from outside of the cluster network, so we can forward the network port 8080 to the host port 8080, which will allow outside traffic to the api service. You can find details of the [port field here](https://www.okteto.com/docs/reference/compose/#ports-int-optional).
    - ```volumes```: We need to specify a volume for this service which will also serve as our host volume (the volume okteto will use to sync files between your local environment and Okteto Cloud). More information regarding the volumes can be found [here](https://www.okteto.com/docs/reference/compose/#volumes-string-optional).
    - ```command```: To get as much value as possible from Okteto, I was able to find an approach that can recompile and restart the Go server on the Okteto Cloud as you make and save changes in your local computer. This required the use of a Go package [CompileDaemon](https://github.com/githubnemo/CompileDaemon) and using it as the command to start up the api Go server. There are two ways to write out the command field that can be found [here](https://www.okteto.com/docs/reference/compose/#command-string-optional).
    ```
      api:
        build: ./
        depends_on:
          - redis
        ports:
          - 8080:8080
        volumes:
          - .:/app
        command: CompileDaemon -build='go build -o /usr/local/bin/app' -command='/usr/local/bin/app'
    ```
    
    For the redis service:
    - ```image```: Instead of the build field, we can use this field to pull in a redis image from the docker public repository. For simplicity I didn't specify a version for this image but I would recommend locking in the version to prevent future breaking changes with updates to this image in the docker public repository. More information on this field [here](https://www.okteto.com/docs/reference/compose/#image-string-optional)
    - ```ports```: Since we only want access to this service from the cluster network and not give access to outside traffic, we can simply specify the redis port without forwarding it to the host port. More information on this field [here](https://www.okteto.com/docs/reference/compose/#ports-int-optional)
    - ```volumes```: I was working in a limited storage environment, so I re-used the same volume as the api service. However, if you'd like to keep the redis data on a separate volume, feel free to specify a different volume based on the information [here](https://www.okteto.com/docs/reference/compose/#volumes-string-optional).
    ```
      redis:
        image: redis
        ports:
          - 6379
        volumes: 
          - .:/app
    ```

4) Now some small updates need to occur in the Dockerfile to match the CompileDaemon approach. 
    - For the ```RUN``` command currently on line 5, we can instead run the Go installation of the CompileDaemon package. The command for the package installation is the following: ```RUN go install -mod=mod github.com/githubnemo/CompileDaemon```
    - Then, we can update the CMD command on line 8 to use the CompileDaemon approach that will start up the Go server ```CompileDaemon -command='/usr/local/bin/app'```

5) Create the [Okteto manifest](https://www.okteto.com/docs/reference/manifest/) with the following steps. This will give you the ability to further customize your Okteto build, deployment, and development environments:
    - Create a new file with the name ```okteto.yml```. This is our Okteto manifest.
    - Name your Okteto Dev Environment to match your project name. This is optional:
        ```
            name: customer-redis-example
        ```
    - For the ```build``` section you can find more information [here](https://www.okteto.com/docs/reference/manifest/#build-object-optional). In this section we get to define how the image for our api service will be created, which will be with the use of the Dockerfile. Since the Dockerfile is on the same directory relative to the Okteto manifest, our context is of ```.``` and the section looks as follows:
        ```
            build:
              api:
                context: .
                dockerfile: Dockerfile
        ```
    NOTE: We do not need to build a custom image for our Redis service, so it's not needed in this section.

    - For the ```deploy``` section you can find more information [here](https://www.okteto.com/docs/reference/manifest/#deploy-string-optional). In this section we define how our services will be deployed. In our case, we will be using the ```docker-compose.yml``` file as our instructions for deployment. We will follow the specific documentation regarding docker-compose deployment, which can be found [here](https://www.okteto.com/docs/reference/manifest/#deploy-a-docker-compose-file). So we can complete the section as follows:
        ```
            deploy:
                compose:
                    file: docker-compose.yml
        ```
    
    - For the ```dev``` section you can find more information [here](https://www.okteto.com/docs/reference/manifest/#dev-object-optional). This is where we define how our development containers will be started when we run ```okteto up``` in the following steps. For this section we must create a block only for the api service since the redis container won't be customized and we can rely on the ```docker-compose.yml``` initialization of the basic Redis container:
        - ```command```: You can find more information about this section [here](https://www.okteto.com/docs/reference/manifest/#command-string-optional). This is the command that will be run to start up the development environment. We need CompileDaemon to run here so that it restarts the server whenever a file changes.
        - ```sync```: You can find more inoformation about this section [here](https://www.okteto.com/docs/reference/manifest/#sync-string-required). This field specified which folder will be synced in the develoment environment running in the Okteto cloud. We will specify the same directory where our files live so we can make changes that are reflected in the Okteto cloud.
        - ```forward```: You can find more inoformation about this section [here](https://www.okteto.com/docs/reference/manifest/#forward-string-optional). Here we will specify the local port forwarding to the remote port so that the api service can be accessed by traffic outside the cluster network.
        ```
            dev:
                api:
                    command:
                    - CompileDaemon
                    - -build=go build -o /usr/local/bin/app
                    - -command=/usr/local/bin/app
                    sync:
                    - .:/app
                    forward:
                    - 8080:8080
        ```
6) Finally, you can run the ```okteto up``` command to start up the development environment. To make updates in the api service, make sure to select ```api``` when the following question comes up:
    ```
    Select which development container to activate:
    Use the arrow keys to navigate: ↓ ↑ → ← 
    ▸ api
        redis
    ```
The output should look similar to the following, which means that you are ready to make changes in the api application, hit save, and check out your changes in your development environment.
```
    $ okteto up
    i  Using jvc5546 @ okteto.assessment.jdm.okteto.net as context
    i  'customer-redis-example' was already deployed. To redeploy run 'okteto deploy' or 'okteto up --deploy'
    ✓  Development container 'api' selected
    i  Images were already built. To rebuild your images run 'okteto build' or 'okteto deploy --build'
    ✓  Images successfully pulled
    ✓  Files synchronized
        Context:   okteto.assessment.jdm.okteto.net
        Namespace: jvc5546
        Name:      api
        Forward:   8080 -> 8080

    2023/02/02 20:58:45 Running build command!
    2023/02/02 20:58:45 Build ok.
    2023/02/02 20:58:45 Restarting the given command.
    2023/02/02 20:58:45 stderr: 2023/02/02 20:58:45 key set as QNMTFJmm
    2023/02/02 20:58:45 stderr: 2023/02/02 20:58:45 PONG <nil>
    2023/02/02 20:58:45 stderr: 2023/02/02 20:58:45 Starting http server on :8080
```

Please let me know if you have additional questions or if you run into any issues implementing the services.

Additionally, I would recommend going over the [docker-compose healthcheck](https://www.okteto.com/docs/reference/compose/#healthcheck-object-optional) and the [okteto probes](https://www.okteto.com/docs/reference/manifest/#probes-boolean-optional) documentation if you are interested in making sure that containers are up and running as expected. Also, fine tuning the services might provide a lot of value for your organization. These settings provided [here](https://www.okteto.com/docs/reference/manifest/#resources-object-optional) in the Okteto manifest can help with that if you do choose to performance test the application to determine the best settings. Let me know if you need any help with these as well. Thank you!"

### What improvments can Okteto make to better support this, or similar, customers in the future.
A few recommendations I would make after going through the assessment are the following:

1) Adding guidance or suggestions regarding where to specify container commands (Dockerfile vs docker-compose vs okteto manifest). Since there are a few ways to accomplish a similar environment setup this was something that I struggled with to really decide on where it was best to place a command. I am sure an approach on this sometimes comes down to preference though.

2) After getting the development environment running in the Okteto Cloud, I assumed there would be a restart after each change in my files. After some digging I realized that the Okteto file syncing is separate from the application recompiling and re-running. When I found this example using [Nodemon here](https://www.okteto.com/docs/using-dev-envs/#development-time), I realized I needed an external library to recompile and rerun. I wonder if it would be possible to document all the common commands to have a recompile and restart of common server languages.

3) I did run into some issues with volumes when trying to use a volume for the api and a separate volume for the redis service. I believe the issue was due to the storage limit of my Okteto Cloud namespace of 5GB. If I understand correctly, I read [here](https://www.okteto.com/docs/cloud/multitenancy/#resource-quotas) each volume tries by default to make a PVC (PersistenVolumeClaim) of 5GB, which gave me an error for going over the limit. I may have missed it in the documentation but I coudln't find a way to make a PVC of less than 5GB. 

Okteto is an amazing service, the documentation is so robust, and it makes development so much easier! I had an awesome time making my very first Okteto development container, thank you for the opportunity to check out Okteto up close.
