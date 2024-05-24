package cli

type BudCLI struct {
	IBudCLI
}

type IBudCLI interface {
	registerDir()
	processDir()
	processFile()
	version()
	help()
	budStatus()
}

func New() *BudCLI {
	return &BudCLI{}
}

func (cli *BudCLI) registerDir() {}
func (cli *BudCLI) processDir()  {}
func (cli *BudCLI) processFile() {}
func (cli *BudCLI) version()     {}
func (cli *BudCLI) hep()         {}
func (cli *BudCLI) budStatus()   {}
