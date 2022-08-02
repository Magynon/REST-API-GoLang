package domain

import (
	"hash/fnv"
)

type Element struct {
	Name      string   `json:"name"`
	Id      string   `json:"id"`
	Price       int      `json:"price"`
	Quantity	int		`json:"quantity"`
	Country string `json:"country"`
}

func (u *Element) GetHash() int {
	h := fnv.New32a()
	h.Write([]byte(u.Name + u.Id))
	return int(h.Sum32())
}
