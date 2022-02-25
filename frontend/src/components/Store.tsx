import React from "react"
import * as storeAPIs from "../apis/StoreAPIs"

const Store = () => {
  return (
    <>
      <button onClick={() => storeAPIs.openStore("email", "YXRlbjEyMzQ=", "name", "Asia/Taipei", ["queue_a", "queue_b"])}>
        openStore
      </button>
      <br />
      <button onClick={() => storeAPIs.signInStore("email", "YXRlbjEyMzQ=")}>
        signInStore
      </button>
      <br />
      <button onClick={() => storeAPIs.refreshToken()}>
        refreshToken
      </button>
      <br />
      <button onClick={() => storeAPIs.closeStore(8, "xxxxx")}>
        closeStore
      </button>
    </>
  )
}

export default Store