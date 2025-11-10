# Debugging Guide

This guide is aimed at developers using JetBrains GoLand IDE.

**Steps**:
1. Set up a Go Remote to `localhost` port `2345`
2. Run the command `make debug LOG_HTTP_TRACE=true LOG_LEVEL=trace LOG_TRACER=true`
3. Run the remote debugger
4. Set breakpoints
5. Run the tests

To stop the app running, press `CTRL+C` in the terminal window. You might also need to terminate
the terminal window.