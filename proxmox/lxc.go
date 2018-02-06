package proxmox


type Lxc struct {
	Cpu float64 	`json:"cpu"`
	Cpus int 		`json:"cpus,string"`
	Disk int64		`json:"disk"`
	DiskRead int64	`json:"diskread,string"`
	DiskWrite int64	`json:"diskwrite,string"`
	MaxDisk	int64	`json:"maxdisk"`
	MaxMem	int64	`json:"maxmem"`
	MaxSwap int64	`json:"maxswap"`
	Mem int64		`json:"mem"`
	Name string		`json:"name"`
	NetIn int64		`json:"netin"`
	NetOut int64	`json:"netout"`
	Pid int			`json:"pid,string"`
	Status string	`json:"status"`
	Swap int64		`json:"swap"`
	Template string	`json:"template"`
	Type string		`json:"type"`
	Uptime int64	`json:"uptime"`
	VmId int64		`json:"vmid,string"`

	BasicObject
}