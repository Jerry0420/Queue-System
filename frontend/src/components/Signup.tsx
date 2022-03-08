import React, {useContext, useEffect} from "react"
import { RefreshTokenContext } from "./contexts"
import { JSONResponse } from "../apis/reducer"

const SignUp = () => {
  const {refreshTokenAction, makeRefreshTokenRequest} = useContext(RefreshTokenContext)
  useEffect(() => {
    makeRefreshTokenRequest().then(jsonRespons => {
      console.log(refreshTokenAction)
    })
  }, [refreshTokenAction])

  return (
    <div>
      {/* <>{console.log(refreshTokenAction)}</>  */}
      <button onClick={makeRefreshTokenRequest}>
        refreshToken
      </button>
      {/* <>{console.log(refreshTokenAction)}</>  */}
    </div>
  )
}

export {
  SignUp
}