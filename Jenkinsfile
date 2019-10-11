@Library('dst-shared@master') _

dockerBuildPipeline {
        repository = "cray"
        imagePrefix = "cray"
        app = "dws-operator"
        name = "cray-dws-operator"
        description = "Operator for the Data Workflow Services stack"
        dockerfile = "build/Dockerfile"
        useLazyDocker = true
}
