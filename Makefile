.PHONY: proto
proto:
	protoc -I=model --go_out=model/proto model/model.proto