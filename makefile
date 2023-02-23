build:
	cd template-validator/ && go build . && cp templateValidator ../ && rm templateValidator
	
test:
	./templateValidator --driverName=aws-ebs-csi-driver --version=0.9.14

check:
	build test

clear_results:
	cd ./RESULTS && rm *

test_ontap_san_brackets_mismatch:
	./templateValidator --driverName=netapp-ontap-san  --version=21.04

test_redhat_local_volume_file:
	./templateValidator --driverName=local-volume-file  --version=4.10

test_all_ok_redhat_odf_local:
	./templateValidator --driverName=odf-local  --version=4.10

test_all_ok_netapp_ontap_nas:
	./templateValidator --driverName=netapp-ontap-nas  --version=22.04


.PHONY: 
	build test clear_results test_ontap_san_brackets_mismatch test_redhat_local_volume_file test_all_ok_redhat_odf_local	
	test_all_ok_netapp_ontap_nas