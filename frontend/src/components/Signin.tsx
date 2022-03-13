import React, {useEffect, useState} from "react"
import { useNavigate } from "react-router-dom"
import classNames from 'classnames'
import '../styles/style.scss'

const SignIn = () => {

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
            
        </div>
    )
}

export {
  SignIn
}