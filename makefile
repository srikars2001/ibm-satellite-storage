build:
	cd template-validator/ && go build . && cp templateValidator ../ && rm templateValidator
	
test:
	./templateValidator --driverName=aws-ebs-csi-driver --version=0.9.14

.PHONY: build test