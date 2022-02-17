import React from "react"
import '../styles/Main.scss'
import classNames from "classnames"

const Main = () => {
  const divClass = classNames("main", "font-bold", "underline")
  return (
    <div className={divClass}>
      hello world
    </div>
  )
}

export default Main