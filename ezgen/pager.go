package ezgen

// Pager 表示一个可执行分页的参数请求
type Pager interface {
	// GetPageSize 分页大小, 需要大于 0
	GetPageSize() uint32

	// GetPageIndex 分页页码, 需要大于 0
	GetPageIndex() uint32
}
