package some

type Usecase interface {
	Execute() error
	SumMethod(a, b int) int
}

type usecase struct {

}

func NewUsecase() Usecase {
	return &usecase{}
}

func (u *usecase) Execute() error {
	// Implementation of the use case logic
	return nil
}

func (u *usecase) SumMethod(a, b int) int {
	return a + b
}