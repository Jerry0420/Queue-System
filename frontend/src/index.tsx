import React from 'react'
import ReactDom from 'react-dom'
// import styles from './index.scss'
import './index.scss'

function Main() {
  return (
    <p className="main">
        Hi JSX！
    </p>
    )
}

ReactDom.render(<Main />, document.getElementById('root'))
