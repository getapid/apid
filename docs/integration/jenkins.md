# Jenkins

Integrating APId is simple because of the flexibility it offers - you can either use the official docker image, or if that doesn't suit your taste, you can use the CLI straight from your shell.

## Docker

Using the official docker image is really simple, the only thing you need to do is use it as the base image for that job.   
  


```groovy
pipeline {
    agent none
    stages {
        stage('Test') {
            agent { docker 'getapid/apid:latest' }
            steps {
                "apid check -c path/to/apid.yaml"
            }
        }
    }
}
```

## CLI

Integrating the CLI is just as simple as using the docker image. Download the latest version of the CLI, make it executable and run it.   
  


```groovy
pipeline {
    agent none
    stages {
        stage('Test') {
            steps {
                "curl -o apid https://cdn.getapid.com/cli/latest/apid-latest-linux-amd64"
                "chmod u+x apid"
                "apid check -c path/to/apid.yaml"
            }
        }
    }
}
```

