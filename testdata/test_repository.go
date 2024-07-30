package testdata

type TestIndexExprRepository[T any] interface {
	Method() T
	Method1() T
}

type TestIndexListExprRepository[T any, T1 any] interface {
	Method(T) T1
	Method1(T) T1
}

type TestRepository[T string, T1 string] interface {
	TestIndexExprRepository[T]
	TestIndexListExprRepository[T, T1]

	SelfMethod(T) T1
}
