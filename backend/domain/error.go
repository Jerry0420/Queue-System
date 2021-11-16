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
	// other jwt token parse error
	ServerError40103 = &ServerError{Code: 40103, Description: "other jwt token parse error"}
	// jwt token expired
	ServerError40104 = &ServerError{Code: 40104, Description: "jwt token expired"}
	// store( refresh jwt token) expired
	ServerError40105 = &ServerError{Code: 40105, Description: "store( refresh jwt token) expired"}
	
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

	// store already exist. (not exceed 24 hrs)
	ServerError40901 = &ServerError{Code: 40901, Description: "store already exist. (not exceed 24 hrs)"}
	// sign_key already exist.
	ServerError40902 = &ServerError{Code: 40902, Description: "sign_key already exist."}
	// store already exist. (but exceed 24 hrs, and already auto close this store.)
	ServerError40903 = &ServerError{Code: 40903, Description: "store already exist. (but exceed 24 hrs, and already auto close this store.)"}

	// ============================================================

	// other internal server error
	ServerError50001 = &ServerError{Code: 50001, Description: "other internal server error"}
	// unexpected database error
	ServerError50002 = &ServerError{Code: 50002, Description: "unexpected database error"}
)
