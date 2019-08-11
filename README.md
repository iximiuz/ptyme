# ptyme - daemonize but keep terminal alive

Simple demonstration of Linux <a href="http://man7.org/linux/man-pages/man7/pty.7.html">PTY</a>
capabilities. PTY is an ancient yet ubiquitous technology. It powers SSH, docker, kubernetes, etc.

This project contains trivial implementation of <a href="https://docs.docker.com/engine/reference/commandline/attach/">attach</a>/<a href="https://docs.docker.com/engine/reference/commandline/exec/">exec</a> (or kubectl <a href="https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#attach">attach</a>/<a href="https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#exec">exec</a>) feature.

```
Server:
    +-----------+                                      +----------------+
    |  shim.c   | <-- [pty] -- read/write -- [pts] --> |   ping ya.ru   |
    +-----------+                                      +----------------+
          |
          |
      [network]
          |
Client:   |
    +-----------+
    | attach.go | <-- [terminal in RAW mode] --> user via xterm (iterm2, etc).
    +-----------+
```

The idea is simple: we want to start an arbitrary executable in background (i.e. as a daemon), but keep its STDIN/STDOUT bound to a controlling terminal to be able to connect to it later on. For that we need a tiny piece of software called _a shim_ (<a href="https://github.com/iximiuz/ptyme/blob/master/shim.c">shim.c</a>). The shim creates a pseudoterminal and `fork/exec`-s a given executable binding its standard streams to the slave side of the pseudoterminal pair. At the same time the parent process keeps reading and writing the master end of the pair. The parent process also starts listening on TCP port. Each byte read from an incomming connection is then forwarded to a master side of the terminal. And other way around - each byte read from the master end of the pseudoterminal has to be written to each incomming connection (i.e. broadcasted).

Simplistic client can be found in <a href="https://github.com/iximiuz/ptyme/blob/master/attach.go">attach.go</a>. Since the controlling (i.e. escape sequence handling, etc) of user interaction is done by the pseudoterminal on the server side, the client sets its controlling terminal to RAW mode and then just blindly forwards bytes from its STDIN to a socket connection and from the connection to its STDOUT.
