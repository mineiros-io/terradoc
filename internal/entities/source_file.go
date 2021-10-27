package entities

import "github.com/hashicorp/hcl/v2"

// TODO: decouple from HCL

type SourceFile struct {
	HCLFile *hcl.File
}

func (sf *SourceFile) HCLBody() hcl.Body {
	return sf.HCLFile.Body
}
