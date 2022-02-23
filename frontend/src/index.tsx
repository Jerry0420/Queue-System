import React from 'react'
import ReactDom from 'react-dom'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Store from './components/Store'
import './tailwind.scss'

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Store />} />
      </Routes>
    </BrowserRouter>
  )
}

ReactDom.render(<App />, document.getElementById('root'))
