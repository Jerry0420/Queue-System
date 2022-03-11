import React, {useEffect} from "react"
import { useNavigate } from "react-router-dom"
import { checkExistenceOfRefreshableCookie } from "../apis/helper"
import classNames from 'classnames'
import '../styles/style.scss'

const SignUp = () => {
  let navigate = useNavigate()

  const inputEmailClassNames = classNames({
    'alertInputField': true
  })

  const handleInputEmail = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
  }

  const handleInputPassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
  }

  const handleInputName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
  }

  Intl.DateTimeFormat().resolvedOptions().timeZone
  window.btoa("im password") // base64 password value

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
      <div>sign up</div>

      <input
        type="text"
        onChange={handleInputEmail}
        placeholder="email"
        className={inputEmailClassNames}
      />
      <input
        type="text"
        onChange={handleInputPassword}
        placeholder="password"
      />
      <input
        type="text"
        onChange={handleInputName}
        placeholder="name"
      />
    </div>
  )
}

export {
  SignUp
}