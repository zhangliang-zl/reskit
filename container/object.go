package container

type Destructible interface {
	Destroy(...interface{}) error
}

// Manager 负责管理带有析构函数的对象
type Manager interface {
	Set(k string, v interface{})
	Get(k string) (v interface{}, exist bool)
	Delete(k string) error
	Destroy() error
}

//Make 把普通对象转换为Destructible对象
func Make(v interface{}, destory ...func(...interface{}) error) {

}
