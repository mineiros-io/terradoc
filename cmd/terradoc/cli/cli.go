package cli

var Cli struct {
	Generate GenerateCmd `cmd:"" help:"Generate a markdown file from .tfdoc.hcl input."`
	Format   FormatCmd   `name:"fmt" cmd:"" help:"Format .tfdoc.hcl file."`
	Validate ValidateCmd `name:"validate" cmd:"" help:"Check if .tfdoc.hcl file is synchronized with Terraform variables and/or outputs files"`
}
