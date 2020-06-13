# jenkins-scheduler [![Build Status](https://api.travis-ci.org/afarid/jenkins-scheduler.svg?branch=master)][travis]



[![Docker Pulls](https://img.shields.io/docker/pulls/amrfarid/jenkins-scheduler.svg?maxAge=604800)][hub]


The jenkins scheduler enables triggering jenkins jobs with predefined scheduled rules. 

## Running this software 

### Create your config file
```yaml
jenkins:
  server: "https://jenkins.example.com/" # The jenkins server on which your jobs are
  user: "example-user" # The jenkins user which has the permissions to trigger this jpb
  token: "jenkins-token" # The token for jenkins user

jobs:
  - name: "jenkins-job-name" # Jenkins job you want to trigger (you can add many jobs)
    schedule: "0 * * * * *" # The schedule you need to configure for this job 
    parameters: # the job custom parameters
      envName: "testing"
```    
### Using the docker image
```shell script
  docker run --rm -d  --name jenkins-scheduler -v `pwd`/config.yaml:/config.yaml amrfarid/jenkins-scheduler:latest
```
    
[hub]: https://hub.docker.com/r/jenkins-scheduler
[travis]: https://travis-ci.org/afarid/jenkins-scheduler
