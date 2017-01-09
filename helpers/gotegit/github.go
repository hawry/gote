package gotegit

type Github struct {
}

type Bitbucket struct {
}

type Gitlab struct {
}

//CreateIssue implements the GitProvider interface
func (g *Github) CreateIssue(accessToken string) (bool, error) {
	return false, nil
}

//ProviderName implements the GitProvider interface
func (g *Github) ProviderName() string {
	return "github"
}

//CreateIssue implements the GitProvider interface
func (b *Bitbucket) CreateIssue(accessToken string) (bool, error) {
	return false, nil
}

//ProviderName implements the GitProvider interface
func (b *Bitbucket) ProviderName() string {
	return "bitbucket"
}

//CreateIssue implements the GitProvider interface
func (g *Gitlab) CreateIssue(accessToken string) (bool, error) {
	return false, nil
}

//ProviderName implements the GitProvider interface
func (g *Gitlab) ProviderName() string {
	return "gitlab"
}
