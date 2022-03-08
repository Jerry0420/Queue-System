import {createContext} from 'react'
import { initialState, Action } from '../apis/reducer'

const initialRefreshTokenContext: {
    refreshTokenAction: Action, 
    makeRefreshTokenRequest: (() => Promise<void>) | undefined
} = {
    refreshTokenAction: initialState,
    makeRefreshTokenRequest: undefined
}

const RefreshTokenContext = createContext(initialRefreshTokenContext)

export {
    RefreshTokenContext
}