pipeline {
    agent any

    stages {
        stage('Check GO Environment') { 
            steps {
                sh 'echo $GOROOT'
                sh 'echo $GOPATH'
            }
        }
        stage('Build') { 
            steps {
                sh 'make build'
            }
        }
        stage('Deploy') {
            steps {
                sh 'echo $PWD'
                sh 'mkdir -p dist'
                sh 'tar -zcvf dist/TLSWebServer.tar.gz TLSWebServer config.yml'
                sh 'aws s3 cp dist/TLSWebServer.tar.gz s3://optimus-deploy/webserver/JenkinsBuilds/'
            }
        }
    }
}
