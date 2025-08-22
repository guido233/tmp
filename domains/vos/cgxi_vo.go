package vos

import "go-app/domains/bos"

type CgxiJxbDataVo struct {
	bos.JxbBo
	Data CgxiJxbVo `json:"data"`
}

type CgxiJxbVo struct {
	CgxiJxbTcp
	CgxiJoint
}

type CgxiJxbTcp struct {
	X  string `json:"x"`
	Y  string `json:"y"`
	Z  string `json:"z"`
	Rx string `json:"rx"`
	Ry string `json:"ry"`
	Rz string `json:"rz"`
}

type CgxiJoint struct {
	Joint1 string `json:"joint1"`
	Joint2 string `json:"joint2"`
	Joint3 string `json:"joint3"`
	Joint4 string `json:"joint4"`
	Joint5 string `json:"joint5"`
	Joint6 string `json:"joint6"`
}
