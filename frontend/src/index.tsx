import React, {useState} from 'react'
import ReactDom from 'react-dom'
import { HashRouter, Route, Routes } from 'react-router-dom'
import { NoMatch } from './components/Default'
import Store from './components/Store'
import './tailwind.scss'

function App() {
  return (
    <HashRouter>
      <Routes>
        <Route path="/" element={<Store />} />
        <Route path="/stores/:storeId" element={(<></>)} />
        <Route path="/stores/:storeId/queues/:queueId" element={(<></>)} />
        <Route path="/signin" element={(<></>)} />

        <Route path="/stores/:storeId/sessions/:sessionId" element={(<></>)} />
        <Route path="/password/forget" element={(<></>)} />
        <Route path="/stores/:sessionId/password/update" element={(<></>)} />

        <Route path="*" element={<NoMatch />} />
      </Routes>
    </HashRouter>
  )
}

ReactDom.render(<App />, document.getElementById('root'))
