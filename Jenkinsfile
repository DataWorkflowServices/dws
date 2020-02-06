@Library('dst-shared@feature/SSM-1208_mm_1') _

dockerBuildPipeline {
        repository = "cray"
        imagePrefix = "cray"
        app = "dws-operator"
        name = "cray-dws-operator"
        description = "Operator for the Data Workflow Services stack"
        dockerfile = "build/Dockerfile"
        useLazyDocker = true
		createSDPManifest = true
}
