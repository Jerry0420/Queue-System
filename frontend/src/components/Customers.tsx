import React, { useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import { createCustomers } from "../apis/CustomerAPIs"
import { CustomerForm } from "../apis/models"
import { ACTION_TYPES, useApiRequest } from "../apis/reducer"
import { scanSession } from "../apis/SessionAPIs"
import { getStoreInfoWithSSE } from "../apis/StoreAPIs"
import classNames from 'classnames'
import { AppBarWDrawer } from "./AppBarWDrawer"
import { Customer, Queue, Store } from "../apis/models"
import { TextField, OutlinedInput, FormControl, InputLabel, SelectChangeEvent, Select, MenuItem, Chip, Button, Stack, CardContent, CardMedia, Container, Card, List, ListItem, ListItemText, ListItemIcon, Divider, FormHelperText, Box, Grid, Paper, Avatar, Typography, Drawer, Toolbar, IconButton } from "@mui/material"
import Table from '@mui/material/Table'
import TableBody from '@mui/material/TableBody'
import TableCell from '@mui/material/TableCell'
import TableContainer from '@mui/material/TableContainer'
import TableHead from '@mui/material/TableHead'
import TableRow from '@mui/material/TableRow'
import AddBoxIcon from '@mui/icons-material/AddBox';

const CreateCustomers = () => {
  let { storeId , sessionId}: {storeId: string, sessionId: string} = useParams()
  
  // ================= store info sse =================
  const [storeInfo, setStoreInfo] = useState<Store>({})
  const [queuesInfo, setQueuesInfo] = useState<Queue[]>([])

  useEffect(() => {
    let getStoreInfoSSE: EventSource
    getStoreInfoSSE = getStoreInfoWithSSE(parseInt(storeId))

    getStoreInfoSSE.onmessage = (event) => {
        // TODO: render customers ui
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
  }, [getStoreInfoWithSSE])

  // ====================== main content ======================
  const [mainContent, setMainContent] = useState<JSX.Element>((<></>))
  const [selectedQueueId, setSelectedQueueId] = useState<number | null>(null)
  // helper function
  const countWaitingOrProcessingCustomers = (customers: Customer[]): Customer[] => {
    return customers.filter((customer: Customer) => customer.status == 'waiting' || customer.status == 'processing')
  }

  const [customerName, setCustomerName] = useState("")
  const [customerNameAlertFlag, setCustomerNameAlertFlag] = useState(false)
  const handleInputCustomerName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setCustomerName(value)
  }

  const [customerPhone, setCustomerPhone] = useState("")
  const [customerPhoneAlertFlag, setCustomerPhoneAlertFlag] = useState(false)
  const handleInputCustomerPhone = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    setCustomerPhone(value)
  }

  // =================
  // const [queueId, setQueueId] = useState("")
  // const [queueIdAlertFlag, setQueueNameAlertFlag] = useState(false)
  // const handleInputQueueName = (e: React.ChangeEvent<HTMLInputElement>) => {
  //   const { value }: { value: string } = e.target
  //   setQueueName(value)
  // }

  // const [addQueueNameAlertFlag, setAddQueueNameAlertFlag] = useState(false)
  // useEffect(() => {
  //   if (queueName) {
  //     setAddQueueNameAlertFlag(false)
  //   } else {
  //     setAddQueueNameAlertFlag(true)
  //   }
  // }, [queueName])

  // const [queueNames, setQueueNames] = useState<string[]>([])
  
  // const addQueueNameToQueueNames = () => {
  //   const _queueNames = [...queueNames]
  //   _queueNames.push(queueName)
  //   setQueueNames(_queueNames)
  //   setQueueName("")
  // }

  // const handleDeleteQueueName = (deletedQueueName: string) => {
  //     var _queueNames = queueNames.filter((value, index, error): boolean => {
  //       return value != deletedQueueName
  //     })
  //     setQueueNames(_queueNames)
  // }

  // =================

  useEffect(() => {
    if (selectedQueueId === null) {
      setMainContent((
        <>
          <Container maxWidth="md">
            <Container fixed>
              <Typography gutterBottom variant="h5" component="h2" align="center">
                Please fill the form to create customers.
                <Typography 
                  gutterBottom 
                  style={{whiteSpace: 'pre-line'}}
                >
                  {storeInfo.description}
                </Typography>
              </Typography>
            </Container>
            <Grid container rowSpacing={2} justifyContent="center" alignItems="center">
              <Grid item key={"all"} xs={10} sm={10} md={6}>
                <Box sx={{ mt: 1 }}>
                  <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="name"
                    label="Name"
                    name="name"
                    autoComplete="name"
                    onChange={handleInputCustomerName}
                    error={customerNameAlertFlag}
                    helperText="Can not add more than 5 customers at a time."
                  />
                  <TextField
                    margin="normal"
                    fullWidth
                    name="phone"
                    label="Phone"
                    type="phone"
                    id="phone"
                    autoComplete="tel"
                    onChange={handleInputCustomerPhone}
                    error={customerPhoneAlertFlag}
                  />
                  <Grid 
                    container 
                    spacing={2}
                    alignItems="center"
                    justifyContent="flex-start"
                  >
                    <Grid item xs={8} sm={8}>
                      <Box component="form" sx={{ display: 'flex', flexWrap: 'wrap' }}>
                        <FormControl 
                          sx={{ mt: 3, minWidth: 160 }} 
                          // error={queueNameAlertFlag}
                        >
                          <InputLabel id="queue-select-label">Queue</InputLabel>
                          <Select
                            labelId="queue-select-label"
                            id="queue-select"
                            // value={customerNewStatus}
                            // onChange={handleChangeCustomerNewStatus}
                            input={<OutlinedInput label="Queue" />}
                          >
                            <MenuItem value="">------</MenuItem>
                            {queuesInfo.map((queue: Queue) => (
                              <MenuItem value={queue.id}>{queue.name}</MenuItem>
                            ))}
                          </Select>
                        </FormControl>
                      </Box>
                    </Grid>
                    <Grid item xs={4} sm={4}>
                      <Button 
                        variant="contained" 
                        startIcon={<AddBoxIcon />}
                        // onClick={addQueueNameToQueueNames}
                        // disabled={addQueueNameAlertFlag}
                      >
                        Add
                      </Button>
                    </Grid>
                  </Grid>              

                  {/* {queueNames.map((queueName: string) => (
                      <Chip 
                        sx={{ mb: 1, ml: 1, mr: 1 }}
                        label={queueName}
                        key={queueName} 
                        onDelete={() => {handleDeleteQueueName(queueName)}}
                      />
                    ))} */}

                  <Button
                    fullWidth
                    variant="contained"
                    sx={{ mt: 3, mb: 2 }}
                    // onClick={doMakeOpenStoreRequest}
                  >
                    Add Customers
                  </Button>
                </Box>
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
                            waiting
                          </TableCell>
                        )}
                        {customer.status === 'processing' && (
                          <TableCell align="right" sx={{color: 'red'}}>
                            {customer.status}
                          </TableCell>
                        )}
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
  }, [
    selectedQueueId, 
    setMainContent, 
    storeInfo, 
    queuesInfo,
  ])

  // ================= scan session =================  
  const [scanSessionAction, makeScanSessionRequest] = useApiRequest(...scanSession(sessionId, parseInt(storeId)))
  useEffect(() => {
    makeScanSessionRequest()
  }, [])
  useEffect(() => {
    if (scanSessionAction.actionType === ACTION_TYPES.ERROR) {
        // TODO: disable create customers block
    }

    // 40007: store_session exist but is already scanned.
    if ((scanSessionAction.response != null) && (scanSessionAction.response["error_code"]) && (scanSessionAction.response["error_code"] !== 40007)) {
        // TODO: disable create customers block
    }
  }, [scanSessionAction.actionType])

  // const [customerName, setCustomerName] = useState("")
  // const [customerNameAlertFlag, setCustomerNameAlertFlag] = useState(false)
  // const handleInputCustomerName = (e: React.ChangeEvent<HTMLInputElement>) => {
  //   const { value }: { value: string } = e.target
  //   if (value) {
  //     setCustomerNameAlertFlag(false)
  //     setCustomerName(value)
  //   } else {
  //     setCustomerNameAlertFlag(true)
  //   }
  // }

  // const [customerPhone, setCustomerPhone] = useState("")
  // const [customerPhoneAlertFlag, setCustomerPhoneAlertFlag] = useState(false)
  // const handleInputCustomerPhone = (e: React.ChangeEvent<HTMLInputElement>) => {
  //   const { value }: { value: string } = e.target
  //   if (value) {
  //     setCustomerPhoneAlertFlag(false)
  //     setCustomerPhone(value)
  //   } else {
  //     setCustomerPhoneAlertFlag(true)
  //   }
  // }

  // const [queueId, setQueueId] = useState<number>(0)
  // const handleInputQueueId = (e: React.ChangeEvent<HTMLInputElement>) => {
  //   const { value }: { value: string } = e.target
  //   if (value) {
  //     setQueueId(parseInt(value))
  //   }
  // }

  // const [addCustomerFlag, setAddCustomerFlag] = useState(false)
  // useEffect(() => {
  //   if (customerName && queueId) {
  //     setAddCustomerFlag(false)
  //   } else {
  //     setAddCustomerFlag(true)
  //   }
  // }, [customerName, queueId])

  // const [customersForm, setCustomersForm] = useState<CustomerForm[]>([])

  // const addCustomer = () => {
  //   const _customersForm = [...customersForm]
  //   _customersForm.push({
  //       name: customerName,
  //       phone: customerPhone,
  //       queue_id: queueId
  //   })
  //   setCustomersForm(_customersForm)

  //   setCustomerName("")
  //   setCustomerPhone("")
  //   setQueueId(0)
  // }

  // const [createCustomersAction, makeCreateCustomersRequest] = useApiRequest(
  //   ...createCustomers(sessionId, parseInt(storeId), customersForm)
  //   )

  // const [createCustomersFlag, setCreateCustomersFlag] = useState(false)
  
  return (
    <AppBarWDrawer
      storeInfo={storeInfo}
      mainContent={mainContent}
      setSelectedQueueId={setSelectedQueueId}
      queuesInfo={queuesInfo}
      StoreDrawer={(<></>)}
    />
    // <div>
    //     <div>
    //         scanned qrcode and create customers
    //     </div>

    //     <input
    //         type="text"
    //         placeholder="customer name"
    //         className={classNames({'alertInputField': customerNameAlertFlag})}
    //         onBlur={handleInputCustomerName}
    //     />
    //     <input
    //         type="text"
    //         onBlur={handleInputCustomerPhone}
    //         className={classNames({'alertInputField': customerPhoneAlertFlag})}
    //         placeholder="customer phone"
    //     />
    //     <input
    //         type="number"
    //         onBlur={handleInputQueueId}
    //         placeholder="queue id"
    //     />

    //     <button 
    //         onClick={addCustomer}
    //         disabled={addCustomerFlag}
    //     >
    //         add customer
    //     </button>

    //     {customersForm.map((customerForm: CustomerForm) => (
    //       // <div id={customer.name} key={customer.name}>{customer.name}</div>
    //       <div id={customerForm.name} key={customerForm.queue_id}>{customerForm.name}</div>
    //     ))}

    //     <br />
    //     <button onClick={clearCustomersForm}>
    //         clear customers
    //     </button>

    //     &nbsp;&nbsp;&nbsp;

    //     <button 
    //         onClick={makeCreateCustomersRequest}
    //         disabled={createCustomersFlag}
    //     >
    //         create customers
    //     </button>
    // </div>
  )
}

export {
    CreateCustomers
}