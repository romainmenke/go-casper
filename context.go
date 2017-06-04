package casper

import "context"

type ctxKey int

const (
	casperKey ctxKey = iota
	hashValuesKey
)

func CasperFromContext(ctx context.Context) *Casper {
	val := ctx.Value(casperKey)
	if val != nil {
		return val.(*Casper)
	}
	return nil
}

func contextWithCasper(ctx context.Context, val *Casper) context.Context {
	if val == nil {
		return ctx
	}
	return context.WithValue(ctx, casperKey, val)
}

func HashValuesFromContext(ctx context.Context) []uint {
	val := ctx.Value(hashValuesKey)
	if val != nil {
		return val.([]uint)
	}
	return nil
}

func contextWithHashValues(ctx context.Context, val []uint) context.Context {
	if val == nil {
		return ctx
	}
	return context.WithValue(ctx, hashValuesKey, val)
}
