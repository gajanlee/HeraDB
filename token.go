package HeraDB

import "errors"

// token.go is used to tokenize the input sentence.

const (
    _TK_CREATE  = "create"
    _TK_FROM    = "from"
    _TK_SELECT  = "select"
    _TK_TABLE   = "table"

    _TK_COMMA   = ","
    _TK_STAR    = "*"
)

var (
    _TOKEN_ERROR    = errors.New("hahah")
)



func expect() string {
    return ""
}

func getNextToken() string {
    return ""
}
