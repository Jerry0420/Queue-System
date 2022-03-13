import React, {useState, useContext} from "react"
import { useNavigate } from "react-router-dom"
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import * as storeAPIs from "../apis/StoreAPIs"
import { RefreshTokenContext } from "./contexts"

// import PropTypes from 'prop-types'

// const Queue = ({name}: {name: string}) => {
//   return (
//     <div>
//       {name}
//     </div>
//   )
// }

// Queue.propTypes = {
//   name: PropTypes.string,
// }

// Queue.defaultProps = {
//   name: '',
// }

const Temp = () => {
  const [normalToken, setNormalToken] = useState("")
  const [storeId, setStoreId] = useState(0)

  const handleInputNormalToken = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setNormalToken(value)
  }

  const handleInputStoreId = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setStoreId(parseInt(value))
  }

  const getCookie = () => {
    console.log(document.cookie)
  }

  let navigate = useNavigate()

  const goToStorePage = () => {
    const storeId = localStorage.getItem("storeId")
    navigate(`/stores/${storeId}`)
  }

  const saveStoreId = (jsonResponse: JSONResponse | null | undefined) => {
      const _jsonResponse = (jsonResponse as JSONResponse)
      const storeId: number = (_jsonResponse["id"] as number)
      localStorage.setItem("storeId", storeId.toString())
  }
  
  const [openStoreAction, makeOpenStoreRequest] = useApiRequest(...storeAPIs.openStore("jeerywa@gmail.com", "YXRlbjEyMzQ=", "name", "Asia/Taipei", ["queue_a", "queue_b"]))
  const [signInStoreAction, makeSignInStoreRequest] = useApiRequest(...storeAPIs.signInStore("jeerywa@gmail.com", "YXRlbjEyMzQ="))
  // const [refreshTokenAction, makeRefreshTokenRequest] = useApiRequest(...storeAPIs.refreshToken())
  const {refreshTokenAction, makeRefreshTokenRequest} = useContext(RefreshTokenContext)
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
      {signInStoreAction.actionType === ACTION_TYPES.SUCCESS && (
        saveStoreId(signInStoreAction.response)
      )}
      
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

      <br />
      <button onClick={goToStorePage}>
        get to storepage
      </button>
    </>
  )
}

export default Temp