package user

type userService struct {
	db     dataStorage
	media  mediaService
	search searchRepository
}

func NewUserService(db dataStorage, media mediaService, search searchRepository) *userService {
	return &userService{
		db:     db,
		media:  media,
		search: search,
	}
}
