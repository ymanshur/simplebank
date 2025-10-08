package common

const (
	DefaultPageID   = 1
	DefaultPageSize = 25
)

func PageID(id int32) int32 {
	if id < 1 {
		return DefaultPageID
	}
	return id
}

func PageSize(size int32) int32 {
	if size < 1 || size > DefaultPageSize {
		return DefaultPageSize
	}
	return size
}
