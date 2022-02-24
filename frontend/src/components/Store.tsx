import React from "react"
import { openStore } from "../apis/StoreAPIs"

const Store = () => {
  return (
    <button onClick={() => openStore("email", "YXRlbjEyMzQ=", "name", "Asia/Taipei", ["queue_a", "queue_b"])}>
      openStore
    </button>
  )
}

export default Store