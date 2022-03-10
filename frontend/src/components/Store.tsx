import React, {useEffect, useContext, useState} from "react"
import { useParams, Link, useNavigate } from "react-router-dom"
import { RefreshTokenContext } from "./contexts"
import { createSessionWithSSE } from "../apis/SessionAPIs"
import { checkAuthFlow, validateResponseSuccess } from "../apis/helper"
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import { toDataURL } from "qrcode"
import { updateStoreDescription } from "../apis/StoreAPIs"

const Store = () => {
  let { storeId }: {storeId: string} = useParams()
  let navigate = useNavigate()
  const [sessionScannedURL, setSessionScannedURL] = useState("")
  const [qrcodeImageURL, setQrcodeImageURL] = useState("")
  const [storeDescription, setStoreDescription] = useState("")
  const [normalToken, setNormalToken] = useState("")

  const {refreshTokenAction, makeRefreshTokenRequest} = useContext(RefreshTokenContext)
  const [updateStoreDescriptionAction, makeUpdateStoreDescriptionRequest] = useApiRequest(
    ...updateStoreDescription(parseInt(storeId), normalToken, storeDescription)
    )

  const handleInputStoreDescription = (e: React.ChangeEvent<HTMLElement>) => {
    const { value }: { value: string } = e.target
    setStoreDescription(value)
  }

  const doCheckAuthFlow = (nextStuff: () => void) => {
    checkAuthFlow(refreshTokenAction.response, makeRefreshTokenRequest, 
      // nextStuff
      () => {
        if (validateResponseSuccess(refreshTokenAction.response) === true) {
          nextStuff()
        }
      }, 
      // redirectToMainPage
      () => {
        // TODO: show error message
        navigate("/")
      }
    )
  }

  const doMakeUpdateStoreDescriptionRequest = () => {
    doCheckAuthFlow(() => {
      makeUpdateStoreDescriptionRequest()
    })
  }

  useEffect(() => {
    toDataURL(sessionScannedURL, (error, url) => {
      if (url != null) {
        setQrcodeImageURL(url)
      }
    })
  }, [sessionScannedURL])

  useEffect(() => {
    if (validateResponseSuccess(refreshTokenAction.response) === true) {
      setNormalToken(refreshTokenAction.response["token"])
    }
  }, [refreshTokenAction.response])

  useEffect(() => {
    let createSessionSSE: EventSource
    doCheckAuthFlow(() => {
      const sessionToken: string = ((refreshTokenAction.response as JSONResponse)["session_token"] as string)
      createSessionSSE = createSessionWithSSE(sessionToken)

      createSessionSSE.onmessage = (event) => {
        setSessionScannedURL(JSON.parse(event.data)["scanned_url"])
      }
      
      createSessionSSE.onerror = (event) => {
        createSessionSSE.close()
      }
    })
    return () => {
      if (createSessionSSE != null) {
        createSessionSSE.close()
      }
    }
  }, [createSessionWithSSE, refreshTokenAction.response, refreshTokenAction.exception])

  return (
    <div>
        <Link to="/temp">to temp</Link>
        {/* <img src={qrcodeImageURL} alt="qrcode image"></img> */}

        <br />
        <input
          type="text"
          onChange={handleInputStoreDescription}
          placeholder="input store description"
        />
        <button onClick={doMakeUpdateStoreDescriptionRequest}>
          update store description
        </button>
    </div>
  )
}

export {
  Store
}