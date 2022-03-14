import React, {useEffect, useContext, useState} from "react"
import { useParams, Link, useNavigate } from "react-router-dom"
import { RefreshTokenContext } from "./contexts"
import { createSessionWithSSE } from "../apis/SessionAPIs"
import { validateResponseSuccess } from "../apis/helper"
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import { toDataURL } from "qrcode"
import { getStoreInfoWithSSE, updateStoreDescription } from "../apis/StoreAPIs"
import { getNormalTokenFromRefreshTokenAction, getSessionTokenFromRefreshTokenAction } from "../apis/validator"

const Store = () => {
  let { storeId }: {storeId: string} = useParams()
  let navigate = useNavigate()
  const [sessionScannedURL, setSessionScannedURL] = useState("")
  const [qrcodeImageURL, setQrcodeImageURL] = useState("")
  const [storeDescription, setStoreDescription] = useState("")

  const {refreshTokenAction, makeRefreshTokenRequest, wrapCheckAuthFlow} = useContext(RefreshTokenContext)
  const [updateStoreDescriptionAction, makeUpdateStoreDescriptionRequest] = useApiRequest(
    ...updateStoreDescription(
      parseInt(storeId), 
      getNormalTokenFromRefreshTokenAction(refreshTokenAction.response), 
      storeDescription
      )
  )

  const handleInputStoreDescription = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value: value }: { value: string } = e.target
    setStoreDescription(value)
  }

  useEffect(() => {
    let getStoreInfoSSE: EventSource
    getStoreInfoSSE = getStoreInfoWithSSE(parseInt(storeId))

    getStoreInfoSSE.onmessage = (event) => {
      // TODO: render to ui
      console.log(JSON.parse(event.data))
    }
    
    getStoreInfoSSE.onerror = (event) => {
      getStoreInfoSSE.close()
    }
    return () => {
      if (getStoreInfoSSE != null) {
        getStoreInfoSSE.close()
      }
    }
  }, [getStoreInfoWithSSE])

  useEffect(() => {
    let createSessionSSE: EventSource
    wrapCheckAuthFlow(
      () => {
        const sessionToken: string = getSessionTokenFromRefreshTokenAction(refreshTokenAction.response)
        createSessionSSE = createSessionWithSSE(sessionToken)

        createSessionSSE.onmessage = (event) => {
          setSessionScannedURL(JSON.parse(event.data)["scanned_url"])
        }
        
        createSessionSSE.onerror = (event) => {
          createSessionSSE.close()
        }
      },
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

  const doMakeUpdateStoreDescriptionRequest = () => {
    wrapCheckAuthFlow(
      () => {
        makeUpdateStoreDescriptionRequest()
      },
      () => {
         // TODO: show error message
         navigate("/")
      }
    )
  }

  useEffect(() => {
    toDataURL(sessionScannedURL, (error, url) => {
      if (url != null) {
        setQrcodeImageURL(url)
      }
    })
  }, [sessionScannedURL])

  useEffect(() => {
    // TODO: handle running, success, error states here.
  }, [updateStoreDescriptionAction.actionType])

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

        <br />
        <a href={sessionScannedURL} target="_blank">{sessionScannedURL}</a>
    </div>
  )
}

export {
  Store
}