package types

// TerraformType represents Terraform type as specified in the `type` attribute of a variable
type TerraformType int

//go:generate stringer -type=TerraformType -linecomment
const (
	TerraformEmptyType TerraformType = iota // empty
	TerraformBool                           // bool
	TerraformString                         // string
	TerraformNumber                         // number
	TerraformList                           // list
	TerraformSet                            // set
	TerraformMap                            // map
	TerraformObject                         // object
	TerraformTuple                          // tuple
	TerraformAny                            // any
	TerraformResource                       // resource
	TerraformModule                         // module
)

var SupportedTerraformTypes = []TerraformType{
	TerraformBool,
	TerraformString,
	TerraformNumber,
	TerraformList,
	TerraformSet,
	TerraformMap,
	TerraformObject,
	TerraformTuple,
	TerraformAny,
	TerraformResource,
	TerraformModule,
}

func TerraformTypes(typename string) (TerraformType, bool) {
	switch typename {
	case TerraformBool.String():
		return TerraformBool, true
	case TerraformString.String():
		return TerraformString, true
	case TerraformNumber.String():
		return TerraformNumber, true
	case TerraformList.String():
		return TerraformList, true
	case TerraformSet.String():
		return TerraformSet, true
	case TerraformMap.String():
		return TerraformMap, true
	case TerraformObject.String():
	case "any": // TODO: check if this make sense as it was included for enabling type checking validation
		return TerraformObject, true
	case TerraformTuple.String():
		return TerraformTuple, true
	case TerraformResource.String():
		return TerraformResource, true
	}

	return TerraformEmptyType, false
}

func (t TerraformType) IsComplex() bool {
	return (t != TerraformEmptyType) &&
		((t == TerraformList) ||
			(t == TerraformSet) ||
			(t == TerraformMap) ||
			(t == TerraformObject) ||
			(t == TerraformResource) ||
			(t == TerraformModule))
}
