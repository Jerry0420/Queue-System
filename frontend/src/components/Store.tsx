import React, {useEffect, useContext, useState} from "react"
import { useParams, Link, useNavigate } from "react-router-dom"
import { RefreshTokenContext } from "./contexts"
import { createSessionWithSSE } from "../apis/SessionAPIs"
import { checkAuthFlow, validateResponseSuccess } from "../apis/helper"
import { JSONResponse } from "../apis/reducer"

const Store = () => {
  let { storeId } = useParams()
  let navigate = useNavigate()
  const {refreshTokenAction, makeRefreshTokenRequest} = useContext(RefreshTokenContext)
  const [sessionScannedURL, setSessionScannedURL] = useState(null)

  useEffect(() => {

    let createSessionSSE: EventSource
    checkAuthFlow(refreshTokenAction.response, makeRefreshTokenRequest, 
      // nextStuff
      () => {
        if (validateResponseSuccess(refreshTokenAction.response) === true) {
          const sessionToken: string = ((refreshTokenAction.response as JSONResponse)["session_token"] as string)
          createSessionSSE = createSessionWithSSE(sessionToken)

          createSessionSSE.onmessage = (event) => {
            setSessionScannedURL(JSON.parse(event.data)["scanned_url"])
          }
          
          createSessionSSE.onerror = (event) => {
            createSessionSSE.close()
          }
        }
      }, 
      // redirectToMainPage
      () => {
        // TODO: show error message
        navigate("/")
      }
    )

    return () => {
      if (createSessionSSE != null) {
        createSessionSSE.close()
      }
    }
  }, [createSessionWithSSE, refreshTokenAction.response, refreshTokenAction.exception])

  return (
    <div>
        <Link to="/temp">to temp</Link>
        {console.log(sessionScannedURL)}
        {/* {console.log(refreshTokenAction)} */}
    </div>
  )
}

export {
  Store
}