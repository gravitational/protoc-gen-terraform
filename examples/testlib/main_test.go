package testlib

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/suite"

	"github.com/gravitational/protoc-gen-terraform/v3/examples/testlib/provider"
)

func TestTerraform(t *testing.T) {
	suite.Run(t, &TerraformSuite{})
}

type TerraformSuite struct {
	suite.Suite

	terraformProvider tfsdk.Provider

	terraformProviders map[string]func() (tfprotov6.ProviderServer, error)
}

func (s *TerraformSuite) SetupSuite() {
	os.Setenv("TF_ACC", "true")
	s.terraformProvider = provider.New()
	s.terraformProviders = make(map[string]func() (tfprotov6.ProviderServer, error))
	s.terraformProviders["example"] = func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6(s.terraformProvider)(), nil
	}
}

//go:embed fixtures/*
var fixtures embed.FS

// getFixture loads fixture and returns it as string or <error> if failed
func (s *TerraformSuite) getFixture(name string, formatArgs ...any) string {
	b, err := fixtures.ReadFile(filepath.Join("fixtures", name))
	if err != nil {
		return fmt.Sprintf("<error: %v fixture not found>", name)
	}
	return fmt.Sprintf(string(b), formatArgs...)
}
