String unique_id = org.apache.commons.lang.RandomStringUtils.random(4, false, true)
pipeline {
  environment {
      versionTag= readFile('./version.conf')
      gitBranch = "${env.BRANCH_NAME}"
      imageName = "memphis-cli"
      repoUrlPrefix = "memphisos"
  }

  agent {label 'big-ec2-fleet' }

  stages {
   stage('Login to Docker Hub') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'docker-hub', usernameVariable: 'DOCKER_HUB_CREDS_USR', passwordVariable: 'DOCKER_HUB_CREDS_PSW')]) {
         sh 'docker login -u $DOCKER_HUB_CREDS_USR -p $DOCKER_HUB_CREDS_PSW'
        }
      }
    }

    stage('Install GoReleaser')
        steps {
            sh """
                echo '[goreleaser]
                name=GoReleaser
                baseurl=https://repo.goreleaser.com/yum/
                enabled=1
                gpgcheck=0' | sudo tee /etc/yum.repos.d/goreleaser.repo
                sudo yum install goreleaser -y
            """
        }

    stage('Run GoReleaser - MASTER') {
        when {branch 'staging'}
        steps {
        }
   }
   stage('Run GoReleaser - LATEST') {
        when {branch 'latest'}
        steps {
            sh """
                git tag -a v${versionTag} -m "Release v${versionTag}"
                git push origin v${versionTag}
            """
            withCredentials([string(credentialsId: 'gh_token', variable: 'GITHUB_TOKEN')]) {
	        sh """
                goreleaser release --clean
            """
          }
        }
   }
   stage('Checkout to version branch') {
        when {branch 'latest'}
        steps {
            withCredentials([sshUserPrivateKey(keyFileVariable:'check',credentialsId: 'main-github')]) {
            sh """
	        git reset --hard origin/latest
	        GIT_SSH_COMMAND='ssh -i $check'  git checkout -b ${versionTag}
       	        GIT_SSH_COMMAND='ssh -i $check' git push --set-upstream origin ${versionTag}
            """
            }
        }
   }
  }
    

  post {
    unsuccessful {
        notifyFailed()
    }
    success {
        notifySuccessful() 
    }
  }
}
    

def notifySuccessful() {
  emailext (
      subject: "SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'",
      body: """<p>SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]':</p>
        <p>Check console output at &QUOT;<a href='${env.BUILD_URL}'>${env.JOB_NAME} [${env.BUILD_NUMBER}]</a>&QUOT;</p>""",
      recipientProviders: [requestor()]
    )
}
def notifyFailed() {
  emailext (
      subject: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'",
      body: """<p>FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]':</p>
        <p>Check console output at &QUOT;<a href='${env.BUILD_URL}'>${env.JOB_NAME} [${env.BUILD_NUMBER}]</a>&QUOT;</p>""",
      recipientProviders: [requestor()]
    )
}
