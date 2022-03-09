import React, {useEffect} from "react"
import { ACTION_TYPES } from "../apis/reducer"
import { useNavigate } from "react-router-dom"
import { checkExistenceOfRefreshableCookie } from "../apis/helper"


const SignUp = () => {
  let navigate = useNavigate()

  useEffect(() => {
    if (checkExistenceOfRefreshableCookie() === true) {
      const storeId = localStorage.getItem("storeId")
      if (storeId != null) {
        // navigate("/temp")
        navigate(`/stores/${storeId}`)
      } else {
        document.cookie = "refreshable=true ; expires = Thu, 01 Jan 1970 00:00:00 GMT"
      }
    }
  }, [])

  return (
    <div>
      <div>sign up</div>
    </div>
  )
}

export {
  SignUp
}