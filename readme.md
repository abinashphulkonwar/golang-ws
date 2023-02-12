# WS Golnag

WS Golnag is a scalable WebSocket service designed to handle a large user base. The system consists of two components:

1. **Node**: This component acts as the WebSocket server that handles TCP connections.

2. **Master**: This component manages the Nodes and the user data.

# How it Works

When a user connects to a Node, the Node sends a request to the Master to store the user's data. When a user sends a message to another user, the connected Node checks its internal cache to determine if the other user is connected to the same Node. If the other user is connected to the same Node, the message is sent directly to that user.

If the other user is not connected to the same Node, the connected Node checks the other users' cache. If the cache contains the other user's data, the Node sends an event to the other user's connected Node about the message, and that Node sends the message to the user.

If the other user's cache is not found or is invalid, the connected Node sends a request to the Master, which returns the user's connected Node data. The Node caches this data and then sends the message event to the other Node for delivery to the user.

# Conclusion

With its efficient caching and distributed architecture, WS Golnag is an effective solution for handling large numbers of WebSocket connections.
