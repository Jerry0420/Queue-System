package domain

type ServerError struct {
	Code        int
	Description string
}

func (serverError *ServerError) Error() string {
	return serverError.Description
}

var (
	// list of all custom errors.
	// check README file for detailed description.

	// ============================================================

	// lack of required params
	ServerError40001 = &ServerError{Code: 40001, Description: "lack of required params"}
	// length of password is not appropriate.
	ServerError40002 = &ServerError{Code: 40002, Description: "length of password is not appropriate."}
	// the incoming password is not equal to the original password.
	ServerError40003 = &ServerError{Code: 40003, Description: "password mismatch."}
	// wrong params
	ServerError40004 = &ServerError{Code: 40004, Description: "wrong params"}

	// ============================================================

	// fail to parse jwt token
	ServerError40101 = &ServerError{Code: 40101, Description: "fail to parse jwt token"}
	// lack of jwt token
	ServerError40102 = &ServerError{Code: 40102, Description: "lack of jwt token"}

	// ============================================================

	// unsupported url route
	ServerError40401 = &ServerError{Code: 40401, Description: "unsupported url route"}
	// store not exist.
	ServerError40402 = &ServerError{Code: 40402, Description: "store not exist."}
	// sign_key not exist.
	ServerError40403 = &ServerError{Code: 40403, Description: "sign_key not exist."}

	// ============================================================

	// method not allowed.
	ServerError40501 = &ServerError{Code: 40501, Description: "method not allowed."}

	// ============================================================

	// store already exist.
	ServerError40901 = &ServerError{Code: 40901, Description: "store already exist."}
	// sign_key not exist.
	ServerError40902 = &ServerError{Code: 40902, Description: "sign_key not exist."}

	// ============================================================

	// other internal server error
	ServerError50001 = &ServerError{Code: 50001, Description: "other internal server error"}
	// unexpected database error
	ServerError50002 = &ServerError{Code: 50002, Description: "unexpected database error"}
)
