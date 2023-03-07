package bloomfilter

type Hash interface {
	Sum() int
}
