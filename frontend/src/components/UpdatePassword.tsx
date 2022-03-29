import React, {useEffect, useState} from "react"
import { useNavigate, useParams, useLocation } from "react-router-dom"
import '../styles/style.scss'
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import { updatePassword } from "../apis/StoreAPIs"
import { Button, Box, Avatar, Typography, TextField } from "@mui/material"

const UpdatePasswordComponent = () => {
  let navigate = useNavigate()
  let { storeId }: {storeId: string} = useParams()
  let location = useLocation()
  let passwordToken = ""

  if (location.search.includes("password_token")) {
    const splittedQueryStrings = location.search.split("=") 
    passwordToken = splittedQueryStrings[splittedQueryStrings.length-1]
  } else {
    passwordToken = ""
  }

  const [password, setPassword] = useState("")
  const [passwordAlertFlag, setPasswordAlertFlag] = useState(false)
  const handleInputPassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setPassword(window.btoa(value))
  }

  const [confirmPassword, setConfirmPassword] = useState("")
  const [confirmPasswordAlertFlag, setConfirmPasswordAlertFlag] = useState(false)
  const handleInputConfirmPassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setConfirmPassword(window.btoa(value))
  }

  const [updatePasswordAction, makeUpdatePasswordRequest] = useApiRequest(...updatePassword(parseInt(storeId), passwordToken, password))

  const doMakeUpdatePasswordRequest = () => {
    const rawPassword = window.atob(password)
    if ((8 <= rawPassword.length) && (rawPassword.length <= 15)) {
      setPasswordAlertFlag(false)
      setPassword(window.btoa(password)) // base64 password value
    } else {
      setPasswordAlertFlag(true)
      return
    }

    const rawConfirmPassword = window.atob(confirmPassword)
    if (rawConfirmPassword === rawPassword) {
        setConfirmPasswordAlertFlag(false)
    } else {
        setConfirmPasswordAlertFlag(true)
        return
    }

    if (rawConfirmPassword && rawPassword) {
        makeUpdatePasswordRequest()
    }
  }

  useEffect(() => {
    // TODO: handle running, success, error states here.
  }, [updatePasswordAction.actionType])

  return (
    <Box
        sx={{
            my: 8,
            mx: 4,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
        }}
        >
        <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }} />
        <Typography component="h1" variant="h5">
            Reset Password
        </Typography>
        <Box sx={{ mt: 1 }}>
            <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="Password"
                type="password"
                id="password"
                autoComplete="current-password"
                onChange={handleInputPassword}
                error={passwordAlertFlag}
                helperText="Use 8 to 15 characters with a mix of letters, numbers & symbols"
            />
            <TextField
                margin="normal"
                required
                fullWidth
                name="passwordConfirm"
                label="Confirm Password"
                id="passwordConfirm"
                onChange={handleInputConfirmPassword}
                error={confirmPasswordAlertFlag}
                helperText="Type again for confirmation of the password."
            />
            <Button
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={doMakeUpdatePasswordRequest}
                >
                Reset Password
            </Button>
            {/* <Copyright sx={{ mt: 5 }} /> */}
        </Box>
    </Box>
  )
}

export {
  UpdatePasswordComponent
}