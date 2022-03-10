import React, {useEffect, useContext} from "react"
import { useParams, Link } from "react-router-dom"
import { RefreshTokenContext } from "./contexts"
import { createSessionWithSSE } from "../apis/SessionAPIs"
import { checkAuthFlow, validateResponseSuccess } from "../apis/helper"
import { JSONResponse } from "../apis/reducer"

const Store = () => {
  let { storeId } = useParams()
  const {refreshTokenAction, makeRefreshTokenRequest} = useContext(RefreshTokenContext)

  useEffect(() => {
    checkAuthFlow(refreshTokenAction.response, makeRefreshTokenRequest, 
      // nextStuff
      () => {
        if (validateResponseSuccess(refreshTokenAction.response) === true) {
          const sessionToken: string = ((refreshTokenAction.response as JSONResponse)["session_token"] as string)
          const createSessionSSE = createSessionWithSSE(sessionToken)

          createSessionSSE.onmessage = (event) => {
            console.log(JSON.parse(event.data))
          }
          
          createSessionSSE.onopen = (event) => {
            console.log('on open', event)
          }
          
          createSessionSSE.onerror = (event) => {
            console.log('on error', event)
          }
        }
      }, 
      // redirectToMainPage
      () => {
        console.log("in redirectToMainPage", refreshTokenAction)
      }
    )
  }, [createSessionWithSSE, refreshTokenAction.response, refreshTokenAction.exception])

  return (
    <div>
        <Link to="/temp">to temp</Link>
        {/* {console.log(refreshTokenAction)} */}
    </div>
  )
}

export {
  Store
}