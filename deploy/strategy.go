package deploy

// Strategy is an interface that describes how to deploy a repository and how to commit changes
type Strategy interface {
	Deploy() error
	Commit() error
}

// Master is a strategy that always deploys from and commits to the master branch
type Master struct{}

func (m Master) Deploy() error {
	return nil
}

func (m Master) Commit() error {
	return nil
}

// ConfigurableBranches is a strategy that deploys from and commits to individually configurable branches
type ConfigurableBranches struct {
	DeployBranch, CommitBranch string
}

func (c ConfigurableBranches) Deploy() error {
	return nil
}

func (c ConfigurableBranches) Commit() error {
	return nil
}
