DIR := ${CURDIR}
.PHONY: dws_crds
dws_crds:	
	make -C "$(DIR)/.." manifests
	mkdir -p $(DIR)/crd/bases
	cp -r $(DIR)/../config/crd/bases $(DIR)/crd/bases

.PHONY: unit_tests
unit_tests:
	docker build --no-cache --target=testing $(DIR)

.PHONY: test_cleanup
test_cleanup:
	rm -rf $(DIR)/crd

tests: dws_crds unit_tests test_cleanup
