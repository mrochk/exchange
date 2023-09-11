# Exchange
A market exchange simulation that can be called through a JSON API.

This is a market exchange I made from scratch, it does not not rely any third-party dependency. \
The engine implements Price/Time (FIFO) algorithm for matching orders. 

For now, it is only possible to place "basic" types of orders (market and limit).

Here are the different API endpoints:

- `/getorderbooks`: returns the ID of all available orderbooks
- `/register`: returns the ID the user will use to place orders
- `/placeorder`: place a limit order and returns the order ID that can be then used to cancel it
- `/cancelorder`: cancel a limit order if it has not already been executed
- `/executeorder`: execute a market order
- `/getorderbookdata`: returns informations about the order-book