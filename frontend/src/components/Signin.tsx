import React, {useEffect, useState} from "react"
import { useNavigate } from "react-router-dom"
import classNames from 'classnames'
import '../styles/style.scss'
import { checkExistenceOfRefreshableCookie } from "../apis/helper"
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import { signInStore } from "../apis/StoreAPIs"

const SignIn = () => {
  let navigate = useNavigate()

  const [email, setEmail] = useState("")
  const [emailAlertFlag, setEmailAlertFlag] = useState(false)
  const handleInputEmail = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    const validateEmail = (inputEmail: string) => {
      return inputEmail.match(/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)
    }
    if (validateEmail(value)) {
      setEmailAlertFlag(false)
      setEmail(value)
    } else {
      setEmailAlertFlag(true)
    }
  }

  const [password, setPassword] = useState("")
  const [passwordAlertFlag, setPasswordAlertFlag] = useState(false)
  const handleInputPassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    if ((8 <= value.length) && (value.length <= 15)) {
      setPasswordAlertFlag(false)
      setPassword(window.btoa(value)) // base64 password value
    } else {
      setPasswordAlertFlag(true)
    }
  }

  const [signInStoreAction, makeSignInStoreRequest] = useApiRequest(...signInStore(email, password))

  useEffect(() => {
    // TODO: handle running, success, error states here.
    if (signInStoreAction.actionType === ACTION_TYPES.SUCCESS) {
      const _jsonResponse = (signInStoreAction.response as JSONResponse)
      const storeId: number = (_jsonResponse["id"] as number)
      localStorage.setItem("storeId", storeId.toString())
      navigate(`/stores/${storeId}`)
    }
  }, [signInStoreAction.actionType])

  const [signInStoreFlag, setSignInStoreFlag] = useState(false)
  useEffect(() => {
    if (email && password) {
      setSignInStoreFlag(false)
    } else {
      setSignInStoreFlag(true)
    }
  }, [email, password])

  useEffect(() => {
    if (checkExistenceOfRefreshableCookie() === true) {
      const storeId = localStorage.getItem("storeId")
      if (storeId != null) {
        navigate(`/stores/${storeId}`)
      } else {
        // remove refreshable cookie for signup again
        document.cookie = "refreshable=true ; expires = Thu, 01 Jan 1970 00:00:00 GMT"
      }
    }
  }, [])

  return (
      <div>
          <div>sign in</div>
          <input
              type="text"
              placeholder="email"
              className={classNames({'alertInputField': emailAlertFlag})}
              onBlur={handleInputEmail}
          />
          <input
              type="text"
              onBlur={handleInputPassword}
              className={classNames({'alertInputField': passwordAlertFlag})}
              placeholder="password"
          />

        <br />
        <button 
          onClick={makeSignInStoreRequest}
          disabled={signInStoreFlag}
        >
          signin store
        </button>
          
      </div>
  )
}

export {
  SignIn
}