package domain

type Get struct {
	Id      int64
	Type    string
	Date    string
	Content string
}

type Lesson struct {
	Id   int
	Time string
	Name string
}

type ByGroup struct {
	Id     int
	Type   string
	Date   string
	Lesson Lesson
}
