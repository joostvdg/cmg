// https://github.com/GoogleContainerTools/kaniko/issues/835
pipeline {
    agent none
    environment {
        REPO        = 'caladreas'
        IMAGE       = 'cmg-preview'
        TAG         = "0.0.0-${BUILD_NUMBER}"
    }
    stages {
        stage('Image Build') {
            when { changeRequest target: 'main' }
            parallel {
                stage('Kaniko') {
                    agent {
                        kubernetes {
                        //cloud 'kubernetes'
                        label 'cmg-kaniko-build'
                        yaml """
kind: Pod
metadata:
  name: kaniko
spec:
  containers:
  - name: golang
    image: golang:1.16
    command:
    - cat
    tty: true
  - name: kaniko
    image: gcr.io/kaniko-project/executor:debug
    imagePullPolicy: Always
    command:
    - /busybox/cat
    tty: true
    volumeMounts:
      - name: jenkins-docker-cfg
        mountPath: /kaniko/.docker
    env:
      - name: DOCKER_CONFIG
        value: /kaniko/.docker
  volumes:
  - name: jenkins-docker-cfg
    projected:
      sources:
      - secret:
          name: docker-credentials
          items:
            - key: .dockerconfigjson
              path: config.json
"""
                        }
                    }
                    stages {
                        stage('Build with Kaniko') {
                            steps {
                                sh 'echo image fqn=${REPO}/${IMAGE}:${TAG}'
                                container(name: 'kaniko', shell: '/busybox/sh') {
                                    withEnv(['PATH+EXTRA=/busybox']) {
                                        sh '''#!/busybox/sh
                                        /kaniko/executor -f `pwd`/Dockerfile -c `pwd` --cleanup --cache=true --destination ${REPO}/${IMAGE}:${TAG} --destination ${REPO}/${IMAGE}:latest
                                        '''
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        stage('Image Test') {
            when { changeRequest target: 'main' }
            parallel {
                stage('Agent') {
                    agent {
                        kubernetes {
                            label 'agent-test'
                            containerTemplate {
                                name 'agent'
                                image "${REPO}/${IMAGE}:${TAG}"
                                ttyEnabled true
                                command 'cat'
                            }
                        }
                    }
                    stages {
                        stage('Verify Image') {
                            steps {
                                container('agent') {
                                    sh 'cmg mapgen --gameType 1 --max 365 --min 156 --minResource 65 --max300 22'
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
