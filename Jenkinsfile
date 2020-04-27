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
                sh 'mkdir -p dist'
                sh 'tar -zcvf dist/TLSWebServer.tar.gz TLSWebServer'
                script {
                    if( env.GIT_BRANCH == "origin/dev" )
                    {
                        sh 'aws s3 cp dist/TLSWebServer.tar.gz s3://optimus-deploy/webserver/JenkinsBuilds/DEV/'
                    } 
                    if( env.GIT_BRANCH == "origin/test" )
                    {
                        sh 'aws s3 cp dist/TLSWebServer.tar.gz s3://optimus-deploy/webserver/JenkinsBuilds/TEST/'
                    } 
                    if( env.GIT_BRANCH == "origin/master" )
                    {
                        sh 'aws s3 cp dist/TLSWebServer.tar.gz s3://optimus-deploy/webserver/JenkinsBuilds/PROD/'
                    } 
                }
            }
        }
    }
}
