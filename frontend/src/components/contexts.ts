import {createContext} from 'react'
import { initialState, Action } from '../apis/reducer'
import { JSONResponse } from '../apis/reducer'

const initialRefreshTokenContext: {
    refreshTokenAction: Action, 
    makeRefreshTokenRequest: (() => Promise<boolean>)
} = {
    refreshTokenAction: initialState,
    makeRefreshTokenRequest: (() => {return new Promise((resolve, reject) => {})})
}

// const {refreshTokenAction, makeRefreshTokenRequest} = useContext(RefreshTokenContext)
const RefreshTokenContext = createContext(initialRefreshTokenContext)

export {
    RefreshTokenContext
}