build:
	protoc -I. --go_out=plugins=micro:. \
		proto/consignment/consignment.proto
	docker build -t consignments-service .

run:
	docker run -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		consignments-service
