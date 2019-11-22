package rand

// 此为 Cgo 创建的特殊的命名空间。目的：用来与 C 的命名空间交流。
// 以下两个函数分别调用 C 的函数，对它们进行类型转换。这样就实现了 Go 里调用 C 的功能。

/*
#include <stdlib.h>
*/
import "C"

func Random() int {
	return int(C.random())
}

func Seed(i int) {
	C.srandom(C.uint(i))
}