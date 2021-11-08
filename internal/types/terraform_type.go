package types

// TerraformType represents Terraform type as specified in the `type` attribute of a variable
type TerraformType int

//go:generate stringer -type=TerraformType -linecomment
const (
	TerraformInvalidType TerraformType = iota // invalid
	TerraformBool                             // bool
	TerraformString                           // string
	TerraformNumber                           // number
	TerraformDynamic                          // dynamic
	TerraformList                             // list
	TerraformObject                           // object
	TerraformTuple                            // tuple
)

var TerraformTypes = map[string]TerraformType{
	TerraformBool.String():    TerraformBool,
	TerraformString.String():  TerraformString,
	TerraformNumber.String():  TerraformNumber,
	TerraformDynamic.String(): TerraformDynamic,
	TerraformList.String():    TerraformList,
	TerraformObject.String():  TerraformObject,
	TerraformTuple.String():   TerraformTuple,
}
