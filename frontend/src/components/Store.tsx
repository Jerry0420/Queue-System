import React, {useEffect, useContext, useState} from "react"
import { useParams, Link, useNavigate } from "react-router-dom"
import { RefreshTokenContext } from "./contexts"
import { createSessionWithSSE } from "../apis/SessionAPIs"
import { validateResponseSuccess } from "../apis/helper"
import { ACTION_TYPES, JSONResponse, useApiRequest } from "../apis/reducer"
import { toDataURL } from "qrcode"
import { getStoreInfoWithSSE, updateStoreDescription } from "../apis/StoreAPIs"
import { getNormalTokenFromRefreshTokenAction, getSessionTokenFromRefreshTokenAction } from "../apis/validator"
import { TextareaAutosize, OutlinedInput, FormControl, InputLabel, SelectChangeEvent, Select, MenuItem, DialogActions, Button, Dialog, DialogTitle, DialogContent, Stack, CardContent, CardMedia, Container, Card, List, ListItem, ListItemText, ListItemIcon, Divider, AppBar, Box, Grid, Paper, Avatar, Typography, Drawer, Toolbar, IconButton } from "@mui/material"
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'
import MenuIcon from '@mui/icons-material/Menu'
import { Customer, Queue, Store } from "../apis/models"
import CloseIcon from '@mui/icons-material/Close'
import RefreshIcon from '@mui/icons-material/Refresh'
import HomeIcon from '@mui/icons-material/Home'
import HailIcon from '@mui/icons-material/Hail'
import EscalatorWarningIcon from '@mui/icons-material/EscalatorWarning'
import ExitToAppIcon from '@mui/icons-material/ExitToApp';
import { updateCustomer } from "../apis/CustomerAPIs"

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

  const {refreshTokenAction, makeRefreshTokenRequest, wrapCheckAuthFlow} = useContext(RefreshTokenContext)

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

  const [openUpdateStoreDescriptionDialog, setOpenUpdateStoreDescriptionDialog] = React.useState(false)
  const [storeNewDescription, setStoreNewDescription] = React.useState('')
  const handleClickUpdateStoreDescription = () => {
    setOpenUpdateStoreDescriptionDialog(true)
    setStoreNewDescription(storeInfo.description) // default description
  }

  const handleChangeStoreNewDescription = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = event.target
    setStoreNewDescription(value)
  }

  const handleCloseUpdateDescriptionDialog = () => {
    setOpenUpdateStoreDescriptionDialog(false)
    setStoreNewDescription('')
  }

  const [updateStoreDescriptionAction, makeUpdateStoreDescriptionRequest] = useApiRequest(
    ...updateStoreDescription(
      parseInt(storeId), 
      getNormalTokenFromRefreshTokenAction(refreshTokenAction.response), 
      storeNewDescription
      )
  )

  const doMakeUpdateStoreDescriptionRequest = () => {
    wrapCheckAuthFlow(
      () => {
        makeUpdateStoreDescriptionRequest()
          .then((response) => {
            handleCloseUpdateDescriptionDialog()
          })
      },
      () => {
         // TODO: show error message
         navigate("/")
      }
    )
  }

  const handleUpdateStoreNewDescription = () => {
    doMakeUpdateStoreDescriptionRequest()
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
  
  const [selectedQueueId, setSelectedQueueId] = useState<number | null>(null) 
  const [mainContent, setMainContent] = useState<JSX.Element>((<></>)) 
  const countWaitingOrProcessingCustomers = (customers: Customer[]): Customer[] => {
    return customers.filter((customer: Customer) => customer.status == 'waiting' || customer.status == 'processing')
  }

  const [openUpdateCustomerStatusDialog, setOpenUpdateCustomerStatusDialog] = React.useState(false)
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null) 
  const [customerNewStatus, setCustomerNewStatus] = React.useState('')
  const handleChangeCustomerNewStatus = (event: SelectChangeEvent<typeof customerNewStatus>) => {
    setCustomerNewStatus(event.target.value)
  }

  const handleClickCustomerStatus = (customer: Customer) => {
    setOpenUpdateCustomerStatusDialog(true)
    setSelectedCustomer(customer)
    setCustomerNewStatus(customer.status) //default status
  }

  const handleCloseCustomerStatusDialog = (event: React.SyntheticEvent<unknown>, reason?: string) => {
    setOpenUpdateCustomerStatusDialog(false)
    setCustomerNewStatus('')
    setSelectedCustomer(null)
  }

  const [updateCustomerAction, makeUpdateCustomerRequest] = useApiRequest(
    ...updateCustomer(
        selectedCustomer === null ? -1 : selectedCustomer.id,
        getNormalTokenFromRefreshTokenAction(refreshTokenAction.response), 
        parseInt(storeId),
        selectedQueueId === null ? -1 : selectedQueueId,
        selectedCustomer === null ? '' : selectedCustomer.status,
        customerNewStatus
      )
  )
  const doMakeUpdateCustomerRequest = () => {
    wrapCheckAuthFlow(
      () => {
        makeUpdateCustomerRequest()
          .then((response) => {
            setOpenUpdateCustomerStatusDialog(false)
            setCustomerNewStatus('')
          })
      },
      () => {
         // TODO: show error message
         navigate("/")
      }
    )
  }

  const handleUpdateCustomerNewStatus = () => {
    if (!customerNewStatus) {
      return
    }
    if ((selectedCustomer as Customer).status === customerNewStatus) {
      return
    }
    doMakeUpdateCustomerRequest()
  }

  useEffect(() => {
    if (selectedQueueId === null) {
      setMainContent((
        <>
          <Container maxWidth="md">
            <Container fixed>
              <Typography gutterBottom variant="h5" component="h2" align="center">
                Please scan the QRCode to join the queue.
              </Typography>
            </Container>
            <Grid container rowSpacing={2} justifyContent="center" alignItems="center">
              <Grid item key={"all"} xs={10} sm={10} md={6}>
                <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                  <CardMedia
                    component="img"
                    sx={{
                      display: 'block',
                      marginLeft: 'auto',
                      marginRight: 'auto',
                      width: '80%'
                    }}
                    src={qrcodeImageURL}
                    alt="qrcode image"
                  />
                  <CardContent sx={{ flexGrow: 1 }}>
                    <Typography 
                      gutterBottom 
                      variant="h5" 
                      component="h2"
                      style={{whiteSpace: 'pre-line'}}
                    >
                      {storeInfo.description}
                    </Typography>
                    <a href={sessionScannedURL} target="_blank">{sessionScannedURL}</a>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item key={"queues"} xs={12} sm={12} md={12}>
                <TableContainer component={Paper}>
                  <Table sx={{ minWidth: '50vw' }} aria-label="simple table">
                    <TableHead>
                      <TableRow>
                        <TableCell>Queue Name</TableCell>
                        <TableCell align="right">Await</TableCell>
                        <TableCell align="right">Next</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {queuesInfo.map((queue: Queue) => (
                        <TableRow
                          key={queue.id}
                          sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                        >
                          <TableCell component="th" scope="row">
                            {queue.name}
                          </TableCell>
                          <TableCell align="right">{countWaitingOrProcessingCustomers(queue.customers).length}</TableCell>
                          {countWaitingOrProcessingCustomers(queue.customers).length === 0 && (
                            <TableCell align="right"> - </TableCell>  
                          )}
                          {countWaitingOrProcessingCustomers(queue.customers).length !== 0 && (
                            <TableCell align="right">{countWaitingOrProcessingCustomers(queue.customers)[0].name}</TableCell>
                          )}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </Grid>
            </Grid>
          </Container>
        </>
      ))
    } else {
      let processedCustomers: Customer[]
      const _selectedQueue = queuesInfo.filter((queue: Queue) => queue.id === selectedQueueId)
      const selectedQueue = _selectedQueue[0]
      if (selectedQueue.customers) {
        processedCustomers = countWaitingOrProcessingCustomers(selectedQueue.customers)
      } else {
        processedCustomers = []
      }
      setMainContent((
        <>
          <Box sx={{ width: '100%' }}>
            <Stack 
              spacing={2}
              justifyContent="center"
              alignItems="center"
            >
              <Typography variant="h2" component="h2">{selectedQueue.name}</Typography>
              <TableContainer component={Paper}>
                <Table sx={{ minWidth: '40vw' }} aria-label="simple table">
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell align="right">Phone</TableCell>
                      <TableCell align="right">Status</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {processedCustomers.map((customer: Customer, index) => (
                      <TableRow
                        key={customer.id}
                        sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                      >
                        <TableCell component="th" scope="row">
                          [{index}] {customer.name}
                        </TableCell>

                        <TableCell align="right">
                          {customer.phone}
                        </TableCell>

                        {customer.status === 'waiting' && (
                          <TableCell align="right">
                            <Button onClick={() => handleClickCustomerStatus(customer)}>waiting</Button>
                          </TableCell>
                        )}
                        {customer.status === 'processing' && (
                          <TableCell align="right">
                            <Button sx={{color: 'red'}} onClick={() => handleClickCustomerStatus(customer)}>{customer.status}</Button>
                          </TableCell>
                        )}

                        <Dialog disableEscapeKeyDown open={openUpdateCustomerStatusDialog} onClose={handleCloseCustomerStatusDialog}>
                          <DialogTitle>Update Customer Status</DialogTitle>
                          <DialogContent>
                            <Box component="form" sx={{ display: 'flex', flexWrap: 'wrap' }}>
                              <FormControl sx={{ m: 1, minWidth: 120 }}>
                                <InputLabel id="dialog-select-label">Status</InputLabel>
                                <Select
                                  labelId="dialog-select-label"
                                  id="dialog-select"
                                  value={customerNewStatus}
                                  onChange={handleChangeCustomerNewStatus}
                                  input={<OutlinedInput label="Status" />}
                                >
                                  <MenuItem value="">------</MenuItem>
                                  <MenuItem value={'waiting'}>Waiting</MenuItem>
                                  <MenuItem value={'processing'}>Processing</MenuItem>
                                  <MenuItem value={'done'}>Done</MenuItem>
                                  <MenuItem value={'delete'}>Delete</MenuItem>
                                </Select>
                              </FormControl>
                            </Box>
                          </DialogContent>
                          <DialogActions>
                            <Button onClick={handleCloseCustomerStatusDialog}>Cancel</Button>
                            <Button onClick={handleUpdateCustomerNewStatus}>Ok</Button>
                          </DialogActions>
                        </Dialog>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </Stack>
          </Box>
        </>
      ))
    }
  }, [selectedQueueId, setMainContent, storeInfo, queuesInfo, qrcodeImageURL, openUpdateCustomerStatusDialog, customerNewStatus, setOpenUpdateCustomerStatusDialog, setCustomerNewStatus])

  const drawer = (
    <div>
      <Toolbar />
      <Divider />
          <ListItem button key={"All"} onClick={() => {setSelectedQueueId(null)}}>
            <ListItemIcon>
              <HomeIcon />
            </ListItemIcon>
            <ListItemText primary={"All"} />
          </ListItem>
      <Divider />
      <List>
        {queuesInfo.map((queue: Queue, index) => (
          <ListItem button key={queue.id} onClick={() => {setSelectedQueueId(queue.id)}}>
            <ListItemIcon>
              {index % 2 === 0 ? <HailIcon /> : <EscalatorWarningIcon />}
            </ListItemIcon>
            <ListItemText primary={queue.name} />
          </ListItem>
        ))}
      </List>
      <Divider />
      <List>
        <ListItem button key={"Update Description"} onClick={handleClickUpdateStoreDescription}>
          <ListItemIcon>
            <RefreshIcon />
          </ListItemIcon>
          <ListItemText primary={"Update Description"} />
        </ListItem>
        <Dialog disableEscapeKeyDown open={openUpdateStoreDescriptionDialog} onClose={handleCloseUpdateDescriptionDialog}>
          <DialogTitle>Update Description</DialogTitle>
          <DialogContent>
            <Box component="form" sx={{ display: 'flex', flexWrap: 'wrap' }}>
              <FormControl sx={{ m: 1, minWidth: 120 }}>
                <TextareaAutosize
                  aria-label="empty textarea"
                  placeholder="Store Description"
                  value={storeNewDescription}
                  style={{ width: 200 }}
                  onChange={handleChangeStoreNewDescription}
                />
              </FormControl>
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleCloseUpdateDescriptionDialog}>Cancel</Button>
            <Button onClick={handleUpdateStoreNewDescription}>Ok</Button>
          </DialogActions>
        </Dialog>
        
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
            keepMounted: true,
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
  )
}


export {
  StoreInfo
}