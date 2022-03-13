package modules

import (
	pgctpl "github.com/everxyz/protoc-gen-connection/templates"
	"github.com/everxyz/protoc-gen-connection/tools"
	"github.com/iancoleman/strcase"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"path/filepath"
	"strings"
	"text/template"
)

type ClientModule struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

type Package struct {
	Name  string
	Alias string
}

type ProtoService struct {
	Name    string
	Package string
	Service string
}

type ServiceConn struct {
	Name   string
	Target string
}

func GRPCClient() *ClientModule { return &ClientModule{ModuleBase: &pgs.ModuleBase{}} }

func (m *ClientModule) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = pgsgo.InitContext(c.Parameters())
}

// Name satisfies the generator.Plugin interface.
func (m *ClientModule) Name() string { return "client" }

func (m *ClientModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	var packages []Package
	var services []ProtoService

	for _, f := range targets {
		packages = append(packages, getGoPackage(f))
		services = append(services, getProtoServices(f)...)
	}

	m.generate(packages, services)

	return m.Artifacts()
}

func (m *ClientModule) generate(packages []Package, services []ProtoService) {
	var packageNames []string
	for _, p := range packages {
		packageNames = append(packageNames, p.Name)
	}
	packagePrefix := tools.GetPrefix(packageNames)
	packageName := filepath.Base(strings.TrimSuffix(packagePrefix, "/"))

	m.AddGeneratorTemplateFile(
		filepath.Join(m.OutputPath(), "client.pb.connection.go"),
		template.Must(template.New("client").Parse(pgctpl.ClientTemplate)),
		struct {
			PackageName string
			Packages    []Package
			Services    []ProtoService
		}{
			PackageName: packageName,
			Packages:    packages,
			Services:    services,
		},
	)

	m.AddGeneratorTemplateFile(
		filepath.Join(m.OutputPath(), "client-conn.pb.connection.go"),
		template.Must(template.New("client-conn").Funcs(template.FuncMap{
			"ToLower": strings.ToLower,
		}).Parse(pgctpl.ClientConnTemplate)),
		struct {
			PackageName string
			Services    []ServiceConn
		}{
			PackageName: packageName,
			Services:    getServiceConn(m),
		},
	)
}

func getGoPackage(f pgs.File) Package {
	name := f.Descriptor().GetOptions().GetGoPackage()
	alias := strcase.ToLowerCamel(f.Package().ProtoName().String())

	return struct {
		Name  string
		Alias string
	}{
		Name:  name,
		Alias: alias,
	}
}

func getProtoServices(f pgs.File) (services []ProtoService) {
	path := filepath.Dir(f.InputPath().String())

	var meta struct {
		Service string `yaml:"service"`
	}

	tools.ParseYaml(&meta)(filepath.Join(path, "meta.yaml"))

	for _, svc := range f.Services() {
		services = append(services, struct {
			Name    string
			Package string
			Service string
		}{
			Name:    svc.Name().String(),
			Package: strcase.ToLowerCamel(svc.Package().ProtoName().String()),
			Service: meta.Service,
		})
	}

	return services
}

func getServiceConn(m *ClientModule) (conn []ServiceConn) {
	var connYaml map[string]string
	tools.ParseYaml(&connYaml)(filepath.Join(m.OutputPath(), "service.conf.connection.yaml"))

	for key, value := range connYaml {
		conn = append(conn, ServiceConn{
			Name:   key,
			Target: value,
		})
	}

	return conn
}
