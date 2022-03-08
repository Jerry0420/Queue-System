import {createContext} from 'react'
import { initialState, Action } from '../apis/reducer'

const initialRefreshTokenContext: {
    refreshTokenAction: Action, 
    makeRefreshTokenRequest: (() => Promise<void>)
} = {
    refreshTokenAction: initialState,
    makeRefreshTokenRequest: (() => {return new Promise((resolve, reject) => {})})
}

const RefreshTokenContext = createContext(initialRefreshTokenContext)

export {
    RefreshTokenContext
}