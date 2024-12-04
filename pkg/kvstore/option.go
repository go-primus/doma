package kvstore

type putOpt struct{}
type PutOption func(*putOpt)

type deleteOpt struct{}
type DeleteOption func(*deleteOpt)

type updateOpt struct{}
type UpdateOption func(*updateOpt)

type getOpt struct{}
type GetOption func(*getOpt)

type listOpt struct{}
type ListOption func(*listOpt)

type countOpt struct{}
type CountOption func(*countOpt)

type watchOpt struct{}
type WatchOption func(*watchOpt)
