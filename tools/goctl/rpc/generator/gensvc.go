package generator

import (
	_ "embed"
	"fmt"
	"path/filepath"

	conf "github.com/jialequ/linux-sdk/tools/goctl/config"
	"github.com/jialequ/linux-sdk/tools/goctl/rpc/parser"
	"github.com/jialequ/linux-sdk/tools/goctl/util"
	"github.com/jialequ/linux-sdk/tools/goctl/util/format"
	"github.com/jialequ/linux-sdk/tools/goctl/util/pathx"
)

//go:embed svc.tpl
var svcTemplate string

// GenSvc generates the servicecontext.go file, which is the resource dependency of a service,
// such as rpc dependency, model dependency, etc.
func (g *Generator) GenSvc(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetSvc()
	svcFilename, err := format.FileNamingFormat(cfg.NamingFormat, "service_context")
	if err != nil {
		return err
	}

	fileName := filepath.Join(dir.Filename, svcFilename+".go")
	text, err := pathx.LoadTemplate(category, svcTemplateFile, svcTemplate)
	if err != nil {
		return err
	}

	return util.With("svc").GoFmt(true).Parse(text).SaveTo(map[string]any{
		"imports": fmt.Sprintf(`"%v"`, ctx.GetConfig().Package),
	}, fileName, false)
}
