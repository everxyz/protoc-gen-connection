package modules

import (
	pgctpl "github.com/everxyz/protoc-gen-connection/templates"
	"github.com/everxyz/protoc-gen-connection/tools"
	"github.com/iancoleman/strcase"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"path/filepath"
	"text/template"
)

type ClientModule struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
	tpl *template.Template
}

type Package struct {
	Name  string
	Alias string
}

type Service struct {
	Name    string
	Package string
	Service string
}

func GRPCClient() *ClientModule { return &ClientModule{ModuleBase: &pgs.ModuleBase{}} }

func (m *ClientModule) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.ctx = pgsgo.InitContext(c.Parameters())
}

// Name satisfies the generator.Plugin interface.
func (m *ClientModule) Name() string { return "client" }

func (m *ClientModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	tpl := template.New("client").Funcs(map[string]interface{}{
		//"package":    p.ctx.PackageName,
		//"name":       p.ctx.Name,
	})

	m.tpl = template.Must(tpl.Parse(pgctpl.ClientTemplate))

	var packages []Package
	var services []Service

	for _, f := range targets {
		packages = append(packages, getGoPackage(f))
		services = append(services, getServices(f)...)
	}

	m.generate(packages, services)

	return m.Artifacts()
}

func (m *ClientModule) generate(packages []Package, services []Service) {
	m.AddGeneratorTemplateFile("client.pb.connection.go", m.tpl, struct {
		Packages []Package
		Services []Service
	}{
		Packages: packages,
		Services: services,
	})
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

func getServices(f pgs.File) (services []Service) {
	path := filepath.Dir(f.InputPath().String())

	var meta struct {
		Service string `yaml:"service"`
	}

	tools.ParseYaml(meta)(filepath.Join(path, "meta.yaml"))

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
