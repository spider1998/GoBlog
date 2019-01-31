package code

import "sort"

type Code int

//响应码
const (
	CodeServerErr    Code = 1
	CodeBadRequest   Code = 2
	CodeNoPermission Code = 3

	CodeUserExist                Code = 10000
	CodeUserNotExist             Code = 10001
	CodeUserInvalidPassword      Code = 10002
	CodeUserAccessSessionInvalid Code = 10003
	CodeVerifyError                   = 10004
	CodeArticleExist                  = 10005
	CodeArticleNotExist               = 10006
	CodeArticleNotChange              = 10007
	CodeDenied						  =	10008
	CodeContentNotExist				  = 10009
	CodeIDNotAllowed				  =	100010

	CodeRecordNotExist Code = 10803
)

//响应码对应信息集合
var codeMap = map[Code]string{
	CodeServerErr:    "internal server error.",
	CodeBadRequest:   "invalid request.",
	CodeNoPermission: "operation is not allowed.",

	CodeUserExist:                "account already exists.",
	CodeRecordNotExist:           "record not exist.",
	CodeUserNotExist:             "user not exist",
	CodeUserInvalidPassword:      "password error",
	CodeUserAccessSessionInvalid: "Session error",
	CodeVerifyError:              "Verification code error",

	CodeArticleExist:     "article exist",
	CodeArticleNotExist:  "article not exist",
	CodeArticleNotChange: "article not change any more",
	CodeDenied:			"Permission denied",
	CodeContentNotExist: "Content not exist",
	CodeIDNotAllowed	: "ID Not Allowed",
}

//解析码的具体信息
func parseCodeMessage(code Code) string {
	if msg, ok := codeMap[code]; ok {
		return msg
	}
	return codeMap[CodeServerErr]
}

//所有响应码列表
func ListCode() []map[string]interface{} {
	var codes []int
	for k := range codeMap {
		codes = append(codes, int(k))
	}
	sort.Ints(codes)

	var list []map[string]interface{}
	for _, k := range codes {
		list = append(list, map[string]interface{}{"code": k, "message": codeMap[Code(k)]})
	}

	return list
}
