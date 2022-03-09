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
    refreshTokenAction: Action,
    makeRefreshTokenRequest: (() => Promise<boolean>),
    nextStuff: (() => void),
    redirectToMainPage: (() => void)
    ) => {
    
    // refresh cookie, refreshTokenAction.response exist
    const doMakeRefreshTokenRequest = () => {
        makeRefreshTokenRequest()
        .then(done => {
            if (refreshTokenAction.actionType === ACTION_TYPES.SUCCESS) {
                if (validateResponseSuccess(refreshTokenAction.response) === true) {
                    console.log("1=========")
                    nextStuff()   
                } else {
                    console.log("2=========")
                    redirectToMainPage()
                }
            }
            if (refreshTokenAction.actionType === ACTION_TYPES.ERROR) {
                console.log("3=========")
                redirectToMainPage()
            }
        })
    }
    
    if (checkExistenceOfRefreshableCookie() === true) {
        if (validateResponseSuccess(refreshTokenAction.response) === false) {
            // refresh cookie exist
            // refreshTokenAction.response not exist.
            console.log("4=========")
            doMakeRefreshTokenRequest()
        } else {
            const jsonResponse = (refreshTokenAction.response as JSONResponse) // jsonResponse would never be null or undefined here
            if (checkUpdatableOfNormalToken(jsonResponse) === true) {
                // refresh cookie, refreshTokenAction.response exist
                // normal token already need to be updated
                console.log("5=========")
                doMakeRefreshTokenRequest()
            } else {
                // refresh cookie, refreshTokenAction.response exist
                // normal token no need to be updated
                console.log("6=========")
                nextStuff()
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