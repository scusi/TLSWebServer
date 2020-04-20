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
                sh 'cd dist && tar -zcvf ../TLSWebServer.tar.gz TLSWebserverDevBuild/* && cd -'
                sh 'aws s3 cp TLSWebServer.tar.gz s3://optimus-deploy/webserver/JenkinsBuilds'
            }
        }
    }
}
