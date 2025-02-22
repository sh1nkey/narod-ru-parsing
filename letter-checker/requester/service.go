package requester

type Checker interface {
	ReadOne(letters string) int
	WriteOne(letters string)
}

func NewCheckerService(repo CheckerRepo) Checker {
	return &сheckerService{repo: repo}
}

type сheckerService struct {
	repo CheckerRepo
}

const NotFountStatus = 204
const FountStatus = 409

const MaxRetry = 10

func (cs *сheckerService) ReadOne(letters string) int {
	isSuccess := cs.repo.ReadOne(letters)

	if isSuccess {
		return NotFountStatus
	}
	return FountStatus

}

func (cs *сheckerService) WriteOne(letters string) {
	cs.repo.WriteOne(letters)
}