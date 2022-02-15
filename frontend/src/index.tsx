import React from 'react'
import ReactDom from 'react-dom'
// import styles from './index.scss'
import './index.scss'

function Main() {
  return (
    <div className="main text-3xl font-bold underline">
      hello world
    </div>
    )
}

ReactDom.render(<Main />, document.getElementById('root'))
