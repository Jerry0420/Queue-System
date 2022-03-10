import React, {useEffect, useContext, useState} from "react"
import { useParams, Link, useNavigate } from "react-router-dom"
import { RefreshTokenContext } from "./contexts"
import { createSessionWithSSE } from "../apis/SessionAPIs"
import { validateResponseSuccess } from "../apis/helper"
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import { toDataURL } from "qrcode"
import { updateStoreDescription } from "../apis/StoreAPIs"

const Store = () => {
  let { storeId: storeId }: {storeId: string} = useParams()
  let navigate = useNavigate()
  const [sessionScannedURL, setSessionScannedURL] = useState("")
  const [qrcodeImageURL, setQrcodeImageURL] = useState("")
  const [storeDescription, setStoreDescription] = useState("")
  const [normalToken, setNormalToken] = useState("")

  const {refreshTokenAction, makeRefreshTokenRequest, wrapCheckAuthFlow} = useContext(RefreshTokenContext)
  const [updateStoreDescriptionAction, makeUpdateStoreDescriptionRequest] = useApiRequest(
    ...updateStoreDescription(parseInt(storeId), normalToken, storeDescription)
    )

  const handleInputStoreDescription = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value: value }: { value: string } = e.target
    setStoreDescription(value)
  }

  useEffect(() => {
    let createSessionSSE: EventSource
    wrapCheckAuthFlow(
      () => {
        const sessionToken: string = ((refreshTokenAction.response as JSONResponse)["session_token"] as string)
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
    if (validateResponseSuccess(refreshTokenAction.response) === true) {
      const response: JSONResponse = refreshTokenAction.response as JSONResponse // refreshTokenAction.response must be JSONResponse here.
      setNormalToken(response["token"] as string)
    }
  }, [refreshTokenAction.response])

  useEffect(() => {
    // handle running, success, error states here.
  }, [updateStoreDescriptionAction])

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