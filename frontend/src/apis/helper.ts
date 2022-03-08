import { JSONResponse } from "./reducer"

const validateResponseSuccess = (jsonResponse: JSONResponse | null | undefined): boolean => {
    if (
        (jsonResponse != null) && 
        (jsonResponse["error_code"] === undefined)
        ) {
        return true
    } else {
        return false
    }
}

export {
    validateResponseSuccess
}