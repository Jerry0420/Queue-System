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

const checkExistenceOfRefreshableCookie = (): boolean => {
    const splittedCookie = document.cookie.split("=")
    if ((splittedCookie.length === 2) && (splittedCookie[0] === "refreshable")) {
        return true
    }
    return false
}

const checkUpdatableOfNormalToken = (jsonResponse: JSONResponse): boolean => {
    const tokenExpiresAt = (jsonResponse["token_expires_at"] as number)
    const now = Date.now()
    if (
        (tokenExpiresAt < now) || 
        (tokenExpiresAt >= now) && 
        (tokenExpiresAt - now <= 5*60)
        ) {
        return true // need to update normal token
    }
    return false // no eed to update normal token
}

// checkAuthFlow(refreshTokenAction.response, makeRefreshTokenRequest, 
//     () => {
        
//     }, () => {

//     }
//   )

const checkAuthFlow = (
    jsonResponse: JSONResponse | undefined | null, 
    makeRefreshTokenRequest: (() => Promise<JSONResponse | undefined>), 
    nextStuff: (() => void),
    redirectToMainPage: (() => void)
    ) => {
    
    const doMakeRefreshTokenRequest = () => {
        makeRefreshTokenRequest()
        .then(jsonResponse => {
            if (validateResponseSuccess(jsonResponse) === true) {
                nextStuff()   
            } else {
                redirectToMainPage()
            }
        })
        .catch((error) => {
            redirectToMainPage()
        })
    }
    
    if (checkExistenceOfRefreshableCookie() === true) {
        if (validateResponseSuccess(jsonResponse) === false) {
            // refresh cookie exist
            // refreshTokenAction.response not exist.
            doMakeRefreshTokenRequest()
        } else {
            const _jsonResponse = (jsonResponse as JSONResponse) // jsonResponse would never be null or undefined here
            if (checkUpdatableOfNormalToken(_jsonResponse) === true) {
                // refresh cookie, refreshTokenAction.response exist
                // normal token already need to be updated
                doMakeRefreshTokenRequest()
            } else {
                // refresh cookie, refreshTokenAction.response exist
                // normal token no need to be updated
                nextStuff()
            }
        }
    } else {
        // refreshable cookie not exist
        redirectToMainPage()
    }
}

export {
    validateResponseSuccess,
    checkExistenceOfRefreshableCookie,
    checkUpdatableOfNormalToken,
    checkAuthFlow
}