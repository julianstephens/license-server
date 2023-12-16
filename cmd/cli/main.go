package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	"github.com/julianstephens/license-server/internal/config"
	"github.com/julianstephens/license-server/pkg/database"
	"github.com/julianstephens/license-server/pkg/licensemanager"
	"github.com/julianstephens/license-server/pkg/logger"
	"github.com/julianstephens/license-server/pkg/service"
)

var (
	conf = config.GetConfig()
	db   = database.GetDB()
	lm   = licensemanager.LicenseManager{Config: conf, DB: db}
)

type GenCmd struct {
	Kp  bool   `short:"k" required:"" xor:"kp" help:"Create new product key pair" `
	Lic bool   `short:"l" required:"" xor:"kp" help:"Create new empty license" `
	Pid string `arg:"" required:"" name:"pid" help:"Product ID"`
}

func (g *GenCmd) Run(ctx *kong.Context) error {
	if g.Kp {
		logger.Infof("creating new key pair for product: %s", g.Pid)
		_, err := lm.CreateProductKeyPair(g.Pid)
		if err != nil {
			fmt.Println("❌  Failed")
			return err
		}

		fmt.Println("✔️  Success")
		fmt.Println("Created new ed25519 key pair for product: ", g.Pid)
	}

	if g.Lic {
		logger.Infof("creating new license for product: %s", g.Pid)
		lic, key, err := lm.GenerateLicense(g.Pid)
		if err != nil {
			fmt.Println("❌  Failed")
			return err
		}

		fmt.Println("✔️  Success")
		fmt.Printf("Generated new license <%s> for %s.\nProduct key:\n%s\n", lic.ExternalId.String(), g.Pid, key)
	}

	return nil
}

type IssCmd struct {
	Key     string `arg:"" required:"" name:"key" help:"Product Key"`
	Machine string `arg:"" optional:"" name:"machine" help:"Machine Binding"`
}

func (i *IssCmd) Run(ctx *kong.Context) error {
	logger.Infof("issuing a new license with product key: %s", i.Key)
	ok, license, err := lm.ValidateKey(i.Key)
	if err != nil || !ok {
		fmt.Println("❌  Failed")
		return service.If(err != nil, err, fmt.Errorf("invalid product key: %+v", i.Key))
	}

	res, err := lm.AssignLicense(i.Machine, license)
	if err != nil {
		fmt.Println("❌  Failed")
		return err
	}

	fmt.Println("✔️  Success")
	fmt.Println("Issued license with product key: ", i.Key)
	resJson, err := jsoniter.MarshalToString(&res)
	if err != nil {
		fmt.Println(res)
		return nil
	}

	fmt.Println(resJson)

	return nil
}

type RevCmd struct {
	Id string `arg:"" required:"" name:"id" help:"ID of the license to revoke"`
}

func (r *RevCmd) Run() error {
	logger.Infof("revoking license: %s", r.Id)
	uid, err := uuid.Parse(r.Id)
	if err != nil {
		fmt.Println("❌  Failed")
		return err
	}

	if err := lm.RevokeLicense(uid); err != nil {
		fmt.Println("❌  Failed")
		return err
	}

	fmt.Println("✔️  Success")
	fmt.Println("Revoked license: ", r.Id)

	return nil
}

var CLI struct {
	Gen GenCmd `cmd:"" help:"create a new key pair or license"`
	Iss IssCmd `cmd:"" help:"Validate and assign a new license" `
	Rev RevCmd `cmd:"" help:"Revoke an existing license"`
}

func main() {
	ctx := kong.Parse(&CLI, kong.Name("license-manager"), kong.Description("A CLI for managing software licenses"), kong.UsageOnError())
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
