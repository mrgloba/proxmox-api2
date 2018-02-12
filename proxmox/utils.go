package proxmox

import "strings"

func parseKeyPairs(str string) [][]string {
	var keypairs [][]string

	a1 := strings.Split(str,",")

	for _,v := range a1 {
		if strings.Index(v,"=") >0 {
			a2 := strings.Split(v, "=")
			keypairs = append(keypairs,a2)
		} else {
			keypairs = append(keypairs,[]string{v})
		}
	}

	return keypairs
}