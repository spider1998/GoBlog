package code

import "sort"

type Code int

//响应码
const (
	CodeServerErr                Code = 1
	CodeBadRequest                    = 2
	CodeNoPermission                  = 3
	CodeUserExist                     = 10000
	CodeUserNotExist                  = 10001
	CodeUserInvalidPassword           = 10002
	CodeUserAccessSessionInvalid      = 10003
	CodeVerifyError                   = 10004
	CodeArticleExist                  = 10005
	CodeArticleNotExist               = 10006
	CodeArticleNotChange              = 10007
	CodeDenied                        = 10008
	CodeContentNotExist               = 10009
	CodeIDNotAllowed                  = 100010
	CodeInvalidCaptcha                = 100011
	CodeInvalidData                   = 100012
	CodeOperatorNotExist              = 100013
	CodeUserDisabled                  = 100014
	CodeTokenNotExist                 = 100015
	CodeStateInvalid                  = 100016
	CodeOperatorTokenRequired         = 100017
	CodeCommentExist                  = 100018
	CodeCommentNotExist               = 100019
	CodeMessageNotExist               = 100020
	CodeCreateMessageError            = 100021
	CodeTaskIsInProgress              = 100022
	CodeStatisticExist                = 100023
	CodeFriendExist                   = 100024
	CodeRecordNotExist                = 10803
)

//响应码对应信息集合
var codeMap = map[Code]string{
	CodeServerErr:                "Internal server error.",
	CodeBadRequest:               "Invalid request.",
	CodeNoPermission:             "Operation is not allowed.",
	CodeUserExist:                "Account already exists.",
	CodeRecordNotExist:           "Record not exist.",
	CodeUserNotExist:             "User not exist",
	CodeUserInvalidPassword:      "password error",
	CodeUserAccessSessionInvalid: "Session error",
	CodeVerifyError:              "Verification code error",
	CodeArticleExist:             "Article exist",
	CodeArticleNotExist:          "Article not exist",
	CodeArticleNotChange:         "Article not change any more",
	CodeDenied:                   "Permission denied",
	CodeContentNotExist:          "Content not exist",
	CodeIDNotAllowed:             "ID Not Allowed",
	CodeInvalidCaptcha:           "Invalid Captcha",
	CodeInvalidData:              "Invalid Data",
	CodeOperatorNotExist:         "Operator Not Exist",
	CodeUserDisabled:             "User Disabled",
	CodeTokenNotExist:            "Token Not Exist",
	CodeStateInvalid:             "State Invalid",
	CodeOperatorTokenRequired:    "Operator Token Required",
	CodeCommentExist:             "Comment Exist",
	CodeCommentNotExist:          "Comment Not Exist",
	CodeMessageNotExist:    "Message Not Exist",
	CodeCreateMessageError: "Create Message Error",
	CodeTaskIsInProgress:   "Task Is In Progress",
	CodeStatisticExist:           "Statistic Exist",
	CodeFriendExist: "Friend Exist",
}

func parseCodeMessage(code Code) string {
	if msg, ok := codeMap[code]; ok {
		return msg
	}
	return codeMap[CodeServerErr]
}

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
