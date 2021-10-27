package types

// TerraformType represents Terraform type as specified in the `type` attribute of a variable
type TerraformType int

//go:generate stringer -type=TerraformType -linecomment
const (
	TerraformEmptyType TerraformType = iota // empty
	TerraformBool                           // bool
	TerraformString                         // string
	TerraformNumber                         // number
	TerraformAny                            // any
	TerraformList                           // list
	TerraformSet                            // set
	TerraformMap                            // map
	TerraformObject                         // object
	TerraformTuple                          // tuple
)

var TerraformTypes = map[string]TerraformType{
	TerraformBool.String():   TerraformBool,
	TerraformString.String(): TerraformString,
	TerraformNumber.String(): TerraformNumber,
	TerraformAny.String():    TerraformAny,
	TerraformList.String():   TerraformList,
	TerraformSet.String():    TerraformSet,
	TerraformMap.String():    TerraformMap,
	TerraformObject.String(): TerraformObject,
	TerraformTuple.String():  TerraformTuple,
}
