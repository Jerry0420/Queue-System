import React, {useState} from "react"
import { useApiRequest } from "../apis/reducer"
import * as storeAPIs from "../apis/StoreAPIs"

const Store = () => {
  const [normalToken, setNormalToken] = useState("")
  const [storeId, setStoreId] = useState(0)

  const handleInputNormalToken = (e: React.ChangeEvent<HTMLElement>) => {
    const { value }: { value: string } = e.target
    setNormalToken(value)
  }

  const handleInputStoreId = (e: React.ChangeEvent<HTMLElement>) => {
    const { value }: { value: number } = e.target
    setStoreId(value)
  }

  const getCookie = () => {
    console.log(document.cookie)
  }
  
  const [openStoreAction, makeOpenStoreRequest] = useApiRequest(...storeAPIs.openStore("jeerywa@gmail.com", "YXRlbjEyMzQ=", "name", "Asia/Taipei", ["queue_a", "queue_b"]))
  const [signInStoreAction, makeSignInStoreRequest] = useApiRequest(...storeAPIs.signInStore("jeerywa@gmail.com", "YXRlbjEyMzQ="))
  const [refreshTokenAction, makeRefreshTokenRequest] = useApiRequest(...storeAPIs.refreshToken())
  const [closeStoreAction, makeCloseStoreRequest] = useApiRequest(...storeAPIs.closeStore(storeId, normalToken))

  return (
    <>
      <button onClick={makeOpenStoreRequest}>
        openStore
      </button>
      <>{console.log(openStoreAction)}</>
      
      <br />
      <button onClick={makeSignInStoreRequest}>
        signInStore
      </button>
      <>{console.log(signInStoreAction)}</>
      
      <br />
      <button onClick={makeRefreshTokenRequest}>
        refreshToken
      </button>
      <>{console.log(refreshTokenAction)}</>
      
      <br />
      <input
          type="text"
          onChange={handleInputStoreId}
          placeholder="storeId"
        />
      <input
          type="text"
          onChange={handleInputNormalToken}
          placeholder="normalToken"
        />
      <button onClick={makeCloseStoreRequest}>
        closeStore
      </button>
      <>{console.log(closeStoreAction)}</>
      
      <hr />

      <br />
      <button onClick={getCookie}>
        get refresh token (cookie)
      </button>
    </>
  )
}

export default Store