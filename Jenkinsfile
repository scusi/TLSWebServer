pipeline {
    agent any
   environment {
       BUILD_ENV = "TEST"
   }
    stages {
        stage('Check GO Environment') { 
            steps {
                sh 'echo $GOROOT'
                sh 'echo $GOPATH'
                // sh 'printenv | sort'
                // script {
                //     if( GIT_BRANCH == "origin/test" )
                //     {
                //         BUILD_ENV = "TEST"
                //     }
                //     if( GIT_BRANCH == "origin/master" )
                //     {
                //         BUILD_ENV = "PROD"
                //     }
                // }
            }
        }
        stage('Build') { 
            steps {
                sh 'make build'
            }
        }
        stage('Deploy') {
            steps {
                // sh 'echo $PWD'
                // sh 'echo --------------'
                // sh 'echo $BUILD_ENV'
                // sh 'echo --------------'
                sh 'mkdir -p dist'
                sh 'tar -zcvf dist/TLSWebServer.tar.gz TLSWebServer'
                sh 'aws s3 cp dist/TLSWebServer.tar.gz s3://optimus-deploy/webserver/JenkinsBuilds/${BUILD_ENV}/'
            }
        }
    }
}
