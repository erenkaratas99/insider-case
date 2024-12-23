# Insider-Case

## Index
- [How to run?](https://github.com/erenkaratas99/insider-case#how-to-run "How to run?")
- [Project Structure](https://github.com/erenkaratas99/insider-case#project-structure "Project Structure")
- [Project Architecture & Working Principle](https://github.com/erenkaratas99/insider-case#project-architecture "Project Architecture & Working Principle")


## How to run?
- clone the project 
	`git clone https://github.com/erenkaratas99/insider-case.git`

i. you might want to change mock data generator params on `insider/cmd/messengerApi.go`
		pkg.GenerateMockData(mc, cfg.MessengerApi.MongoDbName, cfg.MessengerApi.MessagesColName, 10) 
		//10 represents the number of counts that you want to generate in mongo collection
### Docker
1. Run the Docker Desktop
3. Move through project direction on terminal (cd)
4. `docker-compose up --build`
	ENV variable is set by Dockerfile to 'test' by me for to select configs

##### URIs & URLs on Docker 
- mongo (if you want to connect through compass) : 
	`mongodb://root:example@localhost:27018/`  (watch-out for 27018 not 27017)
- redis (if you want to connect through redis insight etc.) :
	 `127.0.0.1:6379`
- swaggers:
	1. job-scheduler service : `http://localhost:3001/job-scheduler/swagger/index.html#/`
	2. messenger API : `http://localhost:3000/messenger/swagger/index.html#/`
- further endpoints can be found on postman collection on the repository.
### Local1
1. set an environment variable ENV=local 
2. start redis cli (i.e `brew services start redis`)
3. you need to have two separate build configurations 
	i. one for messenger API (set a program argument as messengerApi)
	ii. one for job-scheduler service (set a program argument as jobScheduler)
4. run or run on debug mode these separate services
<img width="833" alt="image" src="https://github.com/user-attachments/assets/195510b9-534c-43cd-bb1d-843ed048b582" />

##### URIs & URLs on Local 
- mongo (if you want to connect through compass) : 
	`mongodb://localhost:27017/`  
- redis (if you want to connect through redis insight etc.) :
	 `127.0.0.1:6379`
- swaggers:
	1. job-scheduler service : `http://localhost:3001/job-scheduler/swagger/index.html#/`
	2. messenger API : `http://localhost:3000/messenger/swagger/index.html#/`
- further endpoints can be found on postman collection on the repository.

## Project Structure
```yaml
build/
|-- jobScheduler/
|   |-- Dockerfile
|-- messengerApi/
|   |-- Dockerfile
cmd/
|-- jobScheduler.go
|-- messengerApi.go
|-- root.go
configs/
|-- appConfigs/
|-- errorConfigs/
docs/
|-- jobScheduler/
|-- messengerApi/
internal/
|-- apps/
|   |-- jobScheduler/
|   |   |-- entities/
|   |   |-- handlers/
|   |   |   |-- handler.go
|   |   |   |-- messenger.go
|   |   |-- jobs/
|   |-- messengerApi/
|       |-- entities/
|       |-- handlers/
|           |-- handler.go
|           |-- service.go
|-- clients/
|-- repositories/
pkg/
```
##### build
stores the deployment related files (dockerfiles, further gitlab-ci.yaml's etc.)
##### cmd
stores the init() functions to start apps separately
##### configs
stores the configuration files for both app configs and internal 6-digit error codes just in case there will be an BFF integration to manage response messages etc.
##### docs
swagger docs
##### internal
###### - apps
- stores source codes of both services of messenger API and job-scheduler service
- these two services are separated because the messenger API migth wanted to be used further integrations before sending a message (like GPT integrations, etc.) and besides that, since the built-in job mechanism of golang allocates a thread, it's better to separate it to a different runtime on deployment (best-case was using a tailored job-scheduler like `airflow` or `agendash`)

###### - clients
- stores the baseClient and methods like GET, PUT etc. (using fasthttp besides of main web framework of the project (echo) because it's faster than echo but not efficient as echo when running an app )

###### - repositories
- stores the two separate repositories for the both services
- jobs collection on mongo is related with job records
- messages collection on mongo is related with message records to be sent

##### pkg
stores the project-independent codes, middlewares like logger and utils

## Project Architecture & Working Principle

#### Messenger API
4 endpoints;
- **/messenger-job-toggle** (GET)
	accepts query param ?command=start or stop and sends http requests to job-scheduler to toggle off or on the messenger job
- **/** (GET)
	returns the objects from the messages collection with basic pagination mechanism (?limit=x&offset=y)
- **/get-two** (GET)
	takes only offset to for pagination, limit is fixed to 2
	returns two objects from the database with projection of fields : to, content, id 
	tailored for the messenger job
###### the logic behind separating this endpoint from get-all is to overcome larger objects for not to increase payload size by projection**	**the logic behind separating this endpoint from get-all is to overcome larger objects for not to increase payload size by projection
- **/commit/:messageId** (PUT)
	takes messageId sent by the job, same logic of kafka commit & ack
	if succes, set a redis key with message Id
###### redis get-key part couldn't be integrated due to lack of time but the correct place is right before get-two	redis get-key part couldn't be integrated due to lack of time but the correct place is right before get-two

#### Job-scheduler
- when the job-scheduler service starts, it immediately sets the messenger job that iterates in an interval of 2 minutes (2 mins set via config) 

##### messenger job mechanism

- there are two command endpoints for the messenger job (../handlers/messenger.go) to start and stop the job (for further integrations)
- one query endpoint for the messenger job (health check : /is-working)
- the job can be also controlled via messenger API from its endpoints (just in case for : job-scheduler will be an internal deployment (**i.e not exposed through NodePort nor ingress + ClusterIP**) )

##### messenger job working principle

1. starts iterating with 2 minutes
2. creates a job instance on jobs collection
3. stores the last_offset
4. sends requests to messenger API /get-two
5. takes two messages
6. (runs a for loop with fetched messages from messenger API) sends requests to webhook.site to get 200 or whatever (simulation of whatsapp or slack integration)
	i. a future work might be a retry mechanism
7. if success, sends a requests to the messenger API to commit, messenger API sets the incoming message status to sucess

<img width="933" alt="image" src="https://github.com/user-attachments/assets/1d8dcd12-ae7a-498c-bf46-5936862d68fc" />

