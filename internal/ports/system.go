package ports

type SystemProviderPort interface {
	OrganizeFolder(path string) (int, error)
}