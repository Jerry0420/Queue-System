import React, {useState} from "react"
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

  return (
    <>
      <button onClick={() => storeAPIs.openStore("jeerywa@gmail.com", "YXRlbjEyMzQ=", "name", "Asia/Taipei", ["queue_a", "queue_b"])}>
        openStore
      </button>
      <br />
      <button onClick={() => storeAPIs.signInStore("jeerywa@gmail.com", "YXRlbjEyMzQ=")}>
        signInStore
      </button>
      <br />
      <button onClick={() => storeAPIs.refreshToken()}>
        refreshToken
      </button>
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
      <button onClick={() => storeAPIs.closeStore(storeId, normalToken)}>
        closeStore
      </button>

      <br />
      <button onClick={getCookie}>
        get refresh token (cookie)
      </button>
    </>
  )
}

export default Store