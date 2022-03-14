import React, {useEffect, useState} from "react"
import { useNavigate, Link } from "react-router-dom"
import { checkExistenceOfRefreshableCookie } from "../apis/helper"
import classNames from 'classnames'
import '../styles/style.scss'
import { ACTION_TYPES, useApiRequest } from "../apis/reducer"
import { openStore } from "../apis/StoreAPIs"

const SignUp = () => {
  let navigate = useNavigate()
  
  const timezone: string = Intl.DateTimeFormat().resolvedOptions().timeZone
  
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

  const [name, setName] = useState("")
  const [nameAlertFlag, setNameAlertFlag] = useState(false)
  const handleInputName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    if (value) {
      setNameAlertFlag(false)
      setName(value)
    } else {
      setNameAlertFlag(true)
    }
  }

  const [queueName, setQueueName] = useState("")
  const [queueNameAlertFlag, setQueueNameAlertFlag] = useState(false)
  const handleInputQueueName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    if (value) {
      setQueueNameAlertFlag(false)
      setQueueName(value)
    } else {
      setQueueNameAlertFlag(true)
    }
  }

  const [queueNames, setQueueNames] = useState<string[]>([])
  useEffect(() => {
    if (queueName) {
      const _queueNames = [...queueNames]
      _queueNames.push(queueName)
      setQueueNames(_queueNames)
    }
  }, [queueName])  

  const clearQueueNames = () => {
    setQueueNames([])
  }

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

  const [openStoreAction, makeOpenStoreRequest] = useApiRequest(
    ...openStore(email, password, name, timezone, queueNames)
    )

  const [openStoreFlag, setOpenStoreFlag] = useState(false)
  useEffect(() => {
    if (email && password && name && timezone && queueNames.length > 0) {
      setOpenStoreFlag(false)
    } else {
      setOpenStoreFlag(true)
    }
  }, [email, password, name, timezone, queueNames])

  useEffect(() => {
    // TODO: handle running, success, error states here.
    if (openStoreAction.actionType === ACTION_TYPES.SUCCESS) {
      // TODO: show alert to signin page
      navigate("/signin")
    }
  }, [openStoreAction.actionType])

  return (
    <div>
      <div>sign up</div>

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
      <input
        type="text"
        onBlur={handleInputName}
        className={classNames({'alertInputField': nameAlertFlag})}
        placeholder="name"
      />

      <input
        type="text"
        onBlur={handleInputQueueName}
        className={classNames({'alertInputField': queueNameAlertFlag})}
        placeholder="queue name"
      />
      <button onClick={clearQueueNames}>
        clear queue names
      </button>

      {queueNames.map((queueName: string) => (
        <div id={queueName} key={queueName}>{queueName}</div>
      ))}

      <br />
      <button 
        onClick={makeOpenStoreRequest}
        disabled={openStoreFlag}
      >
        open store
      </button>

      <br />
      
      <Link to="/signin">signin</Link>
    </div>
  )
}

export {
  SignUp
}