import React, { useEffect } from "react"
import { useParams } from "react-router-dom"
import { ACTION_TYPES, useApiRequest } from "../apis/reducer"
import { scanSession } from "../apis/SessionAPIs"
import { getStoreInfoWithSSE } from "../apis/StoreAPIs"

const CreateCustomers = () => {
  let { storeId , sessionId}: {storeId: string, sessionId: string} = useParams()

  const [scanSessionAction, makeScanSessionRequest] = useApiRequest(...scanSession(sessionId, parseInt(storeId)))
  
  useEffect(() => {
    makeScanSessionRequest()
  }, [])

  useEffect(() => {
    if (scanSessionAction.actionType === ACTION_TYPES.ERROR) {
        // TODO: disable create customers button
    }

    // 40007: store_session exist but is already scanned.
    if ((scanSessionAction.response != null) && (scanSessionAction.response["error_code"]) && (scanSessionAction.response["error_code"] !== 40007)) {
        // TODO: disable create customers button
    }
  }, [scanSessionAction.actionType])

  useEffect(() => {
    let getStoreInfoSSE: EventSource
    getStoreInfoSSE = getStoreInfoWithSSE(parseInt(storeId))

    getStoreInfoSSE.onmessage = (event) => {
        // TODO: render customers ui
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
  
  return (
    <div>
        <div>
            scanned qrcode and create customers
        </div>


    </div>
  )
}

export {
    CreateCustomers
}