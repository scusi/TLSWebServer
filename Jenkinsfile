pipeline {
    agent any
   environment {
       BUILD_ENV = "DEV"
   }
    stages {
        stage('Check GO Environment') { 
            steps {
                sh 'echo $GOROOT'
                sh 'echo $GOPATH'
                // sh 'printenv | sort'
                script {
                    if( env.GIT_BRANCH == "origin/test" )
                    {
                        env.BUILD_ENV = "TEST"
                        sh 'echo swithed BUILD_ENV to '
                        
                    }
                    if( env.GIT_BRANCH == "origin/master" )
                    {
                        env.BUILD_ENV = "PROD"
                    }
                }
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
                sh 'echo --------------'
                sh 'echo ${BUILD_ENV}'
                sh 'echo ${env.BUILD_ENV}'
                sh 'echo --------------'
                // sh 'mkdir -p dist'
                // sh 'tar -zcvf dist/TLSWebServer.tar.gz TLSWebServer'
                // sh 'aws s3 cp dist/TLSWebServer.tar.gz s3://optimus-deploy/webserver/JenkinsBuilds/${BUILD_ENV}/'
            }
        }
    }
}
