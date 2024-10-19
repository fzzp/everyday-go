package main

type Blog struct {
	Title   string
	Content string
}

func (b *Blog) Add()    {}
func (b *Blog) Update() {}
func (b *Blog) Delete() {}

// 这块就不合法单一职责
// func (b *Blog) Save(filename string) {}

type Store struct {
	Filename string
}

// 提取出来
func (s *Store) Save(blog *Blog) {

}
