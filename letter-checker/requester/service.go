package requester


type Checker interface {
	Check(letters string) int
}

func NewCheckerService(repo CheckerRepo) Checker {
	return  &сheckerService{repo: repo}
}


type сheckerService struct {
	repo CheckerRepo
}


const SuccessWriteStatusCode = 204
const FailedWriteStatusCode = 409

const MaxRetry = 10
func (cs *сheckerService) Check(letters string) int {
	isSuccess := cs.repo.CheckExisted(letters)

	if isSuccess == true {
		return SuccessWriteStatusCode
	}
	return FailedWriteStatusCode

}