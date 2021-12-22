@Library('dst-shared@master') _

dockerBuildPipeline {
        repository = "cray"
        imagePrefix = "cray"
        app = "dws-operator"
        name = "cray-dws-operator"
        description = "Operator for the Data Workflow Services stack"
        dockerfile = "Dockerfile"
        autoJira = false
        createSDPManifest = true
        product = "kj"
}
