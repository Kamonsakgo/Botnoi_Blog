def skipRemainingStages = false
pipeline {
    agent any
    tools {
        jdk 'jdk17'
    }
    environment {
        REPO_NAME = 'bn-crud-ads'
        REPO_URL = 'registry.gitlab.com/botnoi-text2speech-v2/bn-crud-ads'
        GIT_CREDENTIALS_ID = 'prem-gitlab-login'
        BRANCH_NAME = getGitBranchName()
        TAG = getGitBranchName()
        SCANNER_HOME=tool 'sonar-scanner'
    }
    stages {
        stage('Check Branch start'){
            steps {
                script {
                    sh "echo start pipeline from : " + BRANCH_NAME
                    sh "echo ${REPO_URL}:${TAG}"
                    if (BRANCH_NAME == 'master'){
                        currentBuild.result = 'SUCCESS'
                        skipRemainingStages = true
                    } else if (BRANCH_NAME == 'develop'){
                        currentBuild.result = 'SUCCESS'
                        skipRemainingStages = true
                    } else if (BRANCH_NAME == 'production'){
                        currentBuild.result = 'SUCCESS'
                        skipRemainingStages = true
                    }
                }
            }
        }
        stage("Check text"){
           when {
                anyOf {
                    branch 'master'
                    branch 'main'
                    branch 'develop'
                    branch 'developer'
                    branch 'staging'
                    branch 'production'
                }
            }
            steps {
                script {
                    sh "echo from ${BRANCH_NAME}"
                }
            }
            
        }
        stage("Check User On ubuntu"){
           when {
                anyOf {
                    branch 'master'
                    branch 'main'
                    branch 'develop'
                    branch 'developer'
                    branch 'staging'
                    branch 'production'
                }
            }
            steps {
                script {
                    sh "whoami"
                }
            }
            
        }

        /// DEV SEC OPS///
        stage('Sonarqube Analysis') {
            when { anyOf { branch 'staging'; tag 'release-*' } }
            steps {
                withSonarQubeEnv('sonar-server') {
                    sh ''' $SCANNER_HOME/bin/sonar-scanner \
                    -Dsonar.projectName=${REPO_NAME} \
                    -Dsonar.projectKey=${REPO_NAME} \
                        -Dsonar.sources=. \
                        -Dsonar.exclusions=**/reports/**
                    '''
                }
            }
        }
        stage('Quality Check') {
            when { anyOf { branch 'staging'; tag 'release-*' } }
            steps {
                script {
                    waitForQualityGate abortPipeline: false, credentialsId: 'sonar-token'
                }
            }
        }
        stage('OWASP Dependency-Check Scan') {
            when { anyOf { branch 'staging'; tag 'release-*' } }
            steps {
                    dependencyCheck additionalArguments: '--scan ./ --disableYarnAudit --disableNodeAudit', odcInstallation: 'DP-Check'
                    dependencyCheckPublisher pattern: '**/dependency-check-report.xml'  
            }
        }
        stage('Trivy File Scan') {
            when { anyOf { branch 'staging'; tag 'release-*' } }
            steps {
                sh 'mkdir -p reports'
                sh 'trivy fs --format template --template "@/usr/local/share/trivy/templates/html.tpl" -o reports/file_backend_report.html --ignore-unfixed .'
            }
        }
        //////////////////////////////////////////
        stage('Docker Build Image') {
            when { anyOf { branch 'staging'; tag 'release-*' } }
            steps {
                script {
                    sh "docker build -t ${REPO_NAME} ."
                    sh "docker tag ${REPO_NAME} ${REPO_URL}:${TAG}" 
                }
            }
        }
        stage('TRIVY Image Scan') {
            when { anyOf { branch 'staging'; tag 'release-*' } }
            steps {
                // sh 'trivy image ${REPOSITORY_URI}${AWS_ECR_REPO_NAME}:${BUILD_NUMBER} > trivyimage.txt' 
                sh "echo ${TAG}"
                sh 'pwd'
                // sh 'mkdir -p reports'
                sh 'trivy image --format template --template "@/usr/local/share/trivy/templates/html.tpl" -o reports/img_backend_report.html --ignore-unfixed ${REPO_URL}:${TAG}'
            }
        }
        stage('Docker Push') {
            when { anyOf { branch 'staging'; tag 'release-*' } }
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: GIT_CREDENTIALS_ID, usernameVariable: 'DOCKER_REGISTRY_USER', passwordVariable: 'DOCKER_REGISTRY_PASSWORD')]) {
                        sh "docker login -u ${DOCKER_REGISTRY_USER} -p ${DOCKER_REGISTRY_PASSWORD} registry.gitlab.com"
                        sh "docker push ${REPO_URL}:${TAG}"
                    }
                }
            }
        }
        stage('Delete Docker Image') {
            when { anyOf { branch 'staging'; tag "release-*" } }
            steps {
                sh "docker rmi ${REPO_URL}:${TAG}"    
            }
        }
        stage('Deploy Staging') {
            when { anyOf { branch 'staging'; } }
            steps {
                script {
                    try {
                        // Determine the environment based on the branch
                        def env = (TAG.startsWith('release-')) ? 'production' : 'staging'

                        echo "Environment: ${env}"
                        echo "Repository Name: ${REPO_NAME}"
                        sh "rm -rf app-configs"
                        sh "git clone git@gitlab.com:botnoi-text2speech-v2/app-configs.git"
                        sh "cd app-configs && pwd"
                        sh "cd app-configs && git checkout ${env}"
                        
                        if (env == "staging") {
                            echo "TimeStamp: ${currentBuild.startTimeInMillis}"
                            sh "cd app-configs && sed -i 's|staggingTrigger: .*|staggingTrigger: ${env}-${currentBuild.startTimeInMillis}|' apps/${REPO_NAME}/values-${env}.yaml"
                        } else {
                            sh "cd app-configs && sed -i 's|tag: .*|tag: ${env}|' apps/${REPO_NAME}/values-${env}.yaml"
                        }

                        sh "cd app-configs && git add apps/${REPO_NAME}/values-${env}.yaml"
                        
                        def changes = sh(script: 'cd app-configs && git status --porcelain', returnStdout: true).trim()
                        if (changes) {
                            sh "cd app-configs && git commit -m 'Update Docker image tag for ${env} environment'"
                            sh "cd app-configs && git push origin ${env}"
                        } else {
                            echo "No changes to commit ${env}"
                        }
                    } catch (Exception e) {
                        echo "An error occurred: ${e.message}"
                        // Handle the error, possibly by marking the build as failed
                        error("Failed to deploy: ${e.message}")
                    }
                }
            }
        }
        stage('Post Trivy') {
            when { anyOf { branch 'staging'; tag "release-*" } }
            steps {
                archiveArtifacts artifacts: "reports/*.html", fingerprint: true
                publishHTML target : [
                        allowMissing: true,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: 'reports',
                        reportFiles: 'file_backend_report.html',
                        reportName: 'Trivy Scan',
                        reportTitles: 'Trivy File System Scan'
                    ]
                // archiveArtifacts artifacts: "reports/img_backend_report.html", fingerprint: true
                publishHTML target : [
                        allowMissing: true,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: 'reports',
                        reportFiles: 'img_backend_report.html',
                        reportName: 'Trivy Scan',
                        reportTitles: 'Trivy Images Scan'
                    ]

                script {
                    def message = "Build ${currentBuild.fullDisplayName} (${env.BUILD_URL}) ${currentBuild.currentResult}"
                    def webhookUrl = "https://discord.com/api/webhooks/1209036306953273424/ALfkG2Dquy8DQojBy3f40svpkfabhHyPw6aGgvP1NpwH2bW9PyWcGT54pOriSlAIvE44"
                    def payload = "{\"content\": \"${message}\"}"
                    
                    sh "curl -X POST -H 'Content-type: application/json' --data '${payload}' ${webhookUrl}"
                }
            }
        }
    }
}

def getGitBranchName() {
    branch_name = GIT_BRANCH
    if (branch_name.contains("origin/")) {
        branch_name = branch_name.split("origin/")[1]
    }
    return branch_name
}