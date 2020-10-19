#!/usr/bing/env groovy

pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh 'make install-go'
                sh 'make build MODE=production'
                archiveArtifacts artifacts: 'build/production/*', fingerprint: true
            }
        }
        stage('Deploy') {
            when {
                expression {
                    currentBuild.result == null || currentBuild.result == 'SUCCESS'
                }
            }
            steps {
                sh 'make config'
                sh 'make publish MODE=production'
                sh 'make restart'
            }
        }
    }
}
