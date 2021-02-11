package main

import (
	"github.com/gogo/protobuf/vanity/command"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infof("protoc-gen-terraform %s", Version)
	logrus.Infof("protoc-gen-terraform build hash: %s", Sha)

	p := plugin.NewPlugin()

	req := command.Read()
	//log.Println(req.GetProtoFile()[3].GetName())

	//log.Println(req.GetProtoFile()[4].GetDependency())
	//logrus.Println(req.GetProtoFile()[4].GetPackage())
	// x := req.GetProtoFile()[4].GetMessageType()[2].Field[5]
	// logrus.Println(req.GetProtoFile()[4].GetMessageType()[2].GetName())
	// logrus.Println(x.GetName())
	// logrus.Println(proto.GetCastType(x))

	resp := command.GeneratePlugin(req, p, "_terraform.go")

	command.Write(resp)
}
