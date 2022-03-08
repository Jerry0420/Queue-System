import React, {useContext, useEffect} from "react"
import { RefreshTokenContext } from "./contexts"
import { ACTION_TYPES } from "../apis/reducer"
import { useNavigate } from "react-router-dom"
import { validateResponseSuccess } from "../apis/helper"

const SignUp = () => {
  const {refreshTokenAction, makeRefreshTokenRequest} = useContext(RefreshTokenContext)
  let navigate = useNavigate()

  useEffect(() => {
    if (refreshTokenAction.response == null) {
      makeRefreshTokenRequest()
    }
  }, [])

  return (
    <div>
      {
        (refreshTokenAction.actionType === ACTION_TYPES.SUCCESS) && 
        (validateResponseSuccess(refreshTokenAction.response) === true) &&
        (navigate("/temp"))
      }
      {
        (refreshTokenAction.actionType === ACTION_TYPES.SUCCESS) && 
        (validateResponseSuccess(refreshTokenAction.response) === false) &&
        (
          <div>sign up</div>
        )
      }
    </div>
  )
}

export {
  SignUp
}