package logger

// Router route log messages to multiple streams
type Router struct {
	ifaces []OutputIface
}

// NewRouter create a *Router pointer
func NewRouter(route ...OutputIface) *Router {
	r := &Router{
		ifaces: route,
	}
	return r
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (r *Router) Output(calldepth int, s string) error {
	for _, iface := range r.ifaces {
		err := iface.Output(calldepth+1, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// Error error level message.
func (r *Router) Error(s string) {
	for _, iface := range r.ifaces {
		iface.Error(s)
	}
}

// Errorf error level message.
func (r *Router) Errorf(format string, v ...interface{}) {
	for _, iface := range r.ifaces {
		iface.Errorf(format, v...)
	}
}

// Debug debug level message.
func (r *Router) Debug(s string) {
	for _, iface := range r.ifaces {
		iface.Debug(s)
	}
}

// Debugf debug level message.
func (r *Router) Debugf(format string, v ...interface{}) {
	for _, iface := range r.ifaces {
		iface.Debugf(format, v...)
	}
}

// Printf print error message
func (r *Router) Printf(format string, v ...interface{}) {
	for _, iface := range r.ifaces {
		iface.Printf(format, v...)
	}
}
