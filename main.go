package main

import (
	"github.com/everxyz/protoc-gen-connection/modules"
	"github.com/lyft/protoc-gen-star"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	optional := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	pgs.
		Init(pgs.DebugEnv("DEBUG_PGV"), pgs.SupportedFeatures(&optional)).
		RegisterModule(modules.GRPCClient()).
		//RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
