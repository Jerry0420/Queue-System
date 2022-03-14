import React, {useState} from 'react'
import ReactDom from 'react-dom'
import { HashRouter, Route, Routes } from 'react-router-dom'
import { NoMatch } from './components/Default'
import { Header } from './components/Header'
import { SignUp } from './components/Signup'
import Temp from './components/Temp'
import { useApiRequest } from "./apis/reducer"
import { refreshToken } from './apis/StoreAPIs'
import './tailwind.scss'
import { RefreshTokenContext } from './components/contexts'
import { Store } from './components/Store'
import { checkAuthFlow, validateResponseSuccess } from './apis/helper'
import { SignIn } from './components/Signin'

function App() {
  
  const [refreshTokenAction, makeRefreshTokenRequest] = useApiRequest(...refreshToken())

  const wrapCheckAuthFlow = (nextStuff: () => void, redirectToMainPage: () => void) => {
    checkAuthFlow(refreshTokenAction.response, makeRefreshTokenRequest, 
      // nextStuff
      () => {
        if (validateResponseSuccess(refreshTokenAction.response) === true) {
          nextStuff()
        }
      }, 
      // redirectToMainPage
      () => {
        redirectToMainPage()
      }
    )
  }
  
  return (
    <RefreshTokenContext.Provider value={{refreshTokenAction, makeRefreshTokenRequest, wrapCheckAuthFlow}}>
      <HashRouter>
        <Routes>
          <Route path="/" element={<Header />}>
            <Route path="" element={<SignUp />} />
            <Route path="stores/:storeId" element={(<Store />)} />
            <Route path="stores/:storeId/queues/:queueId" element={(<></>)} />
            <Route path="signin" element={(<SignIn />)} />

            <Route path="stores/:storeId/sessions/:sessionId" element={(<></>)} />
            <Route path="password/forget" element={(<></>)} />
            <Route path="stores/:sessionId/password/update" element={(<></>)} />

            <Route path="temp" element={<Temp />} />

            <Route path="*" element={<NoMatch />} />
          </Route>
        </Routes>
      </HashRouter>
    </RefreshTokenContext.Provider>
  )
}

ReactDom.render(<App />, document.getElementById('root'))
