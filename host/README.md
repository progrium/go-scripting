# host / rpc extensions experiment

This is eventually going to be a sub project that wraps go-scripting in a binary that connects the scripted extensions with a remote application. This binary could be run by go-coproc from the application.

However, for now it's just testing how you can expose go-extensions across applications using Duplex RPC. 

The `demo` project folder is a pretend application that has an `EventObserver` extension point. It will allow other applications to extend it, in this case the application here called `host`.

`demo` provides an `EventObserver` extension point and uses this to produce an event every 5 seconds using the single method `Event`. It also binds a Duplex socket exposing a `Extensions.Register` service. This is for applications to connect to that want to remotely register against the extension point. 

`host` implements an `EventPrinter` with the `EventObserver` interface which just prints out the event payload. It exposes this as a service over Duplex RPC. It calls the `Extensions.Register` service on its remote peer to register it.

`demo` implements a `RemoteProxy` that acts as a wrapper for making the RPC calls for any remote extension to `EventObserver`. This is unfortunate boilerplate and is tightly coupled to the interfaces a remote program wants to implement. There is a commented out attempt at a dynamic version of `RemoteProxy`, however you can see that the static typing involved in RPC calls does not match with the `interface{}` based generic function signature of the proxy. 