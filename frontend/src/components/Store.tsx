import React, {useEffect, useContext, useState} from "react"
import { useParams, Link, useNavigate } from "react-router-dom"
import { RefreshTokenContext } from "./contexts"
import { createSessionWithSSE } from "../apis/SessionAPIs"
import { validateResponseSuccess } from "../apis/helper"
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import { toDataURL } from "qrcode"
import { getStoreInfoWithSSE, updateStoreDescription } from "../apis/StoreAPIs"
import { getNormalTokenFromRefreshTokenAction, getSessionTokenFromRefreshTokenAction } from "../apis/validator"
import { List, ListItem, ListItemText, ListItemIcon, Divider, AppBar, Box, Grid, Paper, Avatar, Typography, Drawer, Toolbar, IconButton } from "@mui/material"
import MenuIcon from '@mui/icons-material/Menu'
import { Customer, Queue, Store } from "../apis/models"
import CloseIcon from '@mui/icons-material/Close'
import RefreshIcon from '@mui/icons-material/Refresh'
import HomeIcon from '@mui/icons-material/Home'
import HailIcon from '@mui/icons-material/Hail'
import EscalatorWarningIcon from '@mui/icons-material/EscalatorWarning'
import ExitToAppIcon from '@mui/icons-material/ExitToApp';

const StoreInfo = () => {
  let { storeId }: {storeId: string} = useParams()
  
  const drawerWidth = 240
  
  const [mobileOpen, setMobileOpen] = useState(false)
  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen)
  }

  let navigate = useNavigate()
  const [sessionScannedURL, setSessionScannedURL] = useState("")
  const [qrcodeImageURL, setQrcodeImageURL] = useState("")
  const [storeDescription, setStoreDescription] = useState("")

  const {refreshTokenAction, makeRefreshTokenRequest, wrapCheckAuthFlow} = useContext(RefreshTokenContext)
  const [updateStoreDescriptionAction, makeUpdateStoreDescriptionRequest] = useApiRequest(
    ...updateStoreDescription(
      parseInt(storeId), 
      getNormalTokenFromRefreshTokenAction(refreshTokenAction.response), 
      storeDescription
      )
  )

  const handleInputStoreDescription = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value: value }: { value: string } = e.target
    setStoreDescription(value)
  }

  const [storeInfo, setStoreInfo] = useState<Store>({})
  const [queuesInfo, setQueuesInfo] = useState<Queue[]>([])
  useEffect(() => {
    let getStoreInfoSSE: EventSource
    getStoreInfoSSE = getStoreInfoWithSSE(parseInt(storeId))

    getStoreInfoSSE.onmessage = (event) => {
      // TODO: render to ui
      setStoreInfo(JSON.parse(event.data))
      setQueuesInfo(JSON.parse(event.data)['queues'])
      console.log(JSON.parse(event.data))
    }
    
    getStoreInfoSSE.onerror = (event) => {
      getStoreInfoSSE.close()
    }
    return () => {
      if (getStoreInfoSSE != null) {
        getStoreInfoSSE.close()
      }
    }
  }, [getStoreInfoWithSSE, setStoreInfo])

  useEffect(() => {
    let createSessionSSE: EventSource
    wrapCheckAuthFlow(
      () => {
        const sessionToken: string = getSessionTokenFromRefreshTokenAction(refreshTokenAction.response)
        createSessionSSE = createSessionWithSSE(sessionToken)

        createSessionSSE.onmessage = (event) => {
          setSessionScannedURL(JSON.parse(event.data)["scanned_url"])
        }
        
        createSessionSSE.onerror = (event) => {
          createSessionSSE.close()
        }
      },
      () => {
         // TODO: show error message
         navigate("/")
      }
    )
    return () => {
      if (createSessionSSE != null) {
        createSessionSSE.close()
      }
    }
  }, [createSessionWithSSE, refreshTokenAction.response, refreshTokenAction.exception])

  const doMakeUpdateStoreDescriptionRequest = () => {
    wrapCheckAuthFlow(
      () => {
        makeUpdateStoreDescriptionRequest()
      },
      () => {
         // TODO: show error message
         navigate("/")
      }
    )
  }

  useEffect(() => {
    toDataURL(sessionScannedURL, (error, url) => {
      if (url != null) {
        setQrcodeImageURL(url)
      }
    })
  }, [sessionScannedURL])

  useEffect(() => {
    // TODO: handle running, success, error states here.
  }, [updateStoreDescriptionAction.actionType])
  
  const [selectQueueId, setSelectQueueId] = useState<number | null>(null) 
  const [mainContent, setMainContent] = useState<JSX.Element>((<></>)) 
  useEffect(() => {
    if (selectQueueId === null) {
      setMainContent((
        <>
          <div>
            hello world
          </div>
        </>
      ))
    } else {
      setMainContent((
        <>
          <Typography paragraph>
            Lorem ipsum ac.
          </Typography>
          <a href={sessionScannedURL} target="_blank">{sessionScannedURL}</a>
        </>
      ))
    }
  }, [selectQueueId, setMainContent, storeInfo, queuesInfo])

  const drawer = (
    <div>
      <Toolbar />
      <Divider />
          <ListItem button key={"All"} onClick={() => {setSelectQueueId(null)}}>
            <ListItemIcon>
              <HomeIcon />
            </ListItemIcon>
            <ListItemText primary={"All"} />
          </ListItem>
      <Divider />
      <List>
        {queuesInfo.map((queue: Queue, index) => (
          <ListItem button key={queue.name} onClick={() => {setSelectQueueId(queue.id)}}>
            <ListItemIcon>
              {index % 2 === 0 ? <HailIcon /> : <EscalatorWarningIcon />}
            </ListItemIcon>
            <ListItemText primary={queue.name} />
          </ListItem>
        ))}
      </List>
      <Divider />
      <List>
        <ListItem button key={"Update Description"}>
          <ListItemIcon>
            <RefreshIcon />
          </ListItemIcon>
          <ListItemText primary={"Update Description"} />
        </ListItem>
        
        <ListItem button key={"Close Store"}>
          <ListItemIcon>
            <CloseIcon />
          </ListItemIcon>
          <ListItemText primary={"Close Store"} />
        </ListItem>

        <ListItem button key={"Sign Out"}>
          <ListItemIcon>
            <ExitToAppIcon />
          </ListItemIcon>
          <ListItemText primary={"Sign Out"} />
        </ListItem>

      </List>
    </div>
  )

  return (
    <Box sx={{ display: 'flex' }}>
      <AppBar
        position="fixed"
        sx={{
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { sm: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div">
            {storeInfo.name}
          </Typography>
        </Toolbar>
      </AppBar>
      <Box
        component="nav"
        sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
        // aria-label="mailbox folders"
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Better open performance on mobile.
          }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', sm: 'block' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{ flexGrow: 1, p: 3, width: { sm: `calc(100% - ${drawerWidth}px)` } }}
      >
        <Toolbar />
        {mainContent}
      </Box>
    </Box>
    // <div>
    //     <Link to="/temp">to temp</Link>
    //     {/* <img src={qrcodeImageURL} alt="qrcode image"></img> */}

    //     <br />
    //     <input
    //       type="text"
    //       onChange={handleInputStoreDescription}
    //       placeholder="input store description"
    //     />
    //     <button onClick={doMakeUpdateStoreDescriptionRequest}>
    //       update store description
    //     </button>

    //     <br />
    //     <a href={sessionScannedURL} target="_blank">{sessionScannedURL}</a>
    // </div>
  )
}


export {
  StoreInfo
}