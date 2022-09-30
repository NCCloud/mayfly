package collectors

type Collector interface {
	Collect() error
}
