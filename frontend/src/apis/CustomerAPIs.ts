import * as httpTools from './base'

interface Customer {
    name: string
    phone: string
    queue_id: number
  }

  const createCustomers = (sessionId: string, storeId: number, customers: Customer[]): [string, httpTools.RequestParams] => {
    const jsonBody: string = JSON.stringify({
        "store_id": storeId,
        "customers": customers
    })
    return [
        httpTools.generateURL("/customers"), { 
            method: httpTools.HTTPMETHOD.POST,
            headers: {...httpTools.CONTENT_TYPE_JSON, ...httpTools.generateAuth(sessionId, false)},
            body: jsonBody
        }
    ]
}

const updateCustomer = (customerId: number, normalToken: string, storeId: number, queueId: number, oldCustomerStatus: string, newCustomerStatus: string): [string, httpTools.RequestParams] => {
    const route = "/customers/".concat(customerId.toString())
    const jsonBody: string = JSON.stringify({
        "store_id": storeId,
        "queue_id": queueId,
        "old_customer_status": oldCustomerStatus,
        "new_customer_status": newCustomerStatus,
    })
    return [
        httpTools.generateURL(route), { 
            method: httpTools.HTTPMETHOD.PUT,
            headers: {...httpTools.CONTENT_TYPE_JSON, ...httpTools.generateAuth(normalToken)},
            body: jsonBody
        }
    ]
}

export {
    createCustomers,
    updateCustomer
}