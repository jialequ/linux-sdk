package upgrade

import "github.com/jialequ/linux-sdk/tools/goctl/internal/cobrax"

// Cmd describes an upgrade command.
var Cmd = cobrax.NewCommand("upgrade", cobrax.WithRunE(upgrade))
