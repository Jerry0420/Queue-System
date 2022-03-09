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

function App() {
  
  const [refreshTokenAction, makeRefreshTokenRequest] = useApiRequest(...refreshToken())
  
  return (
    <RefreshTokenContext.Provider value={{refreshTokenAction, makeRefreshTokenRequest}}>
      <HashRouter>
        <Routes>
          <Route path="/" element={<Header />}>
            <Route path="" element={<SignUp />} />
            <Route path="stores/:storeId" element={(<></>)} />
            <Route path="stores/:storeId/queues/:queueId" element={(<></>)} />
            <Route path="signin" element={(<></>)} />

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
