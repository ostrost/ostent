package enums

// Constants for DF sorting criterion.
const (
	_     int = iota
	FS        // 1
	MP        // 2
	TOTAL     // 3
	USED      // 4
	AVAIL     // 5
)

// Constants for DF tabs.
const (
	_       int = iota
	INODES      // 1
	DFBYTES     // 2
)

// Constants for IF tabs.
const (
	_       int = iota
	PACKETS     // 1
	ERRORS      // 2
	IFBYTES     // 3
)

// Constants for PS sorting criterion.
const (
	_    int = iota
	PID      // 1
	PRI      // 2
	NICE     // 3
	VIRT     // 4
	RES      // 5
	TIME     // 6
	NAME     // 7
	UID      // 8
	USER     // 9
)
