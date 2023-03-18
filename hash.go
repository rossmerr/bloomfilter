package bloomfilter

type Hash interface {
	Sum() uint
}
