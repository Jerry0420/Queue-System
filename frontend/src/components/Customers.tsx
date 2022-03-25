import React, { useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import { createCustomers } from "../apis/CustomerAPIs"
import { CustomerForm } from "../apis/models"
import { ACTION_TYPES, useApiRequest } from "../apis/reducer"
import { scanSession } from "../apis/SessionAPIs"
import { getStoreInfoWithSSE } from "../apis/StoreAPIs"
import classNames from 'classnames'

const CreateCustomers = () => {
  let { storeId , sessionId}: {storeId: string, sessionId: string} = useParams()
  const [queues, setQueues] = useState<{id: number, name: string}[]>([])

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

  useEffect(() => {
    let getStoreInfoSSE: EventSource
    getStoreInfoSSE = getStoreInfoWithSSE(parseInt(storeId))

    getStoreInfoSSE.onmessage = (event) => {
        // TODO: render customers ui
        const _queues: {id: number, name: string}[] = []
        JSON.parse(event.data)["queues"].forEach((queue: {id: number; name: string}) => {
            _queues.push({
                id: queue["id"],
                name: queue["name"]
            })
            setQueues(_queues)
        })

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

  const [customerName, setCustomerName] = useState("")
  const [customerNameAlertFlag, setCustomerNameAlertFlag] = useState(false)
  const handleInputCustomerName = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    if (value) {
      setCustomerNameAlertFlag(false)
      setCustomerName(value)
    } else {
      setCustomerNameAlertFlag(true)
    }
  }

  const [customerPhone, setCustomerPhone] = useState("")
  const [customerPhoneAlertFlag, setCustomerPhoneAlertFlag] = useState(false)
  const handleInputCustomerPhone = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    if (value) {
      setCustomerPhoneAlertFlag(false)
      setCustomerPhone(value)
    } else {
      setCustomerPhoneAlertFlag(true)
    }
  }

  const [queueId, setQueueId] = useState<number>(0)
  const handleInputQueueId = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { value }: { value: string } = e.target
    if (value) {
      setQueueId(parseInt(value))
    }
  }

  const [addCustomerFlag, setAddCustomerFlag] = useState(false)
  useEffect(() => {
    if (customerName && queueId) {
      setAddCustomerFlag(false)
    } else {
      setAddCustomerFlag(true)
    }
  }, [customerName, queueId])

  const [customersForm, setCustomersForm] = useState<CustomerForm[]>([])

  const addCustomer = () => {
    const _customersForm = [...customersForm]
    _customersForm.push({
        name: customerName,
        phone: customerPhone,
        queue_id: queueId
    })
    setCustomersForm(_customersForm)

    setCustomerName("")
    setCustomerPhone("")
    setQueueId(0)
  }

  const [createCustomersAction, makeCreateCustomersRequest] = useApiRequest(
    ...createCustomers(sessionId, parseInt(storeId), customersForm)
    )

  const [createCustomersFlag, setCreateCustomersFlag] = useState(false)

  useEffect(() => {
    if (customersForm.length > 0) {
        setCreateCustomersFlag(false)
    } else {
        setCreateCustomersFlag(true)
    }
  }, [customersForm])

  const clearCustomersForm = () => {
    setCustomersForm([])
  }
  
  return (
    <div>
        <div>
            scanned qrcode and create customers
        </div>

        <input
            type="text"
            placeholder="customer name"
            className={classNames({'alertInputField': customerNameAlertFlag})}
            onBlur={handleInputCustomerName}
        />
        <input
            type="text"
            onBlur={handleInputCustomerPhone}
            className={classNames({'alertInputField': customerPhoneAlertFlag})}
            placeholder="customer phone"
        />
        <input
            type="number"
            onBlur={handleInputQueueId}
            placeholder="queue id"
        />

        <button 
            onClick={addCustomer}
            disabled={addCustomerFlag}
        >
            add customer
        </button>

        {customersForm.map((customerForm: CustomerForm) => (
          // <div id={customer.name} key={customer.name}>{customer.name}</div>
          <div id={customerForm.name} key={customerForm.queue_id}>{customerForm.name}</div>
        ))}

        <br />
        <button onClick={clearCustomersForm}>
            clear customers
        </button>

        &nbsp;&nbsp;&nbsp;

        <button 
            onClick={makeCreateCustomersRequest}
            disabled={createCustomersFlag}
        >
            create customers
        </button>
    </div>
  )
}

export {
    CreateCustomers
}