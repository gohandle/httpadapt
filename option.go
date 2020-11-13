package httpadapt

// Option configures the adapter
type Option func(a *Adapter)

// StripBasePath configures the adapter to strip the prefixing part of the
// resulting url path
func StripBasePath(base string) Option {
	return func(a *Adapter) {
		a.stripBasePath = base
	}
}

// CustomHost configures the custom hostname for the request. If this option
// is not set the framework reverts to `RequestContext.DomainName`. The value
// for a custom host should include a protocol: http://my-custom.host.com
func CustomHost(host string) Option {
	return func(a *Adapter) {
		a.customServerAddress = host
	}
}
