# Exchange
A market exchange simulation that can be called through a JSON API.

I made it from scratch, it does not not rely any third-party dependency. \
The engine implements the Price/Time (First In First Out) algorithm for matching orders. 

For now, it is only possible to place the two "basic" types of orders (market and limit).\
Also, canceling orders is not fully implemented yet, at least not efficiently.

Here are the different API endpoints you can call once the server is launched:

- `/orderbooks`: returns the list of all open orderbooks
- - GET
- `/register`: returns the ID the user must use to place orders
- - POST
- - Body:
- - - "username": string
- `/placeorder`: place a limit order and return the order ID
- - POST
- - Body:
- - - "base": string (e.g "USD")
- - - "quote": string (e.g "EUR")
- - - "type": string ("BUY" | "SELL")
- - - "quantity": float
- - - "issuer": string (e.g "Maxime")
- - - "price": float 
- `/cancelorder`: cancel a limit order if it has not already been executed
- - POST
- - Body:
- - - "base": string 
- - - "quote": string
- - - "type": string
- - - "price": float 
- - - "order_id": int
- `/executeorder`: execute a market order
- - POST
- - Body:
- - - "base": string
- - - "quote": string 
- - - "type": string 
- - - "quantity": float
- - - "issuer": string
- `/orderbookdata`: returns informations about the order-book
- - GET
- - Body:
- - - "base": string
- - - "quote": string 