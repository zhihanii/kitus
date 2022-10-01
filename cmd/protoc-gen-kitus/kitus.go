package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/emicklei/proto"
	"log"
	"os"
	"os/exec"
	"strings"
)

var opts = new(Options)

func run() {
	var err error
	if opts.GoPb {
		err = generateGoPb()
		if err != nil {
			log.Println(err)
			return
		}
	}
	if opts.KitusPb {
		err = generateKitusPb()
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func generateGoPb() error {
	cmd := exec.Command("protoc", "--go_out=.", "--go_opt=paths=source_relative", opts.generatorOpts.FilePath)
	return cmd.Run()
}

func generateKitusPb() error {
	g, err := NewGenerator(&opts.generatorOpts)
	if err != nil {
		return err
	}

	err = g.GenerateFile()
	if err != nil {
		return err
	}

	return nil
}

type Options struct {
	GoPb          bool
	KitusPb       bool
	generatorOpts GeneratorOptions
}

type GeneratorOptions struct {
	PackageName string
	PackagePath string
	FilePath    string
	OutPath     string
}

type Service struct {
	SvcName      string
	CliName      string
	SrvName      string
	MethodNames  []string
	Methods      []string
	CliMethods   []string
	RequestTypes map[string]string
}

type Generator struct {
	opts     *GeneratorOptions
	services []*Service
	buf      bytes.Buffer
}

func NewGenerator(opts *GeneratorOptions) (*Generator, error) {
	strs := strings.Split(opts.FilePath, "/")
	strs = strs[1 : len(strs)-1]
	opts.PackagePath = strings.Join(strs, ".")
	i := strings.LastIndex(opts.FilePath, ".proto")
	if i == -1 {
		return nil, errors.New("file type error, needs .proto")
	}
	opts.OutPath = opts.FilePath[:i] + "_kitus.pb.go"

	g := &Generator{
		opts: opts,
	}

	return g, nil
}

func (g *Generator) P(s string) {
	s = s + "\n"
	g.buf.Write([]byte(s))
}

func (g *Generator) GenerateFile() error {
	reader, err := os.Open(g.opts.FilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	def, err := parser.Parse()
	if err != nil {
		return err
	}

	proto.Walk(def,
		proto.WithPackage(func(p *proto.Package) {
			s := strings.Split(p.Name, ".")
			g.opts.PackageName = s[len(s)-1]
		}),
		proto.WithService(g.handleService),
	)

	g.generatePackage(g.opts.PackageName)
	g.generateImport()
	g.generateServices()
	g.generateRegisters()
	g.generateHandlers()
	g.generateServiceInfos()

	file, err := os.Create(g.opts.OutPath)
	if err != nil {
		return err
	}
	defer file.Close()
	//[]byte是否会溢出
	_, err = file.Write(g.buf.Bytes())
	return err
}

func (g *Generator) handleService(s *proto.Service) {
	var (
		cliName      string
		srvName      string
		methodNames  []string
		methods      []string
		cliMethods   []string
		requestTypes = map[string]string{}
	)

	cliName = strings.ToLower(s.Name[0:1]) + s.Name[1:] + "Client"
	srvName = s.Name + "Server"

	rpcHandler := proto.WithRPC(func(rpc *proto.RPC) {
		var (
			method    string
			cliMethod string
		)

		methodNames = append(methodNames, rpc.Name)
		requestTypes[rpc.Name] = rpc.RequestType
		method = fmt.Sprintf("    %s(context.Context, *%s) (*%s, error)", rpc.Name, rpc.RequestType, rpc.ReturnsType)
		methods = append(methods, method)

		cliMethod = fmt.Sprintf("func (c *%s) %s(ctx context.Context, req *%s) (*%s, error) {\n", cliName, rpc.Name, rpc.RequestType, rpc.ReturnsType) +
			fmt.Sprintf("    resp := new(%s)\n", rpc.ReturnsType) +
			fmt.Sprintf("    err := c.kc.Call(ctx, \"%s\", req, resp)\n", "/"+g.opts.PackagePath+"."+s.Name+"/"+rpc.Name) +
			fmt.Sprintf("    if err != nil {\n") +
			fmt.Sprintf("        return nil, err\n") +
			fmt.Sprintf("    }\n") +
			fmt.Sprintf("    return resp, nil\n") +
			fmt.Sprintf("}\n")
		cliMethods = append(cliMethods, cliMethod)
	})

	for _, element := range s.Elements {
		rpcHandler(element)
	}

	g.services = append(g.services, &Service{
		SvcName:      s.Name,
		CliName:      cliName,
		SrvName:      srvName,
		MethodNames:  methodNames,
		Methods:      methods,
		CliMethods:   cliMethods,
		RequestTypes: requestTypes,
	})
}

func (g *Generator) generatePackage(name string) {
	g.P(fmt.Sprintf("package %s", name))
	g.P("")
}

func (g *Generator) generateImport() {
	g.P("import (")
	g.P("    \"context\"")
	g.P("    \"github.com/zhihanii/kitus\"")
	g.P(")")
	g.P("")
}

func (g *Generator) generateServices() {
	for _, s := range g.services {
		g.generateClient(s.CliName, s.Methods, s.CliMethods)
		g.generateServer(s.SrvName, s.Methods)
	}
}

func (g *Generator) generateClient(cliName string, methods, cliMethods []string) {
	cliNameUppder := strings.ToUpper(cliName[0:1]) + cliName[1:]
	g.P(fmt.Sprintf("type %s interface {", cliNameUppder))
	for _, method := range methods {
		g.P(method)
	}
	g.P("}")
	g.P("")

	g.P(fmt.Sprintf("type %s struct {", cliName))
	g.P("    kc kitus.Client")
	g.P("}")
	g.P("")

	for _, cliMethod := range cliMethods {
		g.P(cliMethod)
	}
}

func (g *Generator) generateServer(srvName string, methods []string) {
	g.P(fmt.Sprintf("type %s interface {", srvName))
	for _, method := range methods {
		g.P(method)
	}
	g.P("}")
	g.P("")
}

func (g *Generator) generateRegisters() {
	for _, s := range g.services {
		g.P(fmt.Sprintf("func Register%s(s kitus.Server, srv %s) {", s.SrvName, s.SrvName))
		g.P(fmt.Sprintf("    s.RegisterService(&%s_ServiceInfo, srv)", s.SvcName))
		g.P("}")
		g.P("")
	}
}

func (g *Generator) generateHandlers() {
	for _, s := range g.services {
		for _, m := range s.MethodNames {
			g.P(fmt.Sprintf("func _%s_%s_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {", s.SvcName, m))
			g.P(fmt.Sprintf("    req := new(%s)", s.RequestTypes[m]))
			g.P("    if err := dec(req); err != nil {")
			g.P("        return nil, err")
			g.P("    }")
			g.P(fmt.Sprintf("    return srv.(%s).%s(ctx, req)", s.SrvName, m))
			g.P("}")
			g.P("")
		}
	}
}

func (g *Generator) generateServiceInfos() {
	for _, s := range g.services {
		g.P(fmt.Sprintf("var %s_ServiceInfo = kitus.ServiceInfo{", s.SvcName))
		g.P(fmt.Sprintf("    ServiceName: \"%s\",", g.opts.PackagePath+"."+s.SvcName))
		g.P(fmt.Sprintf("    HandlerType: (*%s)(nil),", s.SrvName))
		g.P("    Methods: map[string]*kitus.MethodInfo{")
		for _, m := range s.MethodNames {
			g.P(fmt.Sprintf("        \"%s\": {", m))
			g.P(fmt.Sprintf("            MethodName: \"%s\",", m))
			g.P(fmt.Sprintf("            Handler: _%s_%s_Handler,", s.SvcName, m))
			g.P("        },")
		}
		g.P("    },")
		g.P("}")
		g.P("")
	}
}
