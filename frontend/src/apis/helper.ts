import { JSONResponse, ACTION_TYPES, Action } from "./reducer"
import { useContext } from "react"
import { RefreshTokenContext } from "../components/contexts"

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

const checkAuthFlow = ( 
    jsonResponse: JSONResponse | null | undefined,
    makeRefreshTokenRequest: (() => Promise<JSONResponse | null | undefined>),
    nextStuff: ((jsonResponse: JSONResponse) => void),
    redirectToMainPage: (() => void)
    ) => {
    
    // refresh cookie, refreshTokenAction.response exist
    const doMakeRefreshTokenRequest = () => {
        makeRefreshTokenRequest()
        .then(response => {
            if (validateResponseSuccess(response) === true) {
                console.log("1=========")
                nextStuff((response as JSONResponse))   
            } else {
                console.log("2=========")
                redirectToMainPage()
            }
        })
        .catch((error: Error) => {
            console.log("3=========")
            redirectToMainPage()
        })
    }
    
    if (checkExistenceOfRefreshableCookie() === true) {
        if (validateResponseSuccess(jsonResponse) === false) {
            // refresh cookie exist
            // refreshTokenAction.response not exist.
            console.log("4=========")
            doMakeRefreshTokenRequest()
        } else {
            const _jsonResponse = (jsonResponse as JSONResponse) // jsonResponse would never be null or undefined here
            if (checkUpdatableOfNormalToken(_jsonResponse) === true) {
                // refresh cookie, refreshTokenAction.response exist
                // normal token already need to be updated
                console.log("5=========")
                doMakeRefreshTokenRequest()
            } else {
                // refresh cookie, refreshTokenAction.response exist
                // normal token no need to be updated
                console.log("6=========")
                nextStuff((jsonResponse as JSONResponse))
            }
        }
    } else {
        // refreshable cookie not exist
        console.log("7=========")
        redirectToMainPage()
    }
}

export {
    validateResponseSuccess,
    checkExistenceOfRefreshableCookie,
    checkUpdatableOfNormalToken,
    checkAuthFlow
}