package user

import "fmt"

type Info struct {
	Name string `json:"name"`
	Id   uint64 `json:"id"`
}

var db = map[int]Info{

	1: {
		Name: "ss",
		Id:   1,
	},
	2: {
		Name: "ss",
		Id:   2,
	},
	3: {
		Name: "ss",
		Id:   3,
	},
}

type UserService struct {
}

func (t *UserService) SayHello(s string) (string, error) {
	return s, nil
}

func (t *UserService) GetUserIds() ([]int, error) {
	return []int{1, 2, 3}, nil
}

func (t *UserService) GetUserInfoById(id int) (Info, error) {

	if info, ok := db[id]; ok {
		return info, nil
	}
	return Info{}, fmt.Errorf("user not exist")
}
