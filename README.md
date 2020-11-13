# httpadapt
Adapt a http.Handler to a Lambda Gateway Event handler

## features
- Stdlib and lambda event deps only
- Configurable logging
- Only supports context based handling
- Deterministic query params order

## backlog
- [ ] Test errors, possiblty with a package error type
- [ ] Add an functional option to configure a logger
- [ ] Add a functional option to configure stripbasepath
- [ ] Add a functional option for CustomHostVariable env
- [ ] Consider the v2 api format
- [ ] Add support for non base64 encoded bodies

## existing
- [ ] https://github.com/apex/gateway
- [ ] https://github.com/akrylysov/algnhsa