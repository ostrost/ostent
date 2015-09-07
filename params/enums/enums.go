package enums

// Constants for DF sorting criterion.
const (
	_     int = iota
	FS        // 1
	MP        // 2
	AVAIL     // 3
	USED      // 4
	TOTAL     // 5
)

// Constants for DF tabs.
const (
	_       int = iota
	INODES      // 1
	DFBYTES     // 2
)

// Constants for PS sorting criterion.
const (
	_    int = iota
	PID      // 1
	UID      // 2
	USER     // 3
	PRI      // 4
	NICE     // 5
	VIRT     // 6
	RES      // 7
	TIME     // 8
	NAME     // 9
)
