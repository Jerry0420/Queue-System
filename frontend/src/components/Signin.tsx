import React, {useEffect, useState} from "react"
import { 
  Link as RouterLink,
  useNavigate 
} from "react-router-dom"
import classNames from 'classnames'
import '../styles/style.scss'
import { checkExistenceOfRefreshableCookie } from "../apis/helper"
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import { signInStore, forgetPassword } from "../apis/StoreAPIs"
import { Button, Box, Grid, Paper, Avatar, Typography, TextField, Link, DialogActions, Dialog, DialogTitle, DialogContent } from "@mui/material"

const SignIn = () => {
  let navigate = useNavigate()

  const [email, setEmail] = useState("")
  const [emailAlertFlag, setEmailAlertFlag] = useState(false)
  const handleInputEmail = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setEmail(value)
  }

  const [password, setPassword] = useState("")
  const [passwordAlertFlag, setPasswordAlertFlag] = useState(false)
  const handleInputPassword = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setPassword(window.btoa(value))
  }

  const [signInStoreAction, makeSignInStoreRequest] = useApiRequest(...signInStore(email, password))

  const doMakeSignInStoreRequest = () => {
    const validateEmail = (inputEmail: string) => {
      return inputEmail.match(/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)
    }
    if (validateEmail(email)) {
      setEmailAlertFlag(false)
    } else {
      setEmailAlertFlag(true)
      return
    }

    const rawPassword = window.atob(password)
    if ((8 <= rawPassword.length) && (rawPassword.length <= 15)) {
      setPasswordAlertFlag(false)
      setPassword(window.btoa(password)) // base64 password value
    } else {
      setPasswordAlertFlag(true)
      return
    }

    if (email && rawPassword) {
      makeSignInStoreRequest()
    }
  }

  useEffect(() => {
    // TODO: handle running, success, error states here.
    if (signInStoreAction.actionType === ACTION_TYPES.SUCCESS) {
      const _jsonResponse = (signInStoreAction.response as JSONResponse)
      const storeId: number = (_jsonResponse["id"] as number)
      localStorage.setItem("storeId", storeId.toString())
      navigate(`/stores/${storeId}`)
    }
  }, [signInStoreAction.actionType])

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

  const [openForgetPasswordDialog, setOpenForgetPasswordDialog] = React.useState(false)
  const [forgetPasswordEmail, setForgetPasswordEmail] = useState("")
  const [forgetPasswordEmailAlertFlag, setForgetPasswordEmailAlertFlag] = useState(false)
  const handleInputForgetPasswordEmail = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setForgetPasswordEmail(value)
  }
  const [forgetPasswordAction, makeForgetPasswordRequest] = useApiRequest(...forgetPassword(forgetPasswordEmail))
  const handleForgetPassword = () => {
    const validateEmail = (inputEmail: string) => {
      return inputEmail.match(/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)
    }
    if (validateEmail(forgetPasswordEmail)) {
      setForgetPasswordEmailAlertFlag(false)
    } else {
      setForgetPasswordEmailAlertFlag(true)
      return
    }
    makeForgetPasswordRequest().then((response) => {
      setOpenForgetPasswordDialog(false)
    })
  }

  return (
    <Box sx={{flexGrow: 1}}>
      <Grid container direction="row-reverse" component="main" sx={{ height: '100vh' }}>
        <Grid
          item
          xs={false}
          sm={false}
          md={7}
          sx={{
            backgroundImage: 'url(https://images.unsplash.com/photo-1506774518161-b710d10e2733?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80)',
            backgroundRepeat: 'no-repeat',
            backgroundColor: (t) =>
              t.palette.mode === 'light' ? t.palette.grey[50] : t.palette.grey[900],
            backgroundSize: 'cover',
            backgroundPosition: 'center',
          }}
        />
        <Grid item xs={12} sm={12} md={5} component={Paper} elevation={6} square>
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
              Signin Store
            </Typography>
            <Box sx={{ mt: 1 }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                onChange={handleInputEmail}
                error={emailAlertFlag}
              />
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
              />     
              <Button
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={doMakeSignInStoreRequest}
              >
                Sign In
              </Button>
              <Grid container>
                <Grid item xs>
                  <Link variant="body2" sx={{"&:hover": {cursor: "pointer"}}} onClick={() => {setOpenForgetPasswordDialog(true)}}>
                    Forgot password?
                  </Link>
                  <Dialog disableEscapeKeyDown open={openForgetPasswordDialog} onClose={() => {setOpenForgetPasswordDialog(false)}}>
                    <DialogTitle>Forget Password</DialogTitle>
                    <DialogContent>
                      <TextField
                        autoFocus
                        margin="dense"
                        id="email"
                        label="Email Address"
                        type="email"
                        fullWidth
                        variant="standard"
                        autoComplete="email"
                        onChange={handleInputForgetPasswordEmail}
                        error={forgetPasswordEmailAlertFlag}
                      />
                      We'll send a link to the email for resetting password.
                    </DialogContent>
                    <DialogActions>
                      <Button onClick={() => {setOpenForgetPasswordDialog(false)}}>Cancel</Button>
                      <Button onClick={handleForgetPassword}>Ok</Button>
                    </DialogActions>
                  </Dialog>
                </Grid>
                <Grid item>
                  <Link component={RouterLink} variant="body2" to="/">
                    {"Don't have an account? Sign Up"}
                  </Link>
                </Grid>
              </Grid>
              {/* <Copyright sx={{ mt: 5 }} /> */}
            </Box>
          </Box>
        </Grid>
      </Grid>
    </Box>
  )
}

export {
  SignIn
}