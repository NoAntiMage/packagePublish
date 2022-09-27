package constant

import "errors"

var (
	ErrBindJson           = errors.New("Param bind json failed")
	ErrParamIsNotComplete = errors.New("Param is not complete, please check post data")
	ErrValidateFail       = errors.New("validate fail")
	ErrRecordNotfound     = errors.New("Record not found")
	ErrRecordExist        = errors.New("Record exist")
	ErrParamInvalid       = errors.New("Param invalid")
	ErrParamNotFoundInDb  = errors.New("Param not found in db")
	ErrIpInvalid          = errors.New("ip invalid")
	ErrTypeAssertFail     = errors.New("TypeAssertFail")
)
