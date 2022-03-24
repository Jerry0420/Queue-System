import React, {useEffect, useState} from "react"
import {
  Link as RouterLink,
  useNavigate,
} from 'react-router-dom';
import { checkExistenceOfRefreshableCookie } from "../apis/helper"
import { ACTION_TYPES, useApiRequest } from "../apis/reducer"
import { openStore } from "../apis/StoreAPIs"
import { Chip, Button, Box, Grid, Paper, Avatar, Typography, TextField, Link } from "@mui/material"
import AddBoxIcon from '@mui/icons-material/AddBox';

const SignUp = () => {
  let navigate = useNavigate()
  
  const timezone: string = Intl.DateTimeFormat().resolvedOptions().timeZone
  
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
    setPassword(value)
  }

  const [name, setName] = useState("")
  const [nameAlertFlag, setNameAlertFlag] = useState(false)
  const handleInputName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setName(value)
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

  const [addQueueNameAlertFlag, setAddQueueNameAlertFlag] = useState(false)
  useEffect(() => {
    if (queueName) {
      setAddQueueNameAlertFlag(false)
    } else {
      setAddQueueNameAlertFlag(true)
    }
  }, [queueName])

  const [queueNames, setQueueNames] = useState<string[]>([])
  
  const addQueueNameToQueueNames = () => {
    const _queueNames = [...queueNames]
    _queueNames.push(queueName)
    setQueueNames(_queueNames)
    setQueueName("")
  }

  const handleDeleteQueueName = (deletedQueueName: string) => {
      var _queueNames = queueNames.filter((value, index, error): boolean => {
        return value != deletedQueueName
      })
      setQueueNames(_queueNames)
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

  const doMakeOpenStoreRequest = () => {
    const validateEmail = (inputEmail: string) => {
      return inputEmail.match(/^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)
    }
    if (validateEmail(email)) {
      setEmailAlertFlag(false)
    } else {
      setEmailAlertFlag(true)
    }

    if ((8 <= password.length) && (password.length <= 15)) {
      setPasswordAlertFlag(false)
      setPassword(window.btoa(password)) // base64 password value
    } else {
      setPasswordAlertFlag(true)
    }

    if (name) {
      setNameAlertFlag(false)
    } else {
      setNameAlertFlag(true)
    }

    if (queueNames.length > 0) {
      setQueueNameAlertFlag(false)
    } else {
      setQueueNameAlertFlag(true)
    }

    if (email && password && name && timezone && queueNames.length > 0) {
      makeOpenStoreRequest()
    }
  }

  useEffect(() => {
    // TODO: handle running, success, error states here.
    if (openStoreAction.actionType === ACTION_TYPES.SUCCESS) {
      // TODO: show alert to signin page
      navigate("/signin")
    }
  }, [openStoreAction.actionType])

  return (
    <Box sx={{flexGrow: 1}}>
      <Grid container component="main" sx={{ height: '100vh' }}>
        <Grid
          item
          xs={false}
          sm={false}
          md={7}
          sx={{
            backgroundImage: 'url(https://source.unsplash.com/random)',
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
              Open Store
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
                helperText="Use 8 to 15 characters with a mix of letters, numbers & symbols"
              />
              <TextField
                margin="normal"
                required
                fullWidth
                name="name"
                label="Store Name"
                type="text"
                id="name"
                autoComplete="name"
                onChange={handleInputName}
                error={nameAlertFlag}
              />
              <Grid 
                container 
                spacing={2}
                alignItems="center"
                justifyContent="flex-start"
              >
                <Grid item xs={8} sm={8}>
                  <TextField
                    fullWidth
                    required
                    margin="normal"
                    name="queueName"
                    label="Queue Name"
                    type="text"
                    id="queueName"
                    onChange={handleInputQueueName}
                    value={queueName}
                    error={queueNameAlertFlag}
                  />
                </Grid>
                <Grid item xs={4} sm={4}>
                  <Button 
                    variant="contained" 
                    startIcon={<AddBoxIcon />}
                    onClick={addQueueNameToQueueNames}
                    disabled={addQueueNameAlertFlag}
                  >
                    Add
                  </Button>
                </Grid>
              </Grid>              

              {queueNames.map((queueName: string) => (
                  <Chip 
                    sx={{ mb: 1, ml: 1, mr: 1 }}
                    label={queueName}
                    key={queueName} 
                    onDelete={() => {handleDeleteQueueName(queueName)}}
                  />
                ))}

              <Button
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
                onClick={doMakeOpenStoreRequest}
              >
                Sign In
              </Button>
              <Grid container>
                <Grid item xs>
                  <Link component={RouterLink} variant="body2" to="/password/forget">
                    Forgot password?
                  </Link>
                </Grid>
                <Grid item>
                  <Link component={RouterLink} variant="body2" to="/signin">
                    {"Already have an account? Sign In"}
                  </Link>
                </Grid>
              </Grid>
              {/* <Copyright sx={{ mt: 5 }} /> */}
            </Box>
          </Box>
        </Grid>
      </Grid>
    </Box>
    // <div>
    //   <div>sign up</div>

    //   <input
    //     type="text"
    //     placeholder="email"
    //     className={classNames({'alertInputField': emailAlertFlag})}
    //     onBlur={handleInputEmail}
    //   />
    //   <input
    //     type="text"
    //     onBlur={handleInputPassword}
    //     className={classNames({'alertInputField': passwordAlertFlag})}
    //     placeholder="password"
    //   />
    //   <input
    //     type="text"
    //     onBlur={handleInputName}
    //     className={classNames({'alertInputField': nameAlertFlag})}
    //     placeholder="name"
    //   />

    //   <input
    //     type="text"
    //     onBlur={handleInputQueueName}
    //     className={classNames({'alertInputField': queueNameAlertFlag})}
    //     placeholder="queue name"
    //   />
    //   <button onClick={clearQueueNames}>
    //     clear queue names
    //   </button>

      // {queueNames.map((queueName: string) => (
      //   <div id={queueName} key={queueName}>{queueName}</div>
      // ))}

    //   <br />
    //   <button 
    //     onClick={makeOpenStoreRequest}
    //     disabled={openStoreFlag}
    //   >
    //     open store
    //   </button>

    //   <br />
      
    //   <Link to="/signin">signin</Link>
    // </div>
  )
}

export {
  SignUp
}