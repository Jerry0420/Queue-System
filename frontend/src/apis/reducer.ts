import { useReducer, useCallback } from "react"
import {RequestParams} from './base'

const ACTION_TYPES = {
    RUNNING: 'running',
    SUCCESS: 'success',
    ERROR: 'error',
}

const doRunning = (): object => ({ actionType: ACTION_TYPES.RUNNING })
const doSuccess = (response: object): object => ({ actionType: ACTION_TYPES.SUCCESS, response })
const doError = (exception: object): object => ({ actionType: ACTION_TYPES.ERROR, exception })

const initialState = {
    actionType: null,
    response: null,
    exception: null
}

interface Action {
    actionType: string
    response?: object
    exception?: Error
}
  
const reducer = (state = initialState, { actionType, response, exception }: Action) => {
    switch (actionType) {
      case ACTION_TYPES.RUNNING:
        return { ...initialState, actionType: ACTION_TYPES.RUNNING }
      case ACTION_TYPES.SUCCESS:
        return { ...state, actionType: ACTION_TYPES.SUCCESS, response }
      case ACTION_TYPES.ERROR:
        return { ...state, actionType: ACTION_TYPES.ERROR, response, exception }
      default:
        return state;
    }
}

const useApiRequest = (endpoint: string, requestParams: RequestParams): [Action, () => Promise<any>] => {
    const [action, dispatch] = useReducer(reducer, initialState)
  
    const makeRequest = useCallback(async () => {
      dispatch(doRunning())

      try {
        const response = await fetch(endpoint, requestParams)
          .then(response => response.json())
          .then(jsonResponse => {
              return jsonResponse
          })
          .catch(error => {
              console.error(error)
              throw error  
          })

        dispatch(doSuccess(response))
      
      } catch (error) {
        
        dispatch(doError(error));
      
      }
    }, [endpoint, requestParams]);
  
    return [action, makeRequest];
  }

export {
    ACTION_TYPES,
    doRunning,
    doSuccess,
    doError,
    initialState,
    reducer,
    useApiRequest
}