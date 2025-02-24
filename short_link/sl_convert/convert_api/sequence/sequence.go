package sequence


// Sequence 取号器接口
type Sequence interface {
	NextNumber() (uint64, error)
}