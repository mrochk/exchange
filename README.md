# Exchange
A market exchange simulation that can be called through a JSON API.

<img src="sc.png" alt="Screenshot" width="800" height="">


It is made from scratch, it does not not rely any third-party dependency. \
The matching engine implements the Price/Time (First In First Out) algorithm. 

For now, it is only possible to place the two "basic" types of orders (market and limit).\
Canceling orders is not fully implemented yet, at least not efficiently.\
Also, this project will be rewritten in a much cleaner way and with better error handling when I'll have some time.

Here are the different API endpoints you can call once the server is launched,\
with the method you must use and the expected body JSON fields:

- `/orderbooks`: returns the list of all open orderbooks
- - GET
- `/placeorder`: places a limit order and returns its ID
- - POST
- - Body:
- - - "base": string (e.g "USD")
- - - "quote": string (e.g "EUR")
- - - "type": string ("BUY" | "SELL")
- - - "quantity": float
- - - "issuer": string (e.g "Maxime")
- - - "price": float 
- `/cancelorder`: cancels a limit order if it has not already been executed
- - POST
- - Body:
- - - "base": string 
- - - "quote": string
- - - "type": string
- - - "price": float 
- - - "order_id": int
- `/executeorder`: executes a market order
- - POST
- - Body:
- - - "base": string
- - - "quote": string 
- - - "type": string 
- - - "quantity": float
- - - "issuer": string
- `/orderbookdata`: returns information about the order book
- - GET
- - Body:
- - - "base": string
- - - "quote": string 
